package handlers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarthakw7/chat0-backend/models"
	"github.com/sarthakw7/chat0-backend/types"
	"google.golang.org/genai"
)

func HandleChat(c *gin.Context) {
	var req types.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ Failed to parse request: %v\n", err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	// validate atleast one message
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "at least one message is required",
		}) 
		return
	}

	// Get model config
	modelConfig, exists := models.GetModelConfig(req.Model)
	if !exists {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: "unsupported model: "+ req.Model,
		})
		return
	}

	// get the correct key for this model (header first , then environment)
	apiKey := c.GetHeader(modelConfig.HeaderKey)
	if apiKey == "" {
		// Fallback to environment variable based on provider
		switch modelConfig.Provider {
		case "google":
			apiKey = os.Getenv("GOOGLE_API_KEY")
		case "openai":
			apiKey = os.Getenv("OPENAI_API_KEY")
		case "openrouter":
			apiKey = os.Getenv("OPENROUTER_API_KEY")
		}
	}

	if apiKey == "" {
		envVar := ""
		switch modelConfig.Provider {
		case "google":
			envVar = "GOOGLE_API_KEY"
		case "openai":
			envVar = "OPENAI_API_KEY"
		case "openrouter":
			envVar = "OPENROUTER_API_KEY"
		}
		
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error: fmt.Sprintf("API key required for %s. Provide via %s header or %s environment variable.", req.Model, modelConfig.HeaderKey, envVar),
		})
		return
	}

	// Set headers for AI SDK data stream
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	w := c.Writer

	// stream AI response based on provider
	switch modelConfig.Provider {
	case "google":
		streamGoogleResponse(c.Request.Context(), w, req, apiKey, modelConfig)
	case "openrouter":
		streamOpenRouterResponse(c.Request.Context(), w, req, apiKey, modelConfig)
	default:
		// for now, using mock for other providers (OpenAI)
		streamMockResponse(w, req)
	}

	
}

func streamMockResponse(w gin.ResponseWriter, req types.ChatRequest) {
	lastMessage := req.Messages[len(req.Messages)-1]
	
	words := []string{"Mock","streaming","response","to","your","message", lastMessage.Content}

	for _, word := range words {
		// AI SDK data stream format
		textJSON, _ := json.Marshal(word + " ")
		fmt.Fprintf(w, "0:%s\n", textJSON)
		
		// Flush the data immediately
		w.Flush()
		
		// Add a small delay to simulate real streaming
		time.Sleep(200 * time.Millisecond)
	}

	// Send the final chunk to indicate completion
	fmt.Fprintf(w, "d:{\"finishReason\":\"stop\",\"usage\":{\"promptTokens\":0,\"completionTokens\":0}}\n")
	w.Flush()
}

func stringPtr(s string) *string {
	return &s
}

func streamGoogleResponse(ctx context.Context, w gin.ResponseWriter, req types.ChatRequest, apiKey string, modelConfig models.ModelConfig) {

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// create gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		writeErrorChunk(w, "Failed to create AI client: "+ err.Error())
		return
	}

	systemPrompt := `
	You are Chat0, an ai assistant that can answer questions and help with tasks.
	Be helpful and provide relevant information
	Be respectful and polite in all interactions.
	Be engaging and maintain a conversational tone.
	Always use LaTeX for mathematical expressions - 
	Inline math must be wrapped in single dollar signs:$content$
	Display math must be wrapped in double dollar signs: $$content$$
	Display math should be placed on its own line, with nothing else on that line.
	Do not nest math delimiters or mix styles.
	Examples:
	- Inline: The equation $E = mc^2$ shows mass-energy equivalence.
	- Display: 
	$$\frac{d}{dx}\sin(x) = \cos(x)$$`

	var contents []*genai.Content

	// conversation history (sys intructions will be embedded in first message)
	systemAdded := false
	for _, msg := range req.Messages{
		role := msg.Role
		if msg.Role == "assistant" {
			role = "model"
		}

		msgContent := msg.Content
		// sys instructions to first user
		if !systemAdded && msg.Role == "user" {
			msgContent = fmt.Sprintf("%s\n\nUser: %s", systemPrompt, msg.Content)
			systemAdded = true
		}

		content := genai.NewContentFromText(msgContent, genai.Role(role))
		contents = append(contents, content)
	}

	// start streaming generation
	stream := client.Models.GenerateContentStream(ctx, modelConfig.ModelID, contents, &genai.GenerateContentConfig{})

	for resp, err := range stream {
		if err != nil {
			if err == io.EOF {
				break
			}
			// Send error in AI SDK data stream format
			fmt.Fprintf(w, "3:{\"message\":\"%s\"}\n", err.Error())
			w.Flush()
			return
		}

		// process each candidate
		for _ , cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				if part.Text != "" {
					// AI SDK data stream format
					textJSON, _ := json.Marshal(part.Text)
					fmt.Fprintf(w, "0:%s\n", textJSON)
					w.Flush()
				}
			}
		}
	}

	fmt.Fprintf(w, "d:{\"finishReason\":\"stop\",\"usage\":{\"promptTokens\":0,\"completionTokens\":0}}\n")
	w.Flush()
}

