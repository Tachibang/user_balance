.PHONY: build run run-bin

build:
	go build -o /bin/app ./main.go

run:
	go run cmd/main.go

run-bin:
	./bin/main

test:
	go test -v ./...