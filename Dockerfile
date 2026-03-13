# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy internal/infrastructure/db/migrations for database migrations
COPY --from=builder /app/internal/infrastructure/db/migrations ./internal/infrastructure/db/migrations

# ✅ Copy certs folder
COPY --from=builder /app/certs ./certs

# Expose the API port
EXPOSE 8081

# Command to run the application
CMD ["./main"]
