# Stage 1: Build the Go application
FROM golang:1.23.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o webhook-listener cmd/main.go

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/webhook-listener .

# Copy the .env file (if needed)
COPY .env /app/.env

# Expose the port the app runs on
EXPOSE 5001

# Set environment variables
ENV PORT=5001

# Run the binary
CMD ["/app/webhook-listener"]
