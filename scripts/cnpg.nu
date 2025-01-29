#!/usr/bin/env nu

def "main apply cnpg" [
    atlas = true
] {

    (
        helm upgrade --install cnpg cloudnative-pg
            --repo https://cloudnative-pg.github.io/charts
            --namespace cnpg-system --create-namespace --wait
    )

    if $atlas {

        (
            helm upgrade --install atlas-operator 
                oci://ghcr.io/ariga/charts/atlas-operator 
                --namespace atlas-operator --create-namespace
                --wait
        )

    }

}
