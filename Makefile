# run goimports on all go files except *pb.go
lint:
	find . -type f -name "*.go" ! -name "*pb.go" -exec goimports -w {} \;
	go fmt ./...
	pre-commit run --all-files
	golangci-lint run ./...

test:
	go mod tidy
	go test ./... -timeout 10m

check: lint test

# Database commands
db-up:
	docker-compose up -d db

db-down:
	docker-compose down

# Docker commands
build:
	docker build -t knovel-api .

run: build
	docker-compose up -d db
	docker run --network=host -p 8080:8080 -e DATABASE_HOST=localhost -e DATABASE_PORT=5433 -e DATABASE_USER=postgres -e DATABASE_PASSWORD=password -e DATABASE_NAME=knovel knovel-api

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Helper command to start the project from scratch
setup: docker-down
	docker-compose up -d db
	go run cmd/migrate/main.go
