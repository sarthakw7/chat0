// API Configuration for Go Backend Integration
export const API_CONFIG = {
  // Go backend URL - change this for production deployment
  BASE_URL:
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
      : "http://localhost:8080",

  ENDPOINTS: {
    CHAT: "/api/chat",
    COMPLETION: "/api/completion",
    HEALTH: "/api/health",
  },
} as const;

// Helper function to get full API URL
export const getApiUrl = (endpoint: string) => {
  return `${API_CONFIG.BASE_URL}${endpoint}`;
};

// Export individual endpoints for convenience
export const API_ENDPOINTS = {
  CHAT: getApiUrl(API_CONFIG.ENDPOINTS.CHAT),
  COMPLETION: getApiUrl(API_CONFIG.ENDPOINTS.COMPLETION),
  HEALTH: getApiUrl(API_CONFIG.ENDPOINTS.HEALTH),
} as const;

// Debug helper
export const logApiConfig = () => {
  console.log("🔗 API Configuration:", {
    baseUrl: API_CONFIG.BASE_URL,
    endpoints: API_ENDPOINTS,
    environment: process.env.NODE_ENV,
  });
};
