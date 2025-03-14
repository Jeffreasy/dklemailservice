#!/bin/bash

# Script om tests uit te voeren met CGO enabled voor SQLite ondersteuning

# Zorg ervoor dat CGO enabled is
export CGO_ENABLED=1

# Ga naar de project root directory
cd "$(dirname "$0")/.."

# Voer alle tests uit
echo "Running tests with CGO_ENABLED=1..."
go test ./tests/... -v

# Voer tests uit met coverage
if [ "$1" == "--coverage" ]; then
    echo "Running tests with coverage..."
    go test ./tests/... -coverprofile=coverage -v
    go tool cover -html=coverage -o coverage.html
    echo "Coverage report generated at coverage.html"
fi 