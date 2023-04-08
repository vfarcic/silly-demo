FROM golang:1.20-alpine AS build
ADD ./*.go /src
ADD ./go.mod /src
ADD ./go.sum /src
WORKDIR /src
RUN go get -d -v -t
RUN GOOS=linux go build -v -o silly-demo 

FROM alpine:3.17.3
EXPOSE 8080
CMD ["silly-demo"]
ENV VERSION 1.1.4
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
RUN chmod +x /usr/local/bin/silly-demo