func writeErrorChunk(w gin.ResponseWriter, errorMsg string) {
	fmt.Fprintf(w, "3:{\"message\":\"%s\"}\n", errorMsg)
	w.Flush()
}

// response from OPenrouter(Deepseek models)
func streamOpenRouterResponse(ctx context.Context, w gin.ResponseWriter, req types.ChatRequest, apiKey string, modelConfig models.ModelConfig) {

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	systemPrompt := `
		You are Chat0, an ai assistant that can answer questions and help with tasks.
		Be helpful and provide relevant information
		Be respectful and polite in all interactions.
		Be engaging and maintain a conversational tone.
		Always use LaTeX for mathematical expressions - 
		Inline math must be wrapped in single dollar signs: $content$
		Display math must be wrapped in double dollar signs: $$content$$
		Display math should be placed on its own line, with nothing else on that line.
		Do not nest math delimiters or mix styles.
		Examples:
		- Inline: The equation $E = mc^2$ shows mass-energy equivalence.
		- Display: 
		$$\frac{d}{dx}\sin(x) = \cos(x)$$`
	
	var messages []map[string]string
	
	// Add system instruction as first message
	messages = append(messages, map[string]string{
		"role":    "system",
		"content": systemPrompt,
	})

	// convo history
	for _, msg := range req.Messages {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// req body
	requestBody := map[string]interface{}{
		"model":    modelConfig.ModelID,
		"messages": messages,
		"stream":   true,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		writeErrorChunk(w, "Failed to create request: "+err.Error())
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		writeErrorChunk(w, "Failed to create HTTP request: "+err.Error())
		return
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://chat0.dev") 
	httpReq.Header.Set("X-Title", "Chat0")   
	
	// Make request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		writeErrorChunk(w, "Failed to connect to OpenRouter: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response body for more details
		body, _ := io.ReadAll(resp.Body)
		writeErrorChunk(w, fmt.Sprintf("OpenRouter API error %d: %s", resp.StatusCode, string(body)))
		return
	}

	// Read streaming response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		
		// OpenRouter sends Server-Sent Events format
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			
			// Check for completion
			if data == "[DONE]" {
				fmt.Fprintf(w, "d:{\"finishReason\":\"stop\",\"usage\":{\"promptTokens\":0,\"completionTokens\":0}}\n")
				w.Flush()
				break
			}
			
			// Parse the JSON chunk
			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue // Skip malformed chunks
			}

			// Extract content from OpenRouter format and convert to AI SDK format
			if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if delta, ok := choice["delta"].(map[string]interface{}); ok {
						if content, ok := delta["content"].(string); ok && content != "" {
							// AI SDK data stream format: "0:" prefix for text
							textJSON, _ := json.Marshal(content)
							fmt.Fprintf(w, "0:%s\n", textJSON)
							w.Flush()
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		writeErrorChunk(w, "Streaming error: "+err.Error())
	}

}