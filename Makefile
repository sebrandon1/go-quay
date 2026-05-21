APP_NAME=go-quay

# Integration testing:
# Run 'make integration-test' to run all CLI tests
# Set QUAY_TOKEN and QUAY_ORG environment variables for API integration tests
# Example: QUAY_TOKEN=your_token QUAY_ORG=your_org make integration-test

vet:
	go vet ./...

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(APP_NAME)

lint:
	golangci-lint run ./...

test:
	go test ./... -v

integration-test: build
	@./scripts/integration-test.sh

clean:
	rm -f $(APP_NAME)

check-swagger-alignment:
	@echo "Checking API alignment with Quay.io Swagger spec..."
	@go run ./scripts/check-swagger-alignment.go \
		--swagger-url="https://quay.io/api/v1/discovery" \
		--lib-path="./lib" \
		--base-url-var="QuayURL"

.PHONY: vet build lint test integration-test clean check-swagger-alignment
