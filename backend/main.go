package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sarthakw7/chat0-backend/handlers"
)

func main() {
	// load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	// gin router
	r := gin.Default()

	// gets cors origins from environment
	allowedOrigins := []string{"http://localhost:3000"}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	}

	// cors middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, 
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"}, 
		AllowCredentials: true,
	}))
		

	// health check route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "Go backend is running!",
			"status": "ok",
		})
	})

	// API route
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H {
				"status": "ok",
				"service": "chat0-backend",
			})
		})

		// completion endpoint
		api.POST("/completion", handlers.HandleCompletion)

		// streaming endpoint
		api.POST("/chat", handlers.HandleChat)
	}

	// get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Chat0 Go Backend starting on port %s", port)
	log.Printf("📋 Environment: %s", getEnv("ENVIRONMENT", "development"))
	log.Printf("🔑 API Keys loaded: %s", getLoadedKeys())

	r.Run(":" + port)
}

// get env with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// show which API keys are loaded
func getLoadedKeys() string {
	var loaded []string
	
	if os.Getenv("GOOGLE_API_KEY") != "" {
		loaded = append(loaded, "Google")
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		loaded = append(loaded, "OpenAI")
	}
	if os.Getenv("OPENROUTER_API_KEY") != "" {
		loaded = append(loaded, "OpenRouter")
	}
	
	if len(loaded) == 0 {
		return "None (will use headers)"
	}
	
	return strings.Join(loaded, ", ")
}