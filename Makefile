test:
	go test -v ./...
tidy:
	go mod tidy
build:
	go build -o wallet cmd/wallet/main.go
run:
	go run cmd/wallet/main.go