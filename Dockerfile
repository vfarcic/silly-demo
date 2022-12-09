FROM golang:1.19-alpine AS build
ADD . /src
WORKDIR /src
RUN go get -d -v -t
RUN GOOS=linux GOARCH=amd64 go build -v -o silly-demo 

FROM alpine:3.17.0
EXPOSE 8080
CMD ["silly-demo"]
ENV VERSION 1.1.0
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
RUN chmod +x /usr/local/bin/silly-demo
