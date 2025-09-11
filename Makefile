BINARY_NAME=jti-super-app
DOCKER_IMAGE_NAME=jti-super-app
DOCKER_CONTAINER_NAME=jti-app-container

all: build

build:
	@echo "Building application binary..."
	go build -o $(BINARY_NAME) main.go
	@echo "Build complete: $(BINARY_NAME)"

run:
	@echo "Running application locally..."
	go run main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(BINARY_NAME)

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .
	@echo "Docker image '$(DOCKER_IMAGE_NAME)' built successfully."

docker-run:
	@echo "Running Docker container '$(DOCKER_CONTAINER_NAME)'..."
	docker run -d -p 8000:8000 --env-file .env --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

docker-stop:
	@echo "Stopping Docker container '$(DOCKER_CONTAINER_NAME)'..."
	docker stop $(DOCKER_CONTAINER_NAME) && docker rm $(DOCKER_CONTAINER_NAME)

docker-logs:
	@echo "Showing logs for '$(DOCKER_CONTAINER_NAME)'. Press Ctrl+C to exit."
	docker logs -f $(DOCKER_CONTAINER_NAME)

.PHONY: all build run test clean docker-build docker-run docker-stop docker-logs