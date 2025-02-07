APP_NAME=go-quay

vet:
	go vet ./...

build:
	go build -o $(APP_NAME)

lint:
	golangci-lint run ./...
