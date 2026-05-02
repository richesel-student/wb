APP_NAME=warehouse-app

run:
	docker compose up --build

down:
	docker compose down -v

logs:
	docker compose logs -f

test:
	go test -v ./...

deps:
	go mod tidy
fmt:
	go fmt ./...
	go vet ./...

build:
	go build -o $(APP_NAME)

clean:
	rm -f $(APP_NAME)


restart: down run