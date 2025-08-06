package types

// single message in conversation
type ChatMessage struct {
	Role string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// incoming chat req
type ChatRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required"`
	Model string `json:"model" binding:"required"`
}

// streaming response chunk
type StreamChunk struct {
	ID string `json:"id"`
	Object string `json:"object"`
	Created int64 `json:"created"`
	Model string `json:"model"`
	Choices []Choice `json:"choices"`
}

// choice in streaming response
type Choice struct {
	Index int `json:"index"`
	Delta Delta `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// incremanral content is streaming
type Delta struct {
	Content string `json:"content"`
	Role string `json:"role,omitempty"`
}