package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"google.golang.org/genai"
)

func main() {
	err := run(context.Background(), os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(ctx context.Context, apiKey string) error {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	genConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: "You are a helpful assistant. Keep your responses short and concise. Keep the tone friendly and approachable."},
			},
		},
	}

	max := 5

	query := "Why is the sky blue?"

	curr := 0

	var chatSoFar []*genai.Content

	var userMessage string

	for curr < max {

		if curr == 0 {
			userMessage = query
		} else {
			userMessage = "tell me more...."
		}

		chatSoFar = append(chatSoFar, &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: userMessage}},
		})

		fmt.Printf("%s: %s\n", "user", userMessage)

		result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", chatSoFar, genConfig)
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}
		output, err := renderer.Render(result.Text())
		if err != nil {
			return fmt.Errorf("failed to render output: %w", err)
		}
		fmt.Printf("%s:%s\n", "model", output)

		chatSoFar = append(chatSoFar, result.Candidates[0].Content)

		curr++
	}

	return nil
}
