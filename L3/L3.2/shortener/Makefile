# =====================
# BASIC
# =====================

run:
	go run ./cmd

build:
	go build -o app ./cmd

# =====================
# DOCKER
# =====================

docker-run:
	docker-compose up --build

docker-down:
	docker-compose down

logs:
	docker-compose logs -f app

# =====================
# TESTS
# =====================

test:
	go test ./...

cover:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

# =====================
# LINT & FORMAT
# =====================

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

# =====================
# DEV
# =====================

dev:
	docker-compose up --build &
	sleep 3
	open http://localhost:8080

# =====================
# UTILS
# =====================

clean:
	rm -f app
	rm -f cover.out
	rm -f data.db