test:
	go generate ./...
	go test -race -v ./sync ./concurrency

ui:
	go run tests/ui/main.go
