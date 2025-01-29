import k8s.api.networking.v1 as net
import common

Ingress = net.Ingress {
    metadata = common.Metadata
    spec = {
        ingressClassName = common.Values.ingress.className
        rules: [{
            host = common.Values.ingress.host
            http.paths: [{
                path = "/"
                pathType = "ImplementationSpecific"
                backend.service = {
                    name = common.Name
                    port.number = 8080
                }
            }]
        }]
    }
}
