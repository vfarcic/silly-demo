FROM golang:1.25.5-alpine AS build
RUN mkdir /src
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
ADD ./vendor .
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o silly-demo
RUN chmod +x silly-demo

FROM scratch
ARG VERSION
ENV VERSION=$VERSION
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
EXPOSE 8080
CMD ["silly-demo"]
