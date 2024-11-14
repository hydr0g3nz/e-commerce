# Build stage
FROM golang:1.23-alpine AS builder

# Set necessary environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install essential build tools
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/server.go
# RUN go build -ldflags="-w -s" -o main ./cmd/server.go

# Final stage
FROM alpine:3.19

# Add non root user
# RUN addgroup -S app && adduser -S app -G app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main /app
# Copy config file

# Switch to non root user
# USER app

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./main"]