# Silly Demo

# Common

```sh
kubectl create namespace a-team
```

## App Alone

```sh
kubectl apply --filename k8s
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
kubectl create namespace a-team

timoni build silly-demo timoni \
    --values timoni/values-otel.yaml --namespace a-team \
    | kubectl apply --filename -
```
