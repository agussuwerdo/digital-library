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

# Create a default .env file if not exists
RUN if [ ! -f .env ]; then echo "DATABASE_URL=postgres://postgres:postgres@db:5432/digital_library?sslmode=disable\nJWT_SECRET=your_jwt_secret_here" > .env; fi

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