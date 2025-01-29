import k8s.api.core.v1 as core
import common

Service = core.Service {
    metadata = common.Metadata
    spec = {
        ports = [{
            name = "http"
            port = 8080
            protocol = "TCP"
            targetPort = 8080
        }]
        selector = metadata.labels
        type = "ClusterIP"
    }
}
