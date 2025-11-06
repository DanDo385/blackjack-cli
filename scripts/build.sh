#!/usr/bin/env bash
set -euo pipefail

echo "Building blackjack CLI..."
go build -o ./bin/blackjack ./cmd/blackjack
echo "Build complete! Binary at ./bin/blackjack"
