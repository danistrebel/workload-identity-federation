package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody struct {
		Prompt string `json:"prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	ops := []option.ClientOption{}

	// Use custom AWS Security Credentials Supplier
	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_ECS_FARGATE" {
		ecsTokenSource, err := ecs.GetECSTokenSource(ctx)
		if err != nil {
			fmt.Printf("Error getting ECS token source: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		ops = append(ops, option.WithTokenSource(ecsTokenSource))
	}

	response, err := generateContent(reqBody.Prompt, ops...)
	if err != nil {
		fmt.Printf("Error generating content: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func main() {
	http.HandleFunc("/generate", handleGenerate)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port " + port)
	fmt.Println("example curl:")
	fmt.Println("curl localhost:" + port + "/generate -X POST -H 'Content-Type: application/json' -d '{\"prompt\": \"Write a short story about a cat.\"}'")
	http.ListenAndServe(":"+port, nil)
}
