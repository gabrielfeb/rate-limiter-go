# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# -ldflags="-w -s" strips debug information, reducing binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /rate-limiter-app ./main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /rate-limiter-app .

# Expose the port the app runs on
EXPOSE 8080

# Define the command to run the application
CMD ["./rate-limiter-app"]