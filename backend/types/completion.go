package types

// request for generating chat titles
type CompletionRequest struct {
	Prompt 		string 	`json:"prompt" binding:"required"`
	IsTitle 	bool 	`json:"isTitle"`
	MessageID 	string 	`json:"messageId"`
	ThreadID 	string 	`json:"threadId"`
}

// response with generated title
type CompletionResponse struct {
	Title 		string  `json:"title"`
	IsTitle 	bool 	`json:"isTitle"`
	MessageID 	string 	`json:"messageId"`
	ThreadID 	string 	`json:"threadId"`
}

// error response
type ErrorResponse struct {
	Error      string   `json:"error"`
}