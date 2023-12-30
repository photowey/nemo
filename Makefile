.PHONY: tidy fmt test build

tidy:
	go mod tidy

lint: tidy
	golangci-lint run

fmt:
	go fmt ./...

test:
	go test ./... -cover -p 1

build:
	go build ./...