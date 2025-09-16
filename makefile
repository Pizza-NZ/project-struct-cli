# Name of your application binary
BINARY_NAME=project-struct-cli

# Default target executed when you just run `make`
default: build

## --------------------------------------
## Build Targets
## --------------------------------------

# Build the application
build:
	@echo "Building..."
	@go build -o ./out/$(BINARY_NAME) .

# Run the application
run: build
	@echo "Running..."
	@./out/$(BINARY_NAME)

# Run tests
test:
	@echo "Testing..."
	@go test -v ./...

# Remove the binary
clean:
	@echo "Cleaning..."
	@rm -f ./out/$(BINARY_NAME)

.PHONY: default build run test clean