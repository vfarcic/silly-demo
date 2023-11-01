FROM golang:1.20.10-alpine3.17 AS build
RUN mkdir /src
WORKDIR /src
ADD ./go.mod .
ADD ./go.sum .
RUN go mod download
ADD ./*.go .
RUN ls -la
RUN GOOS=linux GOARCH=amd64 go build -v -o silly-demo 
RUN chmod +x silly-demo

FROM scratch
ENV DB_PORT=5432 DB_USER=postgres DB_NAME=silly-demo
COPY --from=build /src/silly-demo /usr/local/bin/silly-demo
EXPOSE 8080
CMD ["silly-demo"]
