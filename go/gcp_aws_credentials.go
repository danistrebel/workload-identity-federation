package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"golang.org/x/oauth2/google/externalaccount"
)

type customAwsSecurityCredentialsSupplier struct{}

func (s customAwsSecurityCredentialsSupplier) AwsRegion(ctx context.Context, options externalaccount.SupplierOptions) (string, error) {
	// Replace with your logic to get the AWS region
	return "us-east-1", nil
}

func (s customAwsSecurityCredentialsSupplier) AwsSecurityCredentials(ctx context.Context, options externalaccount.SupplierOptions) (*externalaccount.AwsSecurityCredentials, error) {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return nil, err
	}

	credentials, err := conf.Credentials.Retrieve(ctx)
	if err != nil {
		fmt.Printf("Error retrieving AWS credentials: %v\n", err)
		return nil, err
	}

	return &externalaccount.AwsSecurityCredentials{
		AccessKeyID:     credentials.AccessKeyID,
		SecretAccessKey: credentials.SecretAccessKey,
		SessionToken:    credentials.SessionToken,
	}, nil
}

func main() {
	ctx := context.Background()

	// Read GCP Workload identity config from env variables
	projectNumber := os.Getenv("GCP_PROJECT_NUMBER")
	if projectNumber == "" {
		fmt.Println("GCP_PROJECT_NUMBER environment variable is not set")
		return
	}
	workloadPoolId := os.Getenv("WORKLOAD_IDENTITY_POOL_ID")
	if workloadPoolId == "" {
		fmt.Println("GCP_WORKLOAD_POOL_ID environment variable is not set")
		return
	}
	providerId := os.Getenv("WORKLOAD_IDENTITY_PROVIDER_ID")
	if providerId == "" {
		fmt.Println("GCP_WORKLOAD_PROVIDER_ID environment variable is not set")
		return
	}

	// Create an instance of your AWS Security Credentials Supplier
	awsSupplier := customAwsSecurityCredentialsSupplier{}

	// Create a GCP token source using the AWS credentials
	// (assumes you have the necessary GCP permissions)
	tokenSource, err := externalaccount.NewTokenSource(ctx, externalaccount.Config{
		SubjectTokenType:               "urn:ietf:params:aws:token-type:aws4_request",
		AwsSecurityCredentialsSupplier: awsSupplier,
		Audience:                       fmt.Sprintf("//iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/providers/%s", projectNumber, workloadPoolId, providerId), // Replace with your GCP project number, pool ID, and provider ID
		Scopes:                         []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		fmt.Printf("Error creating token source: %v\n", err)
		return
	}

	// Now you have a valid GCP token in `creds.TokenSource`
	token, err := tokenSource.Token()
	if err != nil {
		fmt.Printf("Error obtaining token: %v\n", err)
		return
	}

	fmt.Printf("Access token: %s\n", token.AccessToken)
}
