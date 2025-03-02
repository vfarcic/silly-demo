variable "IMAGE" {
    default = "ghcr.io/vfarcic/silly-demo"
}
variable "TAG" {
    default = "dev"
}
target "default" {
    name = item.name
    matrix = {
        item = [{
            name = "backend"
            context = "."
            tags = ["${IMAGE}:${TAG}", "${IMAGE}:latest"]
        }, {
            name = "frontend"
            context = "./frontend"
            tags = ["${IMAGE}-frontend:${TAG}", "${IMAGE}-frontend:latest"]
        }]
    }
    tags = item.tags
    dockerfile = "Dockerfile"
    context = item.context
    # platforms = ["linux/amd64", "linux/arm64/v8"]
    platforms = ["linux/amd64"]
    args = {
        VERSION = TAG
    }
}
