import manifests
import .deployment
import .service
import .ingress
import .db
import .common

manifests.yaml_stream([
    deployment.Deployment,
    service.Service,
    ingress.Ingress,
    if common.Values.db.enabled:
        if common.Values.db.provider == "cnpg":
            db.CNPG,
        db.Schema,
    if common.Values.frontend.enabled:
        deployment.DeploymentFrontEnd,
        service.ServiceFrontEnd,
        ingress.IngressFrontEnd,
])