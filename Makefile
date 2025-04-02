.PHONY: all build build-server build-client run test lint docker-build-server docker-build-client docker-run-server docker-run-client dc

# Build both binaries (server and client)
all: build

build: build-server build-client

build-server:
	@echo "Building server..."
	go build -o wordofwisdom ./cmd/server/main.go

build-client:
	@echo "Building client..."
	go build -o wordofwisdom-client ./cmd/client/main.go

# Run the server locally
run: build-server
	@echo "Running server..."
	./wordofwisdom

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Linting (uses go vet; you can add other tools)
lint:
	@echo "Running lint checks..."
	go vet ./...

# Build Docker image for the server
docker-build-server:
	@echo "Building Docker image for server..."
	docker build -f Dockerfile.server -t wordofwisdom-server .

# Build Docker image for the client
docker-build-client:
	@echo "Building Docker image for client..."
	docker build -f Dockerfile.client -t wordofwisdom-client .

# Run Docker container for the server (port 9000)
docker-run-server:
	@echo "Running Docker container for server..."
	docker run -p 9000:9000 wordofwisdom-server

# Run Docker container for the client
docker-run-client:
	@echo "Running Docker container for client..."
	docker run wordofwisdom-client

# Start services using docker-compose
dc:
	@echo "Starting services with docker-compose..."
	docker-compose up
