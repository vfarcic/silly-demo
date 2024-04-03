package templates

import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
)

#Deployment: appsv1.#Deployment & {
    _config:    #Config
    _secretName: string
    if _config.db.provider == "cnpg" {
        _secretName: _config.metadata.name + "-app"
    }
    if _config.db.provider != "cnpg" {
        _secretName: _config.metadata.name
    }
    apiVersion: "apps/v1"
    kind:       "Deployment"
    metadata:   _config.metadata
    spec:       appsv1.#DeploymentSpec & {
        if !_config.autoscaling.enabled {
            replicas: _config.replicas
        }
        selector: matchLabels: _config.selectorLabels
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
                    name: _config.metadata.name
                    image: "\(_config.image.repository):\(_config.image.tag)"
                    ports: [ { containerPort: _config.service.targetPort } ]
                    livenessProbe: {
                        httpGet: {
                            path: "/"
                            port: _config.service.targetPort
                        }
                    }
                    readinessProbe: {
                        httpGet: {
                            path: "/"
                            port: _config.service.targetPort
                        }
                    }
                    if _config.resources != _|_ {
                        resources: _config.resources
                    }
                    if _config.db.enabled == true {
                        env: [ {
                            name: "DB_ENDPOINT"
                            valueFrom: {
                                secretKeyRef: {
                                    name: _secretName
                                    if _config.db.provider == "cnpg" {
                                        key: "host"
                                    }
                                    if _config.db.provider != "cnpg" {
                                        key: "endpoint"
                                    }
                                }
                            }
                        }, {
                            name: "DB_PORT"
                            valueFrom: {
                                secretKeyRef: {
                                    name: _secretName
                                    key: "port"
                                }
                            }
                        }, {
                            name: "DB_USER"
                            valueFrom: {
                                secretKeyRef: {
                                    name: _secretName
                                    key: "username"
                                }
                            }
                        }, {
                            name: "DB_PASS"
                            valueFrom: {
                                secretKeyRef: {
                                    name: _secretName
                                    key: "password"
                                }
                            }
                        }, {
                            name: "DB_NAME"
                            if _config.db.provider == "cnpg" {
                                value: "app"
                            }
                            if _config.db.provider != "cnpg" {
                                value: _config.metadata.name
                            }
                        }]
                    }
                },
                if _config.otel.enabled == true {
                    {
                        name: _config.metadata.name + "-instrumentation"
                        image: "otel/autoinstrumentation-go"
                        env: [{
                            name: "OTEL_GO_AUTO_TARGET_EXE"
                            value: "/usr/local/bin/silly-demo"
                        }, {
                            name: "OTEL_EXPORTER_OTLP_ENDPOINT"
                            value: _config.otel.jaegerAddr
                        }, {
                            name: "OTEL_SERVICE_NAME"
                            value: _config.metadata.name
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
