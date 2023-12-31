.PHONY: tidy fmt test build

tidy:
	go mod tidy

lint: tidy
	golangci-lint run

fmt:
	go fmt ./...

test: tidy
	go test ./... -cover -p 1

build: tidy
	go build ./...