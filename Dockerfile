FROM golang:1.20.10-alpine3.17 AS build
RUN mkdir /src
WORKDIR /src
ADD ./go.mod .
ADD ./go.sum .
ADD ./vendor .
ADD ./*.go .
RUN GOOS=linux GOARCH=amd64 go build -o silly-demo
RUN chmod +x silly-demo

FROM scratch
ENV DB_PORT=5432 DB_USERNAME=postgres DB_NAME=silly-demo
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
EXPOSE 8080
CMD ["silly-demo"]
