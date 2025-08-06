#!/bin/bash

echo "🧪 Testing Chat0 Go Backend with Real AI Integration"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test health endpoint
echo -e "${YELLOW}1. Testing health endpoint...${NC}"
curl -s http://localhost:8080/api/health | jq .
echo ""

# Test completion endpoint with fake key (should show error)
echo -e "${YELLOW}2. Testing completion endpoint (with fake key - should show error)...${NC}"
curl -s -X POST http://localhost:8080/api/completion \
  -H "Content-Type: application/json" \
  -H "X-Google-API-Key: fake-key" \
  -d '{"prompt":"Hello world","isTitle":true,"messageId":"123","threadId":"456"}' | jq .
echo ""

# Test model validation
echo -e "${YELLOW}3. Testing model validation (unsupported model)...${NC}"
curl -s -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role":"user","content":"Hello"}],"model":"Unknown Model"}' | jq .
echo ""

# Test API key validation
echo -e "${YELLOW}4. Testing API key validation (missing key)...${NC}"
curl -s -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role":"user","content":"Hello"}],"model":"Gemini 2.5 Flash"}' | jq .
echo ""

# Test streaming chat endpoint with fake key (should show streaming error)
echo -e "${YELLOW}5. Testing streaming chat endpoint (with fake key - should show streaming error)...${NC}"
echo "   (This will show real-time streaming error)"
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -H "X-Google-API-Key: fake-key" \
  -d '{"messages":[{"role":"user","content":"Tell me a joke"}],"model":"Gemini 2.5 Flash"}'
echo ""
echo ""

# Show supported models
echo -e "${YELLOW}6. Supported models:${NC}"
echo "   - Gemini 2.5 Pro (Google)"
echo "   - Gemini 2.5 Flash (Google)"
echo "   - GPT-4o (OpenAI)"
echo "   - GPT-4.1-mini (OpenAI)"
echo "   - Deepseek R1 0528 (OpenRouter)"
echo "   - Deepseek V3 (OpenRouter)"
echo ""

echo -e "${GREEN}✅ All tests completed!${NC}"
echo ""
echo -e "${YELLOW}💡 To test with real API keys:${NC}"
echo "   export GOOGLE_API_KEY='your-real-key'"
echo "   curl -X POST http://localhost:8080/api/completion \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -H \"X-Google-API-Key: \$GOOGLE_API_KEY\" \\"
echo "     -d '{\"prompt\":\"Hello world\",\"isTitle\":true}'"