lint:
	golangci-lint run

vet:
	go vet ./...

fmt:
	go fmt ./...

run:
	docker-compose up --build

logs:
	docker-compose logs -f worker

open:

	open http://localhost:8080/ui
dev:
	docker-compose up --build &
	sleep 5
	open http://localhost:8080/ui