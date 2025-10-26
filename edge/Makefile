.PHONY: build run clean deps

BINARY=videodisco-edge

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	go clean

deps:
	go mod download
	go mod tidy

