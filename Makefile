# run goimports on all go files except *pb.go
# buf lint is run in proto-gen; so no need to run here
lint:
	find . -type f -name "*.go" ! -name "*pb.go" -exec goimports -w {} \;
	go fmt ./...
	pre-commit run --all-files
	golangci-lint run ./...

test:
	go mod tidy
	go test ./... -timeout 10m

check: lint test
