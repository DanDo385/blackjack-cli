.PHONY: build run test fmt clean

build:
	@echo "Building blackjack..."
	@go build -o ./bin/blackjack ./cmd/blackjack

run:
	@echo "Running blackjack..."
	@go run ./cmd/blackjack

test:
	@echo "Running tests..."
	@go test ./... -race -count=1

fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
