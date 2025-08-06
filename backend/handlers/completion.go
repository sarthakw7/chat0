package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarthakw7/chat0-backend/types"
	"google.golang.org/genai"
)


func HandleCompletion(c *gin.Context) {
	var req types.CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse {
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	// Get Google API key from header or environment
	googleAPIKey := c.GetHeader("X-Google-API-Key")
	if googleAPIKey == "" {
		// Fallback to environment variable
		googleAPIKey = os.Getenv("GOOGLE_API_KEY")
	}

	if googleAPIKey == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse {
			Error: "Google API key is required. Provide via X-Google-API-Key header or GOOGLE_API_KEY environment variable.",
		})
		return
	}

	// generate title using google gemini
	title, err := generateTitleWithGemini(req.Prompt, googleAPIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error: "failed to generate title: " + err.Error(),
		})
		return
	}

	response := types.CompletionResponse{
		Title : title,
		IsTitle: req.IsTitle,
		MessageID: req.MessageID,
		ThreadID: req.ThreadID,
	}

	c.JSON(http.StatusOK, response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generate chat title
func generateTitleWithGemini(prompt string, apiKey string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}

	systemPrompt := `
	- you will generate a short title based on the first message a user begins a conversation with
	- ensure it is not more than 80 characters long
	- the title should be a summary of the user's message
	- you should NOT answer the user's message, you should only generate a summary/title
	- do not use quotes or colons`

	// content with system instruction embedded in user prompt
	fullPrompt := fmt.Sprintf("%s\n\nUser message: %s", systemPrompt, prompt)
	userContent := genai.NewContentFromText(fullPrompt, "user")

	// generate the title
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{userContent}, &genai.GenerateContentConfig{})
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// extract the text from response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Get the text from the first part
	firstPart := resp.Candidates[0].Content.Parts[0]
	if firstPart.Text != "" {
		return firstPart.Text, nil
	}

	return "", fmt.Errorf("unexpected response format")
}