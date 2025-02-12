import k8s.api.apps.v1 as apps
import k8s.api.core.v1 as core
import common
import db

_tag = common.Values.tag
_probe = core.Probe { httpGet = { path = "/", port = 8080 } }

Deployment = apps.Deployment {
    metadata = common.Metadata
    spec = {
        if not common.Values.autoscaling.enabled:
            replicas = common.Values.replicas
        selector = {
            matchLabels = metadata.labels
        }
        template = {
            metadata.labels = metadata.labels
            spec = {
                shareProcessNamespace = True
                containers = [{
                    name = "main"
                    image = "${common.Values.image}:${common.Values.tag}"
                    ports = [{ containerPort = 8080 }]
                    livenessProbe = _probe
                    readinessProbe = _probe
                    resources = {
                        limits = {
                            cpu = "250m"
                            memory = "256Mi"
                        }
                        requests = {
                            cpu = "125m"
                            memory = "128Mi"
                        }
                    }
                    if common.Values.db.enabled:
                        env = [{
                            name = "DB_URI"
                            valueFrom.secretKeyRef = { name = db.SecretName, key = "uri" }
                        }]
                }]
            }
        }
    }
}
