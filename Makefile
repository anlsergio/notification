.PHONY: build test cmd

# Where the compiled binary will be created at.
BUILD_DIR = ./build
# The name of the compiled binary.
BINARY_NAME = server
# The package where the binary will bootstrap the application.
MAIN_PACKAGE = ./cmd/http
# The relative path to the compiled binary
FULL_BINARY_NAME = $(BUILD_DIR)/$(BINARY_NAME)

# Compile a binary while running the test suite.
all: test build

# Compile the binary
build:
	go build -o $(FULL_BINARY_NAME) $(MAIN_PACKAGE)

# Runs the application locally
run: build
	$(FULL_BINARY_NAME)

# Housekeeping: clean up Go cache,
# temp files and build artifacts.
clean:
	go clean
	rm -f $(FULL_BINARY_NAME)

# Tidy up go modules
tidy:
	go mod tidy

# Run the test suite
test:
	go test -v ./...

# Run the test suite generating a coverage file
# to be interpreted using the "show-doc"
test-cov:
	./scripts/generate-test-cov.sh

# Parse the coverage file generated by "test-cov"
# and render an HTML page.
show-cov:
	go tool cover -html=/tmp/test.cov

# Spin up a Docker container running the application.
docker-up:
	docker compose up -d

# Update the current Docker container with the changes
# made to the project.
docker-update:
	docker compose up --build -d

# Simply build a Docker image.
docker-build:
	docker build -t notification:latest .

# Stops the container execution.
docker-stop:
	docker stop notification

# Removes the container.
docker-rm: docker-stop
	docker rm notification

# Generate API Swagger documentation
api-docs:
	swag init --parseInternal=true \
      --generatedTime=true \
      --dir=cmd/http \
      --parseDependency=true \
      --output=api
