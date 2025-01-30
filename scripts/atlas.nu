#!/usr/bin/env nu

def "main apply atlas" [] {

    (
        helm upgrade --install atlas-operator 
            oci://ghcr.io/ariga/charts/atlas-operator 
            --namespace atlas-operator --create-namespace
            --wait
    )

}
