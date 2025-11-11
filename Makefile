.PHONY: build run test fmt clean

build:
	@./scripts/build.sh

run:
	@./scripts/run.sh

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
