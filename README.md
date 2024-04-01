# Silly Demo

## App Alone

```sh
kubectl apply --filename k8s
```

## App with CNPG PostgreSQL

```sh
timoni build silly-demo timoni \
    --values timoni/values-db-cnpg.yaml \
    | kubectl apply --filename -
```

## App with OTEL

```sh
kubectl create namespace a-team

timoni build silly-demo timoni \
    --values timoni/values-otel.yaml --namespace a-team \
    | kubectl apply --filename -
```
