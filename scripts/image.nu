#!/usr/bin/env nu

# Builds a container image
def "main build image" [
    tag: string                    # The tag of the image (e.g., 0.0.1)
    --registry = "ghcr.io/vfarcic" # Image registry
    --image = "silly-demo"         # Image name
    --push = true                  # Whether to push the image to the registry
] {

    mut pushArg = ""
    if $push {
        $pushArg = "--push"
    }

    TAG=$tag docker buildx bake $pushArg

    # docker image build --tag $"($registry)/($image):latest" .

    # docker image tag $"($registry)/($image):latest" $"($registry)/($image):($tag)"

    # if $push {

    #     docker image push $"($registry)/($image):latest"

    #     docker image push $"($registry)/($image):($tag)"
    # }

}
