VERSION 0.8
FROM ghcr.io/vfarcic/silly-demo-earthly:0.0.5
ARG --global registry=ghcr.io/vfarcic
ARG --global user=vfarcic
ARG --global image=silly-demo
WORKDIR /go-workdir

binary:
    COPY go.mod go.sum vendor .
    COPY *.go .
    RUN go mod vendor
    RUN GOOS=linux GOARCH=amd64 go build --mod vendor -o silly-demo
    SAVE ARTIFACT silly-demo

image:
    BUILD +binary
    ARG tag='latest'
    ARG taglatest='latest'
    ARG base='scratch'
    FROM $base
    ENV DB_PORT=5432 DB_USERNAME=postgres DB_NAME=silly-demo
    COPY cache /cache
    EXPOSE 8080
    CMD ["silly-demo"]
    ENV VERSION=$tag
    COPY +binary/silly-demo /usr/local/bin/silly-demo
    SAVE IMAGE --push \
        $registry/$image:$tag \
        $registry/$image:$taglatest

timoni:
    ARG --required tag
    COPY timoni/values.cue timoni/values.cue
    RUN cat timoni/values.cue \
        | sed -e "s@image: tag:.*@image: tag: \"$tag\"@g" \
        >timoni/values.cue.tmp
    SAVE ARTIFACT timoni/values.cue.tmp AS LOCAL timoni/values.cue
    RUN --push --secret password \
        timoni mod push timoni \
        oci://$registry/$image-package --version $tag \
        --creds $user:$password

cosign:
    ARG --required tag
    RUN --push \
        --secret COSIGN_PASSWORD=cosignpassword \
        --secret cosignkey \
        --secret password \
        cosign sign --yes --key env://cosignkey \
        --registry-username $user \
        --registry-password $password \
        $registry/$image:$tag

helm:
    ARG --required tag
    COPY helm helm
    RUN yq --inplace ".version = \"$tag\"" helm/app/Chart.yaml
    SAVE ARTIFACT helm/app/Chart.yaml AS LOCAL helm/app/Chart.yaml
    RUN yq --inplace ".image.tag = \"$tag\"" helm/app/values.yaml
    SAVE ARTIFACT helm/app/values.yaml AS LOCAL helm/app/values.yaml
    RUN helm package helm/app
    RUN ls -l
    RUN --secret password helm registry login \
        --username $user --password $password $registry
    RUN --push helm push silly-demo-helm-$tag.tgz oci://$registry

kustomize:
    ARG --required tag
    COPY kustomize kustomize
    RUN yq --inplace ".spec.template.spec.containers[0].image = \"ghcr.io/vfarcic/silly-demo:$tag\"" kustomize/base/deployment.yaml
    SAVE ARTIFACT kustomize/base/deployment.yaml AS LOCAL kustomize/base/deployment.yaml

kubernetes:
    ARG --required tag
    COPY k8s k8s
    RUN yq --inplace ".spec.template.spec.containers[0].image = \"ghcr.io/vfarcic/silly-demo:$tag\"" k8s/deployment.yaml
    SAVE ARTIFACT k8s/deployment.yaml AS LOCAL k8s/deployment.yaml

all:
    ARG tag
    WAIT
        BUILD +image --tag $tag --taglatest latest
        BUILD +image --tag $tag-alpine \
            --taglatest latest-alpine --base alpine:3.18.4
    END
    BUILD +cosign --tag latest --tag $tag \
        --tag latest-alpine --tag $tag-alpine
    BUILD +timoni --tag $tag
    BUILD +helm --tag $tag
    BUILD +kustomize --tag $tag
    BUILD +kubernetes --tag $tag

