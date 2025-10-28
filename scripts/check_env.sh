#!/bin/bash

# check_env.sh - Check development environment for quotebox project

set -e

echo "================================"
echo "QuoteBox Environment Checker"
echo "================================"
echo ""

MISSING_DEPS=0

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check Go
echo "Checking Go installation..."
if command_exists go; then
    GO_VERSION=$(go version)
    echo "✓ Go is installed: $GO_VERSION"
else
    echo "✗ Go is NOT installed"
    echo ""
    echo "  Installation instructions:"
    echo "  - Ubuntu/Debian: sudo apt update && sudo apt install -y golang-go"
    echo "  - Or download from: https://golang.org/dl/"
    echo "  - For WSL2 on Windows: wsl --install (if not installed), then in Ubuntu:"
    echo "    sudo apt update && sudo apt install -y golang-go"
    echo ""
    MISSING_DEPS=1
fi

# Check Docker
echo ""
echo "Checking Docker installation..."
if command_exists docker; then
    DOCKER_VERSION=$(docker --version)
    echo "✓ Docker is installed: $DOCKER_VERSION"
    
    # Check if Docker daemon is running
    if docker ps >/dev/null 2>&1; then
        echo "✓ Docker daemon is running"
    else
        echo "⚠ Docker is installed but daemon is not running"
        echo "  Start Docker Desktop or run: sudo systemctl start docker"
    fi
else
    echo "✗ Docker is NOT installed"
    echo ""
    echo "  Installation instructions:"
    echo "  - Windows: Download Docker Desktop from https://docs.docker.com/desktop/install/windows-install/"
    echo "  - Linux: https://docs.docker.com/engine/install/"
    echo "  - macOS: https://docs.docker.com/desktop/install/mac-install/"
    echo ""
    MISSING_DEPS=1
fi

# Check Docker Compose
echo ""
echo "Checking Docker Compose..."
if command_exists docker-compose || docker compose version >/dev/null 2>&1; then
    if command_exists docker-compose; then
        COMPOSE_VERSION=$(docker-compose --version)
    else
        COMPOSE_VERSION=$(docker compose version)
    fi
    echo "✓ Docker Compose is available: $COMPOSE_VERSION"
else
    echo "✗ Docker Compose is NOT available"
    echo "  Docker Compose is usually included with Docker Desktop"
    echo "  Or install separately: https://docs.docker.com/compose/install/"
    MISSING_DEPS=1
fi

# Check Git
echo ""
echo "Checking Git installation..."
if command_exists git; then
    GIT_VERSION=$(git --version)
    echo "✓ Git is installed: $GIT_VERSION"
else
    echo "✗ Git is NOT installed"
    echo "  Install with: sudo apt install -y git (Ubuntu/Debian)"
    MISSING_DEPS=1
fi

# Check for .env file
echo ""
echo "Checking for .env file..."
if [ -f .env ]; then
    echo "✓ .env file exists"
    
    # Check for required variables
    if grep -q "OPENROUTER_API_KEY=" .env; then
        if grep -q "OPENROUTER_API_KEY=your_openrouter_api_key_here" .env; then
            echo "⚠ OPENROUTER_API_KEY is set to placeholder value"
            echo "  Please update .env with your actual OpenRouter API key"
        else
            echo "✓ OPENROUTER_API_KEY is configured"
        fi
    else
        echo "⚠ OPENROUTER_API_KEY not found in .env"
    fi
else
    echo "⚠ .env file not found"
    echo "  Copy .env.example to .env and configure your settings:"
    echo "  cp .env.example .env"
fi

# Summary
echo ""
echo "================================"
if [ $MISSING_DEPS -eq 0 ]; then
    echo "✓ All required dependencies are installed!"
    echo ""
    echo "Next steps:"
    echo "1. Make sure .env file is configured with your OPENROUTER_API_KEY"
    echo "2. Run: docker compose up --build"
    echo "3. Access the app at: http://localhost:8080"
    echo "4. Access Grafana at: http://localhost:3000 (admin/admin)"
    echo "5. Access Prometheus at: http://localhost:9090"
else
    echo "✗ Some dependencies are missing. Please install them and run this script again."
    exit 1
fi
echo "================================"
