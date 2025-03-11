# Variables
SERVICE1_NAME := system
SERVICE3_NAME := runner 
SERVICE1_CMD := ./cmd/$(SERVICE1_NAME)
SERVIVE3_CMD := ./cmd/$(SERVICE3_NAME)
SERVICE1_BINARY := bin/$(SERVICE1_NAME)
SERVICE3_BINARY := bin/$(SERVICE3_NAME)
DOCKERFILE1 := system.dockerfile
DOCKERFILE3 := runner.dockerfile
IMAGE1_NAME := iacmaster_$(SERVICE1_NAME)
IMAGE3_NAME := iacmaster_$(SERVICE3_NAME)
PROTO_SRC_DIR :=./pkg/msg
PROTO_DST_DIR :=./pkg/protos

# Default target
.PHONY: all
all: build

# Build binaries
.PHONY: build
build: proto build-system 

.PHONY: build-system
build-system:
	@echo "Building $(SERVICE1_NAME)..."
	@go build -o $(SERVICE1_BINARY) $(SERVICE1_CMD)

# Create Docker images
.PHONY: docker
docker: docker-system docker-runner 

.PHONY: docker-system
docker-system:
	@echo "Building Docker image for $(SERVICE1_NAME)..."
	@docker build -f $(DOCKERFILE1) -t $(IMAGE1_NAME) .

.PHONY: docker-runner
docker-runner:
	@echo "Building Docker image for $(SERVICE3_NAME)..."
	@docker build -f $(DOCKERFILE3) -t $(IMAGE3_NAME) .

# Run services without Docker
.PHONY: run
run: run-system 

.PHONY: run-system
run-system: build-system
	@echo "Running $(SERVICE1_NAME) without Docker..."
	@./$(SERVICE1_BINARY)

# Clean up binaries
.PHONY: proto
proto:
	@echo "Compiling proto files..."
	@protoc -I=$(PROTO_SRC_DIR) --go_out=$(PROTO_DST_DIR) $(PROTO_SRC_DIR)/*.proto

.PHONY: clean
clean:
	@echo "Cleaning up ..."
	@rm -f $(SERVICE1_BINARY) $(SERVICE2_BINARY) $(SERVICE3_BINARY)
	@echo "End cleaning!"
