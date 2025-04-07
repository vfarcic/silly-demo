# Silly Demo

## Demo Manifests and Code Used in DevOps Toolkit Videos

[![Remote Environments with Dev Containers and Devpod: Are They Worth It?](https://img.youtube.com/vi/T_y3NTDrMHU/0.jpg)](https://youtu.be/T_y3NTDrMHU)
[![Say Goodbye to Tedious Docker Commands: Embrace Docker to Bake Images](https://img.youtube.com/vi/3Fc7YuTWptw/0.jpg)](https://youtu.be/3Fc7YuTWptw)
[![Stop Losing Requests! Learn Graceful Shutdown Techniques](https://img.youtube.com/vi/eQPYsGrZW_E/0.jpg)](https://youtu.be/eQPYsGrZW_E)

## Common

```sh
chmod +x dot.nu

./dot.nu setup --provider kind

source .env
```

## Backend with CNPG PostgreSQL

```sh
./dot.nu apply cnpg

./dot.nu apply atlas

kcl run kcl/main.k -D db.enabled=true \
    | kubectl --namespace a-team apply --filename -

curl "http://silly-demo.127.0.0.1.nip.io"

curl -X POST "http://silly-demo.127.0.0.1.nip.io/video?id=1&title=something"

curl -X POST "http://silly-demo.127.0.0.1.nip.io/video?id=2&title=else"

curl "http://silly-demo.127.0.0.1.nip.io/videos" | jq .
```

## Both frontend and backend with CNPG PostgreSQL

```sh
./dot.nu apply cnpg

./dot.nu apply atlas

kcl run kcl/main.k -D db.enabled=true -D frontend.enabled=true \
    | kubectl --namespace a-team apply --filename -

kubectl --namespace a-team \
    get all,ingresses,secrets,clusters,atlasschemas

curl -X POST "http://silly-demo.127.0.0.1.nip.io/video?id=1&title=something"

curl -X POST "http://silly-demo.127.0.0.1.nip.io/video?id=2&title=else"

curl "http://silly-demo.127.0.0.1.nip.io/videos" | jq .
```

## Unit Tests

```sh
./dot.nu run unit_tests
```

## Integration Tests

```sh
./dot.nu deploy app

source .env

./dot.nu run integration_tests

./dot.nu destroy app
```

## Destroy

```sh
./dot.nu destroy $PROVIDER

exit
```
