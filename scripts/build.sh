#!/usr/bin/env bash
# scripts/build.sh
set -euo pipefail

echo "Building blackjack CLI..."
go build -o ./bin/blackjack ./cmd/blackjack
echo "Build complete! Binary at ./bin/blackjack"
