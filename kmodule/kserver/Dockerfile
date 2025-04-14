# Use the official Golang image as a base
FROM golang:1.23.4 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Use a minimal image for the final build
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x /app/main

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
