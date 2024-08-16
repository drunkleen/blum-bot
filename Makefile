PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := blumbot

.PHONY: build run test clean

# Build for all platforms and architectures
build:
	@mkdir -p ./bin
	@for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCHS); do \
			output_name=./bin/$(BINARY)-$$GOOS-$$GOARCH; \
			if [ $$GOOS = "windows" ]; then \
				output_name=$$output_name.exe; \
			fi; \
			echo "Building for $$GOOS/$$GOARCH..."; \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$output_name .; \
		done \
	done

# Run the built binary (only works for the current OS/ARCH)
run: build
	@./bin/$(BINARY)

test:
	@go test -v ./...

clean:
	@rm -rf ./bin
