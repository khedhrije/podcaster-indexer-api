# Load environment variables if the file exists
include .env
export


default: build

build:
	go build ./...

# Target to run the main Go application
.PHONY: run
run:
	go run main.go

# Target to generate Swagger documentation
.PHONY: generate_swagger
generate_swagger:
	@swag init -g main.go -o ./deployments/swagger --parseDependency
    # Use 'swag' to generate Swagger documentation from 'main.go' and output to './deployments/swagger'
    # The '--parseDependency' flag tells Swag to parse dependencies recursively

# Target to validate Swagger documentation
.PHONY: validate_swagger
validate_swagger:
	docker run --rm -it -e GOPATH=$(HOME)/go:/go -v $(shell pwd)/deployments:/swagger-files -w $(shell pwd) quay.io/goswagger/swagger \
    validate /swagger-files/swagger/swagger.json
    # Run a Docker container with the Go Swagger CLI to validate the Swagger file
    # The '-e' flag sets the GOPATH environment variable
    # The '-v' flag mounts the local directory as a volume inside the container
    # The '-w' flag sets the working directory inside the container

# Target to generate and validate Swagger documentation
.PHONY: swagger
swagger: generate_swagger validate_swagger
    # This is a higher-level target that depends on both 'generate_swagger' and 'validate_swagger'
    # When 'make swagger' is run, it first generates Swagger documentation and then validates it


# Function to get the last commit hash
LAST_COMMIT := $(shell git rev-parse --short HEAD)

# Target to build and push Docker image
.PHONY: docker
docker:
	docker buildx build -t khedhrije/podcaster-indexer-api:latest --build-arg COMMIT_HASH=$(LAST_COMMIT) . && docker push khedhrije/podcaster-indexer-api:latest
    # Use Docker BuildX to build a multi-platform Docker image with the 'latest' tag
    # Pass the last commit hash as a build argument
    # After building, push the Docker image to the Docker registry
