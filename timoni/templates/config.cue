package templates

import (
    "strings"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    corev1 "k8s.io/api/core/v1"
)

dbProvider: "aws" | "azure" | "google" | "cnpg"

#Config: {
    isFrontend: *false | bool
    name: *"silly-demo" | string & =~"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" & strings.MaxRunes(63)
    language: *"Go" | string
    metadata: metav1.#ObjectMeta
    metadata: {
        name:         name
        namespace:    *"default" | string & strings.MaxRunes(63)
        labels:       *selectorLabels | {[ string]: string}
        annotations: {
            "description": "This is a silly demo",
            "owner": "Viktor Farcic (viktor@farcic.com)",
            "team": "dot",
            "language": language,
        }
    }

    replicas:        *2 | int & >0
    selectorLabels:  *{"app.kubernetes.io/name": name} | {[string]: string}
    podAnnotations?: {[ string]: string}
    image: {
        repository: string
        tag:        string
    }
    _resources: {
        limits: {
            cpu:    "500m"
            memory: "512Mi"
        }
        requests: {
            cpu:    "250m"
            memory: "256Mi"
        }
    }
    resources: *_resources | corev1.#ResourceRequirements
    service: port: *8080 | int & >0 & <=65535
    autoscaling: {
        enabled:     *false | bool
        cpu:         *80 | int & >0 & <=100
        memory:      *80 | int & >0 & <=100
        minReplicas: *replicas | int
        maxReplicas: *6 | int & >=minReplicas
    }
    ingress: {
        host:    *"silly-demo.com" | string
        className?: string
    }
    db: {
        enabled: *false | bool
        provider: *"google" | dbProvider
        type: *"postgres" | string
        schema: *"" | string
    }
    otel: {
        enabled: *false | bool
        jaegerAddr: *"http://jaeger.kube-system:4318" | string
    }
    debug: {
        enabled: *false | bool
    }
    frontend: {
        enabled: *true | bool
        image: {
            repository: string
            tag:        string
        }
        service: port: 3000
        ingress: host: *"silly-demo-frontend.com" | string
    }
}

#Instance: {
    config: #Config

    objects: {
        "\(config.metadata.name)-deploy": #Deployment & { _config: config }
        "\(config.metadata.name)-svc": #Service & { _config: config }
        if config.autoscaling.enabled {
            "\(config.metadata.name)-hpa": #HorizontalPodAutoscaler & {_config: config}
        }
        "\(config.metadata.name)-ingress": #Ingress & {_config: config}
        if config.db.enabled {
            if config.db.provider == "cnpg" {
                "\(config.metadata.name)-db-cnpg": #DBCNPG & {_config: config}
                if config.db.schema != "" {
                    "\(config.metadata.name)-db-schema": #DBSchema & {_config: config}
                }
            }
            if config.db.provider != "cnpg" {
                "\(config.metadata.name)-db-secret": #DBSecret & {_config: config}
                "\(config.metadata.name)-db-claim": #DBClaim & {_config: config}
            }
        }
        if config.frontend.enabled {
            "\(config.metadata.name)-deploy-frontend": #Deployment & {
                _config: config
                _config: {
                    isFrontend: true
                    name: "silly-demo-frontend",
                    language: "Node",
                }
            }
            "\(config.metadata.name)-svc-frontend": #Service & {
                _config: config
                _config: {
                    isFrontend: true
                    name: "silly-demo-frontend",
                    language: "Node",
                }
            }
            "\(config.metadata.name)-ingress-frontend": #Ingress & {
                _config: config
                _config: {
                    isFrontend: true
                    name: "silly-demo-frontend",
                    language: "Node",
                }
            }
        }
    }
}
