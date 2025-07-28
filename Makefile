APP_NAME=go-quay

# Integration testing:
# Run 'make integration-test' to run all CLI tests
# Set QUAY_TOKEN and QUAY_ORG environment variables for API integration tests
# Example: QUAY_TOKEN=your_token QUAY_ORG=your_org make integration-test

vet:
	go vet ./...

build:
	go build -o $(APP_NAME)

lint:
	golangci-lint run ./...

test:
	go test ./... -v

integration-test: build
	@./scripts/integration-test.sh

clean:
	rm -f $(APP_NAME)

.PHONY: vet build lint test integration-test clean
