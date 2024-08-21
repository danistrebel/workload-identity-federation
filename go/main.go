package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/danistrebel/workload-identity-federation/go/ecs"
	"google.golang.org/api/option"
)

func generateContent(prompt string, opts ...option.ClientOption) (string, error) {
	ctx := context.Background()

	projectId := os.Getenv("GCP_PROJECT_ID")
	if projectId == "" {
		fmt.Println("GCP_PROJECT_ID environment variable is not set")
		return "", fmt.Errorf("GCP_PROJECT_ID environment variable is not set")
	}

	region := os.Getenv("GCP_REGION")
	if region == "" {
		region = "us-central1"
	}

	client, err := genai.NewClient(ctx, projectId, os.Getenv("GCP_REGION"), opts...)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash-001")
	model.SetTemperature(0.8)

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("unable to generate content: %v", err)
	}

	ret := make([]string, len(res.Candidates[0].Content.Parts))
	for i, part := range res.Candidates[0].Content.Parts {
		ret[i] = fmt.Sprintf("%v", part)
	}
	return strings.Join(ret, ""), nil
}

func main() {
	ctx := context.Background()

	ops := []option.ClientOption{}

	// Use custom AWS Security Credentials Supplier
	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_ECS_FARGATE" {
		ecsTokenSource, err := ecs.GetECSTokenSource(ctx)
		if err != nil {
			fmt.Printf("Error getting ECS token source: %v\n", err)
			return
		}
		ops = append(ops, option.WithTokenSource(ecsTokenSource))
	}

	response, err := generateContent("Tell me a funny joke about food.", ops...)
	if err != nil {
		fmt.Printf("Error generating content: %v\n", err)
		return
	}

	fmt.Println(response)
}
