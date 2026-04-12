.PHONY: run build test race cover fmt dev clean all

APP=comments-service
PORT=8080

all: fmt test build

run:
	go run ./cmd

build:
	go build -o $(APP) ./cmd

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
	pkill -f $(APP) || true
	go run ./cmd & sleep 2 && open http://localhost:$(PORT)

clean:
	rm -f $(APP) cover.out