
# Builder Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all binaries inside cmd/
RUN mkdir -p /app/bin && \
    for cmd in cmd/*/; do \
        binary=$(basename "$cmd"); \
        echo "Building ${binary}..."; \
        CGO_ENABLED=0 GOOS=linux go build \
            -ldflags="-s -w" \
            -o /app/bin/"${binary}" \
            ./"${cmd}"; \
    done


# Runtime Stage

FROM alpine:3.22

# Install runtime dependencies
RUN apk add --no-cache ca-certificates wget

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy application binary
COPY --from=builder /app/bin/* .

# Copy Swagger documentation
COPY --from=builder /app/docs ./docs

# COPY --from=builder /app/uploads ./uploads

# Change ownership
RUN chown -R app:app /app

# Run as non-root
USER app

EXPOSE 1357

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:1357/health || exit 1

CMD ["./api"]