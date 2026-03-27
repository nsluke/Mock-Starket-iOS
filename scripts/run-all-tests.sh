#!/usr/bin/env bash
set -euo pipefail

echo "========================================="
echo "  Mock Starket - Full Test Suite"
echo "========================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

FAILED=0

echo ""
echo "--- Backend (Go) ---"
if cd "$(dirname "$0")/../backend" && [ -f "go.mod" ]; then
    if go test ./... -race -count=1; then
        echo -e "${GREEN}✓ Backend tests passed${NC}"
    else
        echo -e "${RED}✗ Backend tests failed${NC}"
        FAILED=1
    fi
else
    echo "⚠ Backend not initialized, skipping"
fi

echo ""
echo "--- Web (Next.js) ---"
if cd "$(dirname "$0")/../web" && [ -f "package.json" ]; then
    if npm run test -- --run; then
        echo -e "${GREEN}✓ Web tests passed${NC}"
    else
        echo -e "${RED}✗ Web tests failed${NC}"
        FAILED=1
    fi
else
    echo "⚠ Web not initialized, skipping"
fi

echo ""
echo "========================================="
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All test suites passed!${NC}"
else
    echo -e "${RED}Some test suites failed.${NC}"
    exit 1
fi
