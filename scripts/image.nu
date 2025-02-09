#!/usr/bin/env nu

# Builds a container image using `docker buildx bake`
def "main build image" [
    tag: string                    # The tag of the image (e.g., 0.0.1)
    --registry = "ghcr.io/vfarcic" # Image registry
    --image = "silly-demo"         # Image name
    --push = true                  # Whether to push the image to the registry
    --bake_target = default        # Which Docker Bake target to use if `--bake` is `true`
] {

    if $push {

        TAG=$tag IMAGE=$"($registry)/($image)" docker buildx bake $bake_target --push

    } else {

        TAG=$tag IMAGE=$"($registry)/($image)" docker buildx bake $bake_target

    }

}

# Retrieves a container registry address
def "main get container_registry" [] {

    mut registry = ""
    if "CONTAINER_REGISTRY" in $env {
        $registry = $env.CONTAINER_REGISTRY
    } else {
        let value = input $"(ansi green_bold)Enter container image registry \(e.g., `ghcr.io/vfarcic`\):(ansi reset) "
        $registry = $value
    }
    $"export CONTAINER_REGISTRY=($registry)\n" | save --append .env

    $registry

}
