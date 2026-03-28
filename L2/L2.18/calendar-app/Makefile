APP=calendar
PORT=8080

.PHONY: setup build run start test-api test vet lint check clean

# права на скрипты
setup:
	chmod 755 scripts/*.sh

# сборка
build: setup
	go build -o $(APP) ./cmd

# запуск
run: setup
	PORT=$(PORT) go run ./cmd/main.go

# запуск бинаря
start: build
	PORT=$(PORT) ./$(APP)

# API тесты (curl)
test-api: setup

	./scripts/create_event.sh
	@sleep 1
	./scripts/get_day.sh
	@sleep 1
	./scripts/delete_event.sh

# unit тесты Go
test:
	go test ./... -v

# go vet
vet:
	go vet ./...

# lint (нужен golangci-lint)
lint:
	golangci-lint run || echo "⚠️ golangci-lint not installed"

# полный чек (как в CI)
check: vet lint test


# очистка
clean:
	rm -f $(APP)