build:
	@go build -o bin/ecom_backend

test:
	@go test -v ./...

run:
	@go run cmd/main.go