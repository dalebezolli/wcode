.PHONY: fmt vet build

build: vet
	go build -o bin/wcode ./cmd/wcode

vet: fmt
	go vet ./...

fmt:
	go fmt ./...
