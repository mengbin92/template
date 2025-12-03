# Build stage
FROM golang:1.25 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    make \
    gcc \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /src

# Set GOPROXY
ENV GOPROXY=https://goproxy.cn

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Runtime stage
FROM debian:stable-slim

# Install CA certificates and timezone data
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -s /bin/sh -m appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /src/bin ./bin

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8000
EXPOSE 9000

# Create volume mount point
VOLUME /data/conf

# Run the application
CMD ["./bin/app", "-conf", "/data/conf"]

