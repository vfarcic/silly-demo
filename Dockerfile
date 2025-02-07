# Build stage
FROM golang:1.23.4-alpine AS build
RUN mkdir /src
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o silly-demo
RUN chmod +x silly-demo

FROM scratch
ARG VERSION
ENV VERSION=$VERSION
ENV DB_PORT=5432 DB_USERNAME=postgres DB_NAME=silly-demo
COPY cache /cache
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
