# Build stage
FROM golang:1.22.4-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install tzdata and set timezone
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tokyo

# Copy the binary from builder
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x ./main

# Set DOCKER_ENV
ENV DOCKER_ENV=true

EXPOSE 8080

CMD ["./main"]
