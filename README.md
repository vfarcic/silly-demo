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