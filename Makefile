test:
	go generate
	go test -race -v ./...

ui:
	go run ui_tester.go
