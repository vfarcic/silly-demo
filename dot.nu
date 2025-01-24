#!/usr/bin/env nu

source  scripts/image.nu
source  scripts/tests.nu
source  scripts/kubernetes.nu
source  scripts/ingress.nu
source  scripts/cert-manager.nu

def main [] {}

# Creates a local Kubernetes cluster
def "main setup" [] {

    main create kubernetes kind 

    main apply ingress nginx --hyperscaler kind

    kubectl create namespace a-team
    
}

# Updates Timoni files
def "main update timoni" [
    tag: string # The tag of the image (e.g., 0.0.1)
] {

    cat timoni/values.cue
        | sed -e $"s@image: tag:.*@image: tag: \"($tag)\"@g"
        | save timoni/values.cue.tmp --force

    mv timoni/values.cue.tmp timoni/values.cue

}

# Signs the image
def "main sign image" [
    tag: string                    # The tag of the image (e.g., `0.0.1`)
    --registry_pass: string,       # Registry password. Overwrites environment variable `REGISTRY_PASSWORD`.
    --registry_user = "vfarcic",   # Registry username
    --cosign_private_key: string,  # Cosign private key. Overwrites environment variable `COSIGN_PRIVATE_KEY`.
    --registry = "ghcr.io/vfarcic" # Image registry
    --image = "silly-demo"         # Image name
] {

    mut registry_pass = get_registry_pass $registry_pass

    if $cosign_private_key != null {
        $env.COSIGN_PRIVATE_KEY = $cosign_private_key
    }

    (
        cosign sign --yes --key env://COSIGN_PRIVATE_KEY
            --registry-username $registry_user
            --registry-password $registry_pass
            $"($registry)/($image):($tag)"
    )

}

# Updates Helm files
def "main build helm" [
    tag: string                    # The tag of the image (e.g., 0.0.1)
    --push = true                  # Whether to push the chart to the registry
    --registry = "ghcr.io/vfarcic" # Image registry
    --registry_pass: string,       # Registry password. Overwrites environment variable `REGISTRY_PASSWORD`.
    --registry_user = "vfarcic"    # Registry username
] {

    mut registry_pass = get_registry_pass $registry_pass

    open helm/app/Chart.yaml
        | upsert version $tag
        | save helm/app/Chart.yaml --force

    open helm/app/values.yaml
        | upsert image.tag $tag
        | save helm/app/values.yaml --force

    helm package helm/app

    if $push {
        (
            helm registry login
                --username $registry_user
                --password $registry_pass
                $registry
        )
    }

    helm push $"silly-demo-helm-($tag).tgz" $"oci://($registry)"

}

def get_registry_pass [registry_pass] {
    mut registry_pass = $registry_pass
    if $registry_pass == "" or $registry_pass == null {
        $registry_pass = $env.REGISTRY_PASSWORD
    }
    $registry_pass
}

# Updates Kustomize files
def "main update kustomize" [
    tag: string                    # The tag of the image (e.g., 0.0.1)
    --registry = "ghcr.io/vfarcic" # Image registry
    --image = "silly-demo"         # Image name
] {

    open kustomize/base/deployment.yaml
        | upsert spec.template.spec.containers.0.image $"($registry)/($image):($tag)"
        | save kustomize/base/deployment.yaml --force

}

# Updates YAML files
def "main update yaml" [
    tag: string                    # The tag of the image (e.g., 0.0.1)
    --registry = "ghcr.io/vfarcic" # Image registry
    --image = "silly-demo"         # Image name
] {

    open k8s/deployment.yaml
        | upsert spec.template.spec.containers.0.image $"($registry)/($image):($tag)"
        | save k8s/deployment.yaml --force

}

# Runs all CI tasks
def "main run ci" [
    tag: string                    # The tag of the image (e.g., 0.0.1)
    --registry = "ghcr.io/vfarcic" # Image registry
    --image = "silly-demo"         # Image name
] {

    main run tests --language go

    main build image $tag

    main sign image $tag

    main update timoni $tag

    main build helm $tag

    main update kustomize $tag

    main update yaml $tag

}
