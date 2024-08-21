package ecs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google/externalaccount"
)

// CustomAwsSecurityCredentialsSupplier implements the externalaccount.Supplier interface
type customAwsSecurityCredentialsSupplier struct{}

// AwsRegion retrieves the AWS region from the environment
func (s customAwsSecurityCredentialsSupplier) AwsRegion(ctx context.Context, options externalaccount.SupplierOptions) (string, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return "", fmt.Errorf("AWS_REGION environment variable is not set")
	}
	return region, nil
}

// AwsSecurityCredentials retrieves AWS credentials from the default config
func (s customAwsSecurityCredentialsSupplier) AwsSecurityCredentials(ctx context.Context, options externalaccount.SupplierOptions) (*externalaccount.AwsSecurityCredentials, error) {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	credentials, err := conf.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AWS credentials: %w", err)
	}

	return &externalaccount.AwsSecurityCredentials{
		AccessKeyID:     credentials.AccessKeyID,
		SecretAccessKey: credentials.SecretAccessKey,
		SessionToken:    credentials.SessionToken,
	}, nil
}

func GetECSTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	// Read GCP Workload identity config from env variables
	projectNumber := os.Getenv("GCP_PROJECT_NUMBER")
	if projectNumber == "" {
		return nil, fmt.Errorf("GCP_PROJECT_NUMBER environment variable is not set")
	}

	workloadPoolId := os.Getenv("WORKLOAD_IDENTITY_POOL_ID")
	if workloadPoolId == "" {
		return nil, fmt.Errorf("GCP_WORKLOAD_POOL_ID environment variable is not set")
	}
	providerId := os.Getenv("WORKLOAD_IDENTITY_PROVIDER_ID")
	if providerId == "" {
		return nil, fmt.Errorf("GCP_WORKLOAD_PROVIDER_ID environment variable is not set")
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
		return nil, err
	}

	// For Debug purposes only
	// token, err := tokenSource.Token()
	// if err != nil {
	// 	fmt.Printf("Error obtaining token: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Access token: %s\n", token.AccessToken)

	return tokenSource, nil
}
