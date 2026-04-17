.PHONY: run-api run-worker build test race cover fmt dev clean all

APP_API=image-processor-api
APP_WORKER=image-processor-worker
PORT=8080

all: fmt test build

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

build:
	mkdir -p bin
	go build -o bin/$(APP_API) ./cmd/api
	go build -o bin/$(APP_WORKER) ./cmd/worker

test:
	go test ./...

race:
	go test -race ./...

cover:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

fmt:
	go fmt ./...
	go vet ./...

dev:
	pkill -f $(APP_API) || true
	go run ./cmd/api & sleep 2 && open http://localhost:$(PORT)

clean:
	rm -rf bin cover.out