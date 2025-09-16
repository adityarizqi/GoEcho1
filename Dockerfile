# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-krs-app main.go

# Stage 2: Create a minimal final image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /go-krs-app .

# Copy templates and .env file
COPY templates ./templates
COPY .env .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./go-krs-app"]