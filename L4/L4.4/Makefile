# Makefile for GC Monitor Project

.PHONY: run build test clean


all: build


run:
	go run main.go


build:
	go build -o gc-monitor main.go

# Test
test:
	go test ./...


clean:
	rm -f gc-monitor


deps:
	go mod tidy
	go mod download