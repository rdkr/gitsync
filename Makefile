test:
	go generate
	go test -race -v ./sync ./concurrency

ui:
	go run ui_tester.go
