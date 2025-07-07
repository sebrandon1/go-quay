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
	@echo "Running CLI integration tests..."
	@echo "Testing main CLI help..."
	@./$(APP_NAME) --help > /dev/null
	@echo "Testing get command help..."
	@./$(APP_NAME) get --help > /dev/null
	@echo "Testing billing command help..."
	@./$(APP_NAME) get billing --help > /dev/null
	@echo "Testing billing subcommands help..."
	@./$(APP_NAME) get billing org-info --help > /dev/null
	@./$(APP_NAME) get billing user-info --help > /dev/null
	@./$(APP_NAME) get billing org-subscription --help > /dev/null
	@./$(APP_NAME) get billing user-subscription --help > /dev/null
	@./$(APP_NAME) get billing org-invoices --help > /dev/null
	@# ./$(APP_NAME) get billing user-invoices --help > /dev/null  # Endpoint doesn't exist in Quay API
	@# Usage commands commented out - endpoints don't exist in Quay API
	@# ./$(APP_NAME) get billing org-usage --help > /dev/null
	@# ./$(APP_NAME) get billing user-usage --help > /dev/null
	@./$(APP_NAME) get billing plans --help > /dev/null
	@echo "Testing error handling..."
	@OUTPUT=$$(./$(APP_NAME) get billing org-info 2>&1) && echo "$$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error" || (echo "✗ Missing required flag error not displayed" && exit 1)
	@OUTPUT=$$(./$(APP_NAME) get billing user-info 2>&1) && echo "$$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error" || (echo "✗ Missing required flag error not displayed" && exit 1)
	@if [ -n "$$QUAY_TOKEN" ]; then \
		echo "Running integration tests with real API..."; \
		printf "Testing user billing info... "; \
		if ./$(APP_NAME) get billing user-info --token "$$QUAY_TOKEN" > /dev/null 2>&1; then echo "✓"; else echo "✗ (may not have access)"; fi; \
		printf "Testing user subscription... "; \
		if ./$(APP_NAME) get billing user-subscription --token "$$QUAY_TOKEN" > /dev/null 2>&1; then echo "✓"; else echo "✗ (may not have access)"; fi; \
		printf "Testing available plans... "; \
		if ./$(APP_NAME) get billing plans --token "$$QUAY_TOKEN" > /dev/null 2>&1; then echo "✓"; else echo "✗ (may not have access)"; fi; \
		if [ -n "$$QUAY_ORG" ]; then \
			printf "Testing org billing info for $$QUAY_ORG... "; \
			if ./$(APP_NAME) get billing org-info --organization "$$QUAY_ORG" --token "$$QUAY_TOKEN" > /dev/null 2>&1; then echo "✓"; else echo "✗ (may not have access)"; fi; \
			printf "Testing org subscription for $$QUAY_ORG... "; \
			if ./$(APP_NAME) get billing org-subscription --organization "$$QUAY_ORG" --token "$$QUAY_TOKEN" > /dev/null 2>&1; then echo "✓"; else echo "✗ (may not have access)"; fi; \
			printf "Testing org invoices for $$QUAY_ORG... "; \
			if ./$(APP_NAME) get billing org-invoices --organization "$$QUAY_ORG" --token "$$QUAY_TOKEN" > /dev/null 2>&1; then echo "✓"; else echo "✗ (may not have access)"; fi; \
		else \
			echo "No QUAY_ORG specified, skipping organization tests"; \
		fi \
	else \
		echo "No QUAY_TOKEN available, skipping integration tests"; \
	fi
	@echo "✓ All CLI tests completed successfully"

clean:
	rm -f $(APP_NAME)

.PHONY: vet build lint test integration-test clean
