VERSION 0.8
FROM golang:1.22.2-alpine
# WORKDIR /go-workdir

binary:
    COPY go.mod go.sum vendor .
    COPY *.go .
    RUN GOOS=linux GOARCH=amd64 go build -o silly-demo
    SAVE ARTIFACT silly-demo

timoni:
    ARG user=vfarcic
    RUN go install github.com/stefanprodan/timoni/cmd/timoni@latest
    COPY timoni/values.cue timoni/values.cue
    RUN cat timoni/values.cue | sed -e "s@image: tag:.*@image: tag: \"9.9.9\"@g" >timoni/values.cue.tmp
    SAVE ARTIFACT timoni/values.cue.tmp AS LOCAL timoni/values.cue
    RUN --push --secret password=password timoni mod push timoni oci://ghcr.io/$user/silly-demo-package --version 9.9.9 --creds $user:$password

image-common:
    BUILD +binary
    ARG tag='latest'
    ARG base='scratch'
    FROM $base
    ENV DB_PORT=5432 DB_USERNAME=postgres DB_NAME=silly-demo
    EXPOSE 8080
    CMD ["silly-demo"]
    ENV VERSION=$tag
    COPY +binary/silly-demo /usr/local/bin/silly-demo
    SAVE IMAGE --push ghcr.io/vfarcic/silly-demo:$tag

image:
    ARG tag=latest
    BUILD +image-common --tag $tag --base scratch

image-alpine:
    ARG tag=latest
    BUILD +image-common --tag $tag-alpine --base alpine:3.18.4

image-all:
    ARG tag=latest
    BUILD +image --tag latest --tag $tag
    BUILD +image-alpine --tag latest --tag $tag