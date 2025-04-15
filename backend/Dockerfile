# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

# Expose port 3001
EXPOSE 3001

# Command to run the executable
CMD ["./main"] 