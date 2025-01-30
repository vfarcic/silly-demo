import common

Ingress = {
    apiVersion = "networking.k8s.io/v1"
    kind = "Ingress"
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
