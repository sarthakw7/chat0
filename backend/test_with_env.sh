#!/bin/bash

echo "🧪 Testing Chat0 Go Backend with Environment Variables"
echo "===================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ .env file not found!${NC}"
    echo -e "${YELLOW}💡 Copy .env.example to .env and add your API keys:${NC}"
    echo "   cp .env.example .env"
    echo "   # Then edit .env and add your API keys"
    exit 1
fi

# Load environment variables
source .env

# Check which API keys are available
echo -e "${BLUE}🔑 Checking available API keys...${NC}"
AVAILABLE_KEYS=""

if [ ! -z "$GOOGLE_API_KEY" ]; then
    AVAILABLE_KEYS="$AVAILABLE_KEYS Google"
    echo -e "${GREEN}✅ Google API Key: Available${NC}"
else
    echo -e "${YELLOW}⚠️  Google API Key: Not set${NC}"
fi

if [ ! -z "$OPENAI_API_KEY" ]; then
    AVAILABLE_KEYS="$AVAILABLE_KEYS OpenAI"
    echo -e "${GREEN}✅ OpenAI API Key: Available${NC}"
else
    echo -e "${YELLOW}⚠️  OpenAI API Key: Not set${NC}"
fi

if [ ! -z "$OPENROUTER_API_KEY" ]; then
    AVAILABLE_KEYS="$AVAILABLE_KEYS OpenRouter"
    echo -e "${GREEN}✅ OpenRouter API Key: Available${NC}"
else
    echo -e "${YELLOW}⚠️  OpenRouter API Key: Not set${NC}"
fi

if [ -z "$AVAILABLE_KEYS" ]; then
    echo -e "${RED}❌ No API keys found in .env file!${NC}"
    echo -e "${YELLOW}💡 Add your API keys to .env file:${NC}"
    echo "   GOOGLE_API_KEY=your_key_here"
    echo "   OPENAI_API_KEY=your_key_here"
    echo "   OPENROUTER_API_KEY=your_key_here"
    exit 1
fi

echo ""

# Test health endpoint
echo -e "${YELLOW}1. Testing health endpoint...${NC}"
curl -s http://localhost:8080/api/health | jq .
echo ""

# Test completion endpoint (only if Google key is available)
if [ ! -z "$GOOGLE_API_KEY" ]; then
    echo -e "${YELLOW}2. Testing completion endpoint with Google API key from environment...${NC}"
    curl -s -X POST http://localhost:8080/api/completion \
      -H "Content-Type: application/json" \
      -d '{"prompt":"How to learn programming?","isTitle":true,"messageId":"123","threadId":"456"}' | jq .
    echo ""
else
    echo -e "${YELLOW}2. Skipping completion test (Google API key not available)${NC}"
    echo ""
fi

# Test streaming chat endpoint (only if Google key is available)
if [ ! -z "$GOOGLE_API_KEY" ]; then
    echo -e "${YELLOW}3. Testing streaming chat endpoint with Google API key from environment...${NC}"
    echo "   (This will show real-time streaming)"
    curl -X POST http://localhost:8080/api/chat \
      -H "Content-Type: application/json" \
      -d '{"messages":[{"role":"user","content":"Tell me a short joke"}],"model":"Gemini 1.5 Flash"}'
    echo ""
    echo ""
else
    echo -e "${YELLOW}3. Skipping streaming test (Google API key not available)${NC}"
    echo ""
fi

echo -e "${GREEN}✅ Tests completed!${NC}"
echo ""
echo -e "${BLUE}💡 Available API keys:${AVAILABLE_KEYS}${NC}"
echo -e "${YELLOW}🔧 To add more API keys, edit the .env file${NC}"