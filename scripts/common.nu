#!/usr/bin/env nu

def "main get provider" [] {
    let hyperscaler = [aws azure google local upcloud]
        | input list $"(ansi yellow_bold)Which Hyperscaler do you want to use?(ansi green_bold)"
    print $"(ansi reset)"

    $"export HYPERSCALER=($hyperscaler)\n" | save --append .env

    $hyperscaler
}

def "main print source" [] {

    print $"
Execute `(ansi yellow_bold)source .env(ansi reset)` to load the environment variables.
"

}

def "main delete temp_files" [] {

    rm --force .env

    rm --force kubeconfig*.yaml

}

def --env "main get creds" [
    hyperscaler: string,
] {

    mut creds = {hyperscaler: $hyperscaler}

    if $hyperscaler == "google" {

        gcloud auth login


    } else if $hyperscaler == "aws" {

        mut aws_access_key_id = ""
        if AWS_ACCESS_KEY_ID in $env {
            $aws_access_key_id = $env.AWS_ACCESS_KEY_ID
        } else {
            $aws_access_key_id = input $"(ansi green_bold)Enter AWS Access Key ID: (ansi reset)"
        }
        $"export AWS_ACCESS_KEY_ID=($aws_access_key_id)\n"
            | save --append .env
        $creds = ( $creds | upsert aws_access_key_id $aws_access_key_id )

        mut aws_secret_access_key = ""
        if AWS_SECRET_ACCESS_KEY in $env {
            $aws_secret_access_key = $env.AWS_SECRET_ACCESS_KEY
        } else {
            $aws_secret_access_key = input $"(ansi green_bold)Enter AWS Secret Access Key: (ansi reset)" --suppress-output
            print ""
        }
        $"export AWS_SECRET_ACCESS_KEY=($aws_secret_access_key)\n"
            | save --append .env
        $creds = ( $creds | upsert aws_secret_access_key $aws_secret_access_key )

        mut aws_account_id = ""
        if AWS_ACCOUNT_ID in $env {
            $aws_account_id = $env.AWS_ACCOUNT_ID
        } else {
            $aws_account_id = input $"(ansi green_bold)Enter AWS Account ID: (ansi reset)"
        }
        $"export AWS_ACCOUNT_ID=($aws_account_id)\n"
            | save --append .env
        $creds = ( $creds | upsert aws_account_id $aws_account_id )

    } else if $hyperscaler == "azure" {

        mut tenant_id = ""

        if AZURE_TENANT in $env {
            $tenant_id = $env.AZURE_TENANT
        } else {
            $tenant_id = input $"(ansi green_bold)Enter Azure Tenant ID: (ansi reset)"
        }
        $creds = ( $creds | upsert tenant_id $tenant_id )

        az login --tenant $tenant_id
    
    } else {

        print $"(ansi red_bold)($hyperscaler)(ansi reset) is not a supported."
        exit 1

    }

    $creds

}
