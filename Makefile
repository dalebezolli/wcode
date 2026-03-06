.PHONY: fmt vet build

build: vet
	go build

vet: fmt
	go vet ./...

fmt:
	go fmt ./...
