#!/usr/bin/env bash
set -euo pipefail

echo "========================================="
echo "  Mock Starket - Development Setup"
echo "========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_dependency() {
    if command -v "$1" &> /dev/null; then
        echo -e "  ${GREEN}✓${NC} $1 found"
        return 0
    else
        echo -e "  ${RED}✗${NC} $1 not found"
        return 1
    fi
}

echo "Checking dependencies..."
MISSING=0
check_dependency "go" || MISSING=1
check_dependency "node" || MISSING=1
check_dependency "npm" || MISSING=1
check_dependency "docker" || MISSING=1
check_dependency "docker compose" 2>/dev/null || check_dependency "docker-compose" || MISSING=1

if [ $MISSING -eq 1 ]; then
    echo ""
    echo -e "${YELLOW}Some dependencies are missing. Install them before continuing.${NC}"
    echo ""
fi

echo ""
echo "Setting up backend..."
cd "$(dirname "$0")/../backend"
if [ -f "go.mod" ]; then
    go mod download
    echo -e "  ${GREEN}✓${NC} Go dependencies installed"
else
    echo -e "  ${YELLOW}⚠${NC} No go.mod found (backend not initialized yet)"
fi

echo ""
echo "Setting up web frontend..."
cd "$(dirname "$0")/../web"
if [ -f "package.json" ]; then
    npm ci
    echo -e "  ${GREEN}✓${NC} Node dependencies installed"
else
    echo -e "  ${YELLOW}⚠${NC} No package.json found (web not initialized yet)"
fi

echo ""
echo "Setting up environment..."
cd "$(dirname "$0")/../deploy"
if [ ! -f ".env" ]; then
    cp .env.example .env
    echo -e "  ${GREEN}✓${NC} Created deploy/.env from template"
else
    echo -e "  ${YELLOW}⚠${NC} deploy/.env already exists, skipping"
fi

echo ""
echo "========================================="
echo -e "  ${GREEN}Setup complete!${NC}"
echo ""
echo "  Quick start:"
echo "    make dev          # Start full stack"
echo "    make dev-backend  # Backend + DB only"
echo "    make test-all     # Run all tests"
echo "========================================="
