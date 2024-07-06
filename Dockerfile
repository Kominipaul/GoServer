# Use the official Golang image as the base image
FROM golang:1.22.3-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/GoServer

# Use a minimal base image for the final image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the builder image
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

