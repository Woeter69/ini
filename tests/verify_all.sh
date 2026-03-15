#!/bin/bash
set -e

# verify_all.sh — Run all INI integration tests

echo "🚀 Building ini..."
go build -o ini main.go

echo "🧪 Running internal integration tests..."
# We run with -v to see the progress of all 39+ languages
go test ./tests/... -v

echo "✅ All tests passed!"
rm ini
