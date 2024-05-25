VERSION 0.8
# FROM golang:1.22.2-alpine
FROM ghcr.io/vfarcic/silly-demo-earthly:0.0.5
ARG --global user=vfarcic
# WORKDIR /go-workdir

binary:
    COPY go.mod go.sum vendor .
    COPY *.go .
    RUN go mod init
    RUN go mod vendor
    RUN GOOS=linux GOARCH=amd64 go build -mod vendor -o silly-demo
    SAVE ARTIFACT silly-demo

timoni:
    COPY timoni/values.cue timoni/values.cue
    RUN cat timoni/values.cue \
        | sed -e "s@image: tag:.*@image: tag: \"9.9.9\"@g" \
        >timoni/values.cue.tmp
    SAVE ARTIFACT timoni/values.cue.tmp AS LOCAL timoni/values.cue
    RUN --push --secret password \
        timoni mod push timoni \
        oci://ghcr.io/$user/silly-demo-package --version 9.9.9 \
        --creds $user:$password

cosign:
    ARG tag='latest'
    RUN echo "USER: $user"
    RUN --push \
        --secret COSIGN_PASSWORD=cosignpassword \
        --secret cosignkey \
        cosign sign --yes --key env://cosignkey \
        --registry-username $user \
        ghcr.io/vfarcic/silly-demo:$tag

helm:
    ARG tag='latest'
    COPY helm helm
    RUN yq --inplace ".version = \"$tag\"" helm/app/Chart.yaml
    SAVE ARTIFACT helm/app/Chart.yaml AS LOCAL helm/app/Chart.yaml
    RUN yq --inplace ".image.tag = \"$tag\"" helm/app/values.yaml
    SAVE ARTIFACT helm/app/values.yaml AS LOCAL helm/app/values.yaml
    RUN helm package helm/app
    RUN --push helm push silly-demo-helm-$tag.tgz oci://ghcr.io/vfarcic

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
    ARG tag
    BUILD +image --tag latest --tag $tag
    BUILD +image-alpine --tag latest --tag $tag
    # BUILD +cosign --tag latest --tag $tag
    # BUILD +cosign --tag latest-alpine --tag $tag-alpin

cosign-all:
    ARG --required tag
    BUILD +cosign --tag latest --tag $tag
    BUILD +cosign --tag latest-alpine --tag $tag-alpine
