package models

// config of AI model
type ModelConfig struct {
	ModelID string `json:"modelId"`
	Provider string `json:"provider"`
	HeaderKey string `json:"headerKey"`
}

// supported maps model names
var SupportedModels = map[string]ModelConfig{
	"Deepseek R1 0528": {
		ModelID:   "deepseek/deepseek-r1-0528:free",
		Provider:  "openrouter",
		HeaderKey: "X-OpenRouter-API-Key",
	},
	"Deepseek V3": {
		ModelID:   "deepseek/deepseek-chat-v3-0324:free",
		Provider:  "openrouter",
		HeaderKey: "X-OpenRouter-API-Key",
	},
	"Gemini 2.5 Pro": {
		ModelID:   "gemini-2.5-pro",
		Provider:  "google",
		HeaderKey: "X-Google-API-Key",
	},
	"Gemini 2.5 Flash": {
		ModelID:   "gemini-2.5-flash",
		Provider:  "google",
		HeaderKey: "X-Google-API-Key",
	},
	"Gemini 1.5 Flash": {
		ModelID:   "gemini-1.5-flash",
		Provider:  "google",
		HeaderKey: "X-Google-API-Key",
	},
	"GPT-4o": {
		ModelID:   "gpt-4o",
		Provider:  "openai",
		HeaderKey: "X-OpenAI-API-Key",
	},
	"GPT-4o-mini": {
		ModelID:   "gpt-4o-mini",
		Provider:  "openai",
		HeaderKey: "X-OpenAI-API-Key",
	},
}

func GetModelConfig(modelName string) (ModelConfig, bool) {
	config, exists := SupportedModels[modelName]
	return config, exists
}

func GetSupportedModelNames() []string {
	names := make([]string, 0, len(SupportedModels))
	for name := range SupportedModels {
		names = append(names, name)
	}
	return names
}