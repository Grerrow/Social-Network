# Build stage
FROM golang:1.24-alpine AS builder  

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o social-network .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/social-network .

# Copy migrations
COPY --from=builder /app/db/migrations ./db/migrations

# Create directories for data persistence
RUN mkdir -p /app/db /app/public/images/avatars /app/public/images/posts /app/public/images/messages /app/public/images/comments /app/public/images/groups

# Expose port
EXPOSE 8080

# Set environment variables with defaults
ENV PORT=:8080
ENV DB_PATH=./db/social-network.db
ENV UPLOAD_PATH=./public/images
# ENV ALLOWED_ORIGIN=http://localhost

# Run the application
CMD ["./social-network"]