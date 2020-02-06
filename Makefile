test:
	go generate
	go test -race -v ./sync

ui:
	go run ui_tester.go
