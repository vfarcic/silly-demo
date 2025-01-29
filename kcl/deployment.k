import k8s.api.apps.v1 as apps
import k8s.api.core.v1 as core

_tag = option("tag")
_probe = core.Probe { httpGet = { path = "/", port = 8080 } }

apps.Deployment {
    metadata: {
        annotations = {
            description = "This is a silly demo"
            language = "Go"
            owner = "Viktor Farcic (viktor@farcic.com)"
            team = "dot"
        }
        labels: {
            "app.kubernetes.io/name" = "silly-demo"
        }
        name = "silly-demo"
    }
    spec = {
        if not option("autoscaling.enabled"):
            replicas = option("replicas")
        selector = {
            matchLabels = metadata.labels
        }
        template = {
            metadata.labels = metadata.labels
            spec = {
                shareProcessNamespace = True
                containers = [{
                    name = "silly-demo"
                    image = "${option('image')}:${option('tag')}"
                    ports = [{ containerPort = 8080 }]
                    livenessProbe = _probe
                    readinessProbe = _probe
                    resources = {
                        limits = {
                            cpu = "500m"
                            memory = "512Mi"
                        }
                        requests = {
                            cpu = "250m"
                            memory = "256Mi"
                        }
                    }
                }]
            }
        }
    }
}
