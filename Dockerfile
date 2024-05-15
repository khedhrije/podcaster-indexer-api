# Use the official Golang image as builder
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /go/src

# Copy the current directory contents into the container
COPY . .

# Install git and build the Go application
RUN apt-get update \
    && apt-get install -y --no-install-recommends git \
    && cd /go/src \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -ldflags="-s -w" -o /app .

#-----------------------------------------------------------------------------------

# Use Distroless as the final base image for a minimal container
FROM gcr.io/distroless/base

# Copy the built executable from the builder stage into the final image
COPY --from=builder /app .

# Specify the command to run when the container starts
CMD ["/app"]
