import k8s.apimachinery.pkg.apis.meta.v1 as meta
import file
import yaml

Name = "silly-demo"

_values = yaml.decode(file.read("kcl/values.yaml"))
if option("tag"):
    _values.tag = option("tag")
if option("db.enabled"):
    _values.db.enabled = option("db.enabled")
if option("db.provider"):
    _values.db.provider = option("db.provider")
if option("frontend.enabled"):
    _values.frontend.enabled = option("frontend.enabled")
if option("frontend.host"):
    _values.frontend.host = option("frontend.host")

Values = _values

Metadata = meta.ObjectMeta {
    annotations = {
        description = "This is a silly demo"
        language = "Go"
        owner = "Viktor Farcic (viktor@farcic.com)"
        team = "dot"
    }
    labels: {
        "app.kubernetes.io/name" = Name
    }
    name = Name
}

MetadataFrontEnd = meta.ObjectMeta {
    annotations = {
        description = "This is a silly demo frontend"
        language = "TypeScript"
        owner = "Viktor Farcic (viktor@farcic.com)"
        team = "dot"
    }
    labels: {
        "app.kubernetes.io/name" = Name + "-frontend"
    }
    name = Name + "-frontend"
}
