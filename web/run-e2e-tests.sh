#!/bin/bash

# SBOM Sentinel E2E Test Runner
# This script starts the backend server and runs end-to-end tests

set -e

echo "🚀 SBOM Sentinel E2E Test Runner"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to cleanup background processes
cleanup() {
    echo -e "\n${YELLOW}🧹 Cleaning up processes...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        echo "✅ Backend server stopped"
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo "✅ Frontend server stopped"
    fi
    # Clean up test database
    rm -f ../test_e2e_sentinel.db 2>/dev/null || true
    echo "✅ Test database cleaned up"
    exit 0
}

# Set up trap to cleanup on script exit
trap cleanup EXIT INT TERM

# Step 1: Start Backend Server
echo -e "${YELLOW}📡 Starting backend server...${NC}"
cd ..
export PORT=8080
export DATABASE_PATH="./test_e2e_sentinel.db"

go run ./cmd/sentinel-server/main.go &
BACKEND_PID=$!

echo "Backend server PID: $BACKEND_PID"

# Wait for backend to be ready
echo -e "${YELLOW}⏳ Waiting for backend server to be ready...${NC}"
for i in {1..30}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend server is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ Backend server failed to start within 30 seconds${NC}"
        exit 1
    fi
    sleep 1
done

# Step 2: Start Frontend Server (in background)
echo -e "${YELLOW}🌐 Starting frontend server...${NC}"
cd web
npm run dev &
FRONTEND_PID=$!

echo "Frontend server PID: $FRONTEND_PID"

# Wait for frontend to be ready
echo -e "${YELLOW}⏳ Waiting for frontend server to be ready...${NC}"
for i in {1..30}; do
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Frontend server is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ Frontend server failed to start within 30 seconds${NC}"
        exit 1
    fi
    sleep 1
done

# Step 3: Run E2E Tests
echo -e "${YELLOW}🎭 Running E2E tests...${NC}"
npx playwright test

# Check test results
if [ $? -eq 0 ]; then
    echo -e "${GREEN}🎉 All E2E tests passed!${NC}"
else
    echo -e "${RED}❌ Some E2E tests failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ E2E test run completed successfully!${NC}"