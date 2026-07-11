package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

type GenerateRequest struct {
	Prompt string
}

type GenerateResponse struct {
	Text string
}

func New(ctx context.Context, apiKey string, model string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key cannot be empty")
	}
	if model == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := g.Generate(ctx, GenerateRequest{Prompt: prompt})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

func (g *GeminiClient) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	if g.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	resp, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(req.Prompt), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("received an empty response from Gemini")
	}

	generatedText := resp.Text()

	return &GenerateResponse{
		Text: generatedText,
	}, nil
}
