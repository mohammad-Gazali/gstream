run: build
	@./bin/main

build:
	@go build -o bin/main

test-publish:
	@go run cmd/testpublish/main.go

test-consumer:
	@go run cmd/testconsumer/main.go

test:
	@go test ./... -v