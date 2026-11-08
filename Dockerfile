# HostMK Security Framework Dockerfile
# Multi-stage build for production-ready container

# ============================================================================
# Stage 1: Builder
# ============================================================================
FROM golang:1.26-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct \
    GONOSUMDB=* \
    GONOPROXY=*

# Install dependencies
RUN apk add --no-cache \
    ca-certificates \
    git \
    gcc \
    musl-dev \
    && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Verify dependencies
RUN go mod verify

# Build the application
RUN go build \
    -ldflags="-s -w -X main.Version=${VERSION:-1.0.0} -X main.Commit=${COMMIT:-unknown} -X main.Date=${DATE:-$(date +%Y-%m-%dT%H:%M:%SZ)} -X main.BuiltBy=docker" \
    -o /app/hostmk-security \
    ./cmd/server

# ============================================================================
# Stage 2: Runtime
# ============================================================================
FROM alpine:3.20

# Set environment variables
ENV APP_NAME=hostmk-security \
    APP_VERSION=${VERSION:-1.0.0} \
    APP_ENV=${ENVIRONMENT:-production} \
    PORT=8080 \
    HOST=0.0.0.0

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Set timezone
ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/hostmk-security /app/hostmk-security

# Copy configuration files
COPY --chown=appuser:appgroup config.yaml /app/config.yaml
COPY --chown=appuser:appgroup .env /app/.env

# Set permissions
RUN chmod 755 /app/hostmk-security && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set security labels
LABEL \
    org.opencontainers.image.title="HostMK Security Framework" \
    org.opencontainers.image.description="Enterprise Security Framework for SaaS Platforms" \
    org.opencontainers.image.source="https://github.com/hostpanel22/security-hostmk22-enterprise" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.vendor="HostMK Team" \
    org.opencontainers.image.version="${VERSION:-1.0.0}" \
    org.opencontainers.image.created="${DATE:-$(date +%Y-%m-%dT%H:%M:%SZ)}" \
    org.opencontainers.image.revision="${COMMIT:-unknown}"

# Security hardening
RUN mkdir -p /app/logs && \
    chown appuser:appgroup /app/logs && \
    chmod 750 /app/logs

# Set security options
# Note: These are set in docker run command for better compatibility
# docker run --security-opt no-new-privileges --cap-drop=ALL ...

# Entrypoint
ENTRYPOINT ["/app/hostmk-security"]

# Default command (can be overridden)
CMD ["serve"]
