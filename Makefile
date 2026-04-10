.PHONY: build run clean

BINARY_NAME=gemma-agent

build:
	go build -o $(BINARY_NAME) main.go

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

deps:
	go mod tidy
	go mod download
