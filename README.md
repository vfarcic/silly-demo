```bash
export VERSION=1.2.22

cat Dockerfile \
    | sed -e "s@ENV VERSION .*@ENV VERSION $VERSION@g" \
    | tee Dockerfile.tmp

mv Dockerfile.tmp Dockerfile

docker image build -t c8n.io/vfarcic/silly-demo:latest .

docker image tag c8n.io/vfarcic/silly-demo:latest c8n.io/vfarcic/silly-demo:$VERSION

docker image push c8n.io/vfarcic/silly-demo:latest

docker image push c8n.io/vfarcic/silly-demo:$VERSION

cosign sign --key cosign/cosign.key c8n.io/vfarcic/silly-demo:$VERSION

cosign verify --key c8n.io/vfarcic/silly-demo:$VERSION

yq --inplace ".values.image.tag = \"$VERSION\"" timoni/values.yaml

cat timoni/values.cue \
    | sed -e "s@image: tag:.*@image: tag: \"$VERSION\"@g" \
    | tee timoni/values.cue.tmp

mv timoni/values.cue.tmp timoni/values.cue

timoni mod push timoni oci://c8n.io/vfarcic/silly-demo-package --version $VERSION
```