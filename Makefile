
.PHONY: build
build: 
	@go generate ./...
	@go build ./cmd/server
