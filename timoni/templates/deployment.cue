package templates

import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
)

#Deployment: appsv1.#Deployment & {
    _config:     #Config
    _secretName: string
    _image:      string
    _port:       int
    _env:        [...]
    if _config.isFrontend == false {
        _image: _config.image.repository + ":" + _config.image.tag
        _port: _config.service.port
        if (_config.db.enabled == true || _config.debug.enabled == true) {
            _env: [
                if _config.db.enabled == true {
                    {
                        name: "DB_ENDPOINT"
                        valueFrom: secretKeyRef: {
                            name: _secretName
                            if _config.db.provider == "cnpg" {
                                key: "host"
                            }
                            if _config.db.provider != "cnpg" {
                                key: "endpoint"
                            }
                        }
                    }
                }
                if _config.db.enabled == true {
                    {
                        name: "DB_PORT"
                        valueFrom: secretKeyRef: {
                            name: _secretName
                            key: "port"
                        }
                    }
                }
                if _config.db.enabled == true {
                    {
                        name: "DB_USER"
                        valueFrom: secretKeyRef: {
                            name: _secretName
                            key: "username"
                        }
                    }
                }
                if _config.db.enabled == true {
                    {
                        name: "DB_PASS"
                        valueFrom: secretKeyRef: {
                            name: _secretName
                            key: "password"
                        }
                    }
                }
                if _config.db.enabled == true {
                    {
                        name: "DB_NAME"
                        if _config.db.provider == "cnpg" {
                            value: "app"
                        }
                        if _config.db.provider != "cnpg" {
                            value: _config.name
                        }
                    }
                }
                if _config.debug.enabled == true {
                    {name: "DEBUG", value: "true" }
                }
            ]
        }
    }
    if _config.isFrontend == true {
        _image: _config.frontend.image.repository + ":" + _config.frontend.image.tag
        _port: _config.frontend.service.port
        _env: [
            {
                name: "REACT_APP_BACKEND_URL"
                value: "http://" + _config.ingress.host
            }
        ]
    }
    if _config.db.provider == "cnpg" {
        _secretName: _config.name + "-app"
    }
    if _config.db.provider != "cnpg" {
        _secretName: _config.name
    }

    apiVersion: "apps/v1"
    kind:       "Deployment"
    metadata: {
        name:        _config.name
        namespace:   _config.metadata.namespace
        labels:      _config.metadata.labels
        annotations: _config.metadata.annotations
    }
    spec:       appsv1.#DeploymentSpec & {
        if !_config.autoscaling.enabled {
            replicas: _config.replicas
        }
        selector: matchLabels: {"app.kubernetes.io/name": _config.name}
        template: {
            metadata: {
                labels: _config.selectorLabels
                if _config.podAnnotations != _|_ {
                    annotations: _config.podAnnotations
                }
            }
            spec: corev1.#PodSpec & {
                if _config.otel.enabled == true {
                    shareProcessNamespace: true
                }
                containers: [{
                    name: _config.name
                    image: _image
                    ports: [ { containerPort: _port } ]
                    livenessProbe: {
                        httpGet: {
                            path: "/"
                            port: _port
                        }
                    }
                    readinessProbe: {
                        httpGet: {
                            path: "/"
                            port: _port
                        }
                    }
                    if _config.resources != _|_ {
                        resources: _config.resources
                    }
                    env: _env
                },
                if _config.otel.enabled == true {
                    {
                        name: _config.name + "-instrumentation"
                        image: "otel/autoinstrumentation-go"
                        env: [{
                            name: "OTEL_GO_AUTO_TARGET_EXE"
                            value: "/usr/local/bin/silly-demo"
                        }, {
                            name: "OTEL_EXPORTER_OTLP_ENDPOINT"
                            value: _config.otel.jaegerAddr
                        }, {
                            name: "OTEL_SERVICE_NAME"
                            value: _config.name
                        }, {
                            name: "OTEL_PROPAGATORS"
                            value: "tracecontext,baggage"
                        }]
                        securityContext: {
                            runAsUser: 0
                            privileged: true
                        }
                    }
                }
                ]
            }
        }
    }
}
