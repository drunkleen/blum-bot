build:
	@go build -o ./bin/blumbot-linux

run: build
	@./bin/blumbot

test:
	@go test -v ./...
