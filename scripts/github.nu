#!/usr/bin/env nu

def --env "main get github" [] {

    mut github_token = ""
    if "GITHUB_TOKEN" in $env {
        $github_token = $env.GITHUB_TOKEN
    } else if "REGISTRY_PASSWORD" in $env {
        $github_token = $env.REGISTRY_PASSWORD
    } else {
        let value = input $"(ansi green_bold)Enter GitHub token:(ansi reset) "
        $github_token = $value
    }
    $"export GITHUB_TOKEN=($github_token)\n" | save --append .env

    mut github_org = ""
    if "GITHUB_ORG" in $env {
        $github_org = $env.GITHUB_ORG
    } else if "REGISTRY_USER" in $env {
        $github_org = $env.REGISTRY_USER
    } else {
        let value = input $"(ansi green_bold)Enter GitHub user or organization where you forked the repo:(ansi reset) "
        $github_org = $value
    }
    $"export GITHUB_ORG=($github_org)\n" | save --append .env

    {org: $github_org, token: $github_token}

}
