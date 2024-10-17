# Silly Demo

## Demo Manifests and Code Used in DevOps Toolkit Videos

[![Stop Losing Requests! Learn Graceful Shutdown Techniques](https://img.youtube.com/vi/eQPYsGrZW_E/0.jpg)](https://youtu.be/eQPYsGrZW_E)

## Common

```sh
kind create cluster --config kind.yaml

kubectl apply \
    --filename https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kubectl create namespace a-team
```

## App Alone

```sh
kubectl --namespace a-team apply --filename k8s
```

## App with CNPG PostgreSQL

```sh
helm upgrade --install cnpg cloudnative-pg \
    --repo https://cloudnative-pg.github.io/charts \
    --namespace cnpg-system --create-namespace --wait

helm upgrade --install atlas-operator \
    oci://ghcr.io/ariga/charts/atlas-operator \
    --namespace atlas-operator --create-namespace --wait

timoni build silly-demo timoni \
    --values timoni/values-db-cnpg.yaml --namespace a-team \
    | kubectl apply --filename -
```

## App with OTEL

```sh
timoni build silly-demo timoni \
    --values timoni/values-otel.yaml --namespace a-team \
    | kubectl apply --filename -
```

## App with NATS

```sh
timoni build silly-demo timoni \
    --values timoni/values-nats.yaml --namespace a-team \
    | kubectl apply --filename -
```
