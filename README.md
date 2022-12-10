```bash
export VERSION=1.1.4

cat Dockerfile \
    | sed -e "s@ENV VERSION .*@ENV VERSION $VERSION@g" \
    | tee Dockerfile

docker image build -t vfarcic/silly-demo:latest .

docker image tag vfarcic/silly-demo:latest vfarcic/silly-demo:$VERSION

docker image push vfarcic/silly-demo:latest

docker image push vfarcic/silly-demo:$VERSION

cosign sign --key cosign/cosign.key vfarcic/silly-demo:latest

cosign sign --key cosign/cosign.key vfarcic/silly-demo:$VERSION

cosign verify --key cosign/cosign.pub vfarcic/silly-demo:latest vfarcic/silly-demo:$VERSION
```