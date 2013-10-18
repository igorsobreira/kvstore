test:
	go test -race ./...

lint:
	@golint `find . -name "*.go"`

fmt:
	@go fmt ./...
