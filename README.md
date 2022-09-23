```bash
export VERSION=1.0.9

docker image build -t vfarcic/silly-demo:latest .

docker image tag vfarcic/silly-demo:latest vfarcic/silly-demo:$VERSION

docker image push vfarcic/silly-demo:latest

docker image push vfarcic/silly-demo:$VERSION

cosign sign --key cosign.key vfarcic/silly-demo:latest

cosign sign --key cosign.key vfarcic/silly-demo:$VERSION

cosign verify --key cosign.pub vfarcic/silly-demo:latest vfarcic/silly-demo:$VERSION
```