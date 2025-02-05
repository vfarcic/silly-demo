# Build stage
FROM golang:1.23.4-alpine AS build
RUN apk add --no-cache ca-certificates git tzdata
RUN adduser -D -g '' appuser

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build with optimizations and proper version
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/silly-demo

# Final stage
FROM scratch
ARG VERSION
ENV VERSION=$VERSION
ENV DB_PORT=5432 DB_USERNAME=postgres DB_NAME=silly-demo
COPY cache /cache

# Use non-root user
USER appuser

# Add metadata labels
LABEL org.opencontainers.image.source="https://github.com/vfarcic/silly-demo" \
      org.opencontainers.image.description="A silly demo application" \
      org.opencontainers.image.version="${VERSION}"

# Add healthcheck using the ping endpoint
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1

EXPOSE 8080
CMD ["silly-demo"]
