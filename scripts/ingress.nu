#!/usr/bin/env nu

# Applies Ingress
#
# Examples:
# > main apply ingress "contour" --hyperscaler aws
def --env "main apply ingress" [
    class = "traefik" # The class of Ingress controller to apply. Available options: traefik, contour, nginx
    --hyperscaler = "none"
    --env_prefix = ""
] {

    if $class == "traefik" {

        (
            helm upgrade --install traefik traefik
                --repo https://helm.traefik.io/traefik
                --namespace traefik --create-namespace --wait
        )

    } else if $class == "contour" {

        (
            helm upgrade --install contour 
                oci://registry-1.docker.io/bitnamicharts/contour
                --namespace contour --create-namespace --wait
        )
    
    } else if $class == "nginx" {

        if $hyperscaler == "kind" {

            (
                kubectl apply
                    --filename https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
            )

            sleep 10sec

            (
                kubectl --namespace ingress-nginx wait
                    --for=condition=Available
                    deployment ingress-nginx-controller
            )

            sleep 5sec

            (
                kubectl --namespace ingress-nginx wait
                    --for=condition=Complete
                    job ingress-nginx-admission-create
            )

            (
                kubectl --namespace ingress-nginx wait
                    --for=condition=Complete
                    job ingress-nginx-admission-patch
            )

        }

    } else {

        print $"(ansi red_bold)($class)(ansi reset) is not a supported."
        exit 1

    }

    main get ingress $class --hyperscaler $hyperscaler --env_prefix $env_prefix

}

def "main get ingress" [
    class = "traefik" # The class of Ingress controller to apply. Available options: traefik, contour, nginx
    --hyperscaler: string
    --env_prefix = ""
] {

    mut service_name = $class

    if $class == "contour" {
        $service_name = "contour-envoy"
    }
    
    mut ingress_ip = ""
  
    if $hyperscaler == "aws" {

        sleep 30sec

        let ingress_hostname = (
            kubectl --namespace $class
                get service $service_name --output yaml
                | from yaml
                | get status.loadBalancer.ingress.0.hostname
        )

        while $ingress_ip == "" {
            print "Waiting for Ingress Service IP..."
            sleep 10sec
            $ingress_ip = (dig +short $ingress_hostname)
        }

        $ingress_ip = $ingress_ip | lines | first

    } else if $hyperscaler == "kind" {

        $ingress_ip = "127.0.0.1"

    } else {

        while $ingress_ip == "" {

            print $"Waiting for ($class) Ingress IP from ($service_name) Service..."

            sleep 10sec

            $ingress_ip = (
                kubectl --namespace $class
                    get service $service_name --output yaml
                    | from yaml
                    | get status.loadBalancer.ingress.0.ip
            )

        }
    }

    $"export ($env_prefix)INGRESS_IP=($ingress_ip)\n" | save --append .env
    $"export ($env_prefix)INGRESS_HOST=($ingress_ip).nip.io\n" | save --append .env

    {ip: $ingress_ip, host: $"($ingress_ip).nip.io", class: $class}

}