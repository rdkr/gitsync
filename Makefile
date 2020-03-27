.PHONY: test ui

test: mocks/mock_git.go
	go test -race -v ./sync ./concurrency
	staticcheck ./sync ./concurrency ./tests/ui

ui:
	go run tests/ui/main.go

mocks/mock_git.go: sync/git.go
	go generate ./...
