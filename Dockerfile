# Start with a base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and build files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Explicitly copy the web directory (assuming it's in the project root)
COPY web/ ./web/

# Build the Go app (adjust the path based on your actual main package location)
RUN go build -o app ./cmd/GoServer

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/app .

# Explicitly copy the web directory from the builder stage
COPY --from=builder /app/web ./web

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
