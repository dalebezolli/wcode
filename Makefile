.PHONY: fmt vet build

build: vet
	go build -o bin/wcode cmd/wcode/main.go

vet: fmt
	go vet ./...

fmt:
	go fmt ./...
