FROM golang:1.20.10-alpine3.17 AS build
RUN mkdir /src
ADD ./*.go /src
ADD ./go.mod /src
ADD ./go.sum /src
WORKDIR /src
RUN go get -d -v -t
RUN GOOS=linux GOARCH=amd64 go build -v -o silly-demo 
RUN chmod +x silly-demo

FROM scratch
ENV DB_PORT=5432 DB_USER=postgres DB_NAME=silly-demo
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
EXPOSE 8080
CMD ["silly-demo"]
