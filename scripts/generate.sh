#!/bin/bash

# Solast-go Parser Generation Script
# This script generates Go parser code from Solidity ANTLR grammar

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
GRAMMAR_DIR="$PROJECT_ROOT/grammar"
GEN_DIR="$PROJECT_ROOT/internal/gen"
ANTLR_VERSION="4.13.1"
ANTLR_JAR="antlr-${ANTLR_VERSION}-complete.jar"
ANTLR_URL="https://www.antlr.org/download/${ANTLR_JAR}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Solast-go Parser Generator${NC}"
echo "================================"

# Check for Java
if ! command -v java &> /dev/null; then
    echo -e "${RED}Error: Java is required but not installed.${NC}"
    echo "Please install Java Runtime Environment (JRE) to generate the parser."
    exit 1
fi

# Download ANTLR if not present
if [ ! -f "$PROJECT_ROOT/$ANTLR_JAR" ]; then
    echo -e "${YELLOW}Downloading ANTLR ${ANTLR_VERSION}...${NC}"
    curl -o "$PROJECT_ROOT/$ANTLR_JAR" "$ANTLR_URL"
fi

# Check if grammar files exist
if [ ! -f "$GRAMMAR_DIR/SolidityLexer.g4" ] || [ ! -f "$GRAMMAR_DIR/SolidityParser.g4" ]; then
    echo -e "${RED}Error: Grammar files not found in $GRAMMAR_DIR${NC}"
    echo "Expected files: SolidityLexer.g4, SolidityParser.g4"
    exit 1
fi

# Create output directory
mkdir -p "$GEN_DIR"

# Generate parser
echo -e "${YELLOW}Generating Go parser from grammar...${NC}"
java -jar "$PROJECT_ROOT/$ANTLR_JAR" \
    -Dlanguage=Go \
    -visitor \
    -package gen \
    -o "$GEN_DIR" \
    "$GRAMMAR_DIR/SolidityLexer.g4" \
    "$GRAMMAR_DIR/SolidityParser.g4"

echo -e "${GREEN}Parser generated successfully!${NC}"
echo "Output directory: $GEN_DIR"

