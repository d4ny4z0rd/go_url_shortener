build:
	@go build -o bin/url_shortener_1 *.go

test:
	@go test -v ./...

run: build
	@./bin/url_shortener_1