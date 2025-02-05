# Silly Demo

## Demo Manifests and Code Used in DevOps Toolkit Videos

[![Stop Losing Requests! Learn Graceful Shutdown Techniques](https://img.youtube.com/vi/eQPYsGrZW_E/0.jpg)](https://youtu.be/eQPYsGrZW_E)

## Common

```sh
chmod +x dot.nu

./dot.nu setup

source .env
```

## App Alone

```sh
kubectl --namespace a-team apply --filename k8s
```

## App with CNPG PostgreSQL

```sh
./dot.nu apply cnpg

./dot.nu apply atlas

kcl run kcl/main.k -D db.enabled=true \
    | kubectl --namespace a-team apply --filename -

kubectl --namespace a-team \
    get all,ingresses,secrets,clusters,atlasschemas

curl -X POST "http://silly-demo.127.0.0.1.nip.io/video?id=1&title=something"

curl -X POST "http://silly-demo.127.0.0.1.nip.io/video?id=2&title=else"

curl "http://silly-demo.127.0.0.1.nip.io/videos" | jq .
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

## Unit Tests

```sh
./dot.nu run unit_tests
```

## Integration Tests

```sh
./dot.nu deploy app

./dot.nu run integration_tests

./dot.nu destroy app
```

## Destroy

```sh
./dot.nu destroy
```
