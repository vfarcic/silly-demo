#!/usr/bin/env nu

def "main apply certmanager" [] {

    helm repo add jetstack https://charts.jetstack.io --force-update

    helm repo update

    (
        helm upgrade --install cert-manager jetstack/cert-manager
            --namespace cert-manager --create-namespace
            --set crds.enabled=true --wait
    )

}