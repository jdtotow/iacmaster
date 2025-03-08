# Variables
SERVICE1_NAME := system
SERVICE2_NAME := service
SERVICE3_NAME := runner 
SERVICE1_CMD := ./cmd/$(SERVICE1_NAME)
SERVICE2_CMD := ./cmd/$(SERVICE2_NAME)
SERVIVE3_CMD := ./cmd/$(SERVICE3_NAME)
SERVICE1_BINARY := bin/$(SERVICE1_NAME)
SERVICE2_BINARY := bin/$(SERVICE2_NAME)
SERVICE3_BINARY := bin/$(SERVICE3_NAME)
DOCKERFILE1 := system.dockerfile
DOCKERFILE2 := service.dockerfile
DOCKERFILE3 := runner.dockerfile
IMAGE1_NAME := iacmaster_$(SERVICE1_NAME)
IMAGE2_NAME := iacmaster_$(SERVICE2_NAME)
IMAGE3_NAME := iacmaster_$(SERVICE3_NAME)

# Default target
.PHONY: all
all: build

# Build binaries
.PHONY: build
build: build-system build-service 

.PHONY: build-system
build-system:
	@echo "Building $(SERVICE1_NAME)..."
	@go build -o $(SERVICE1_BINARY) $(SERVICE1_CMD)

.PHONY: build-service
build-service:
	@echo "Building $(SERVICE2_NAME)..."
	@go build -o $(SERVICE2_BINARY) $(SERVICE2_CMD)

# Create Docker images
.PHONY: docker
docker: docker-system docker-service docker-runner 

.PHONY: docker-system
docker-system:
	@echo "Building Docker image for $(SERVICE1_NAME)..."
	@docker build -f $(DOCKERFILE1) -t $(IMAGE1_NAME) .

.PHONY: docker-service
docker-service:
	@echo "Building Docker image for $(SERVICE2_NAME)..."
	@docker build -f $(DOCKERFILE2) -t $(IMAGE2_NAME) .

.PHONY: docker-runner
docker-runner:
	@echo "Building Docker image for $(SERVICE3_NAME)..."
	@docker build -f $(DOCKERFILE3) -t $(IMAGE3_NAME) .

# Run services without Docker
.PHONY: run
run: run-system run-service 

.PHONY: run-system
run-system: build-system
	@echo "Running $(SERVICE1_NAME) without Docker..."
	@./$(SERVICE1_BINARY)

.PHONY: run-service
run-service: build-service
	@echo "Running $(SERVICE2_NAME) without Docker..."
	@./$(SERVICE2_BINARY)

# Clean up binaries
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(SERVICE1_BINARY) $(SERVICE2_BINARY) $(SERVICE3_BINARY)
