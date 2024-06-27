build:
	@go build -o ./bin/blumbot

run: build
	@./bin/blumbot

test:
	@go test -v ./...
