package templates

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

#Deployment: appsv1.#Deployment & {
	_config:    #Config
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
				containers: [
					{
						name: _config.metadata.name
						image: "\(_config.image.repository):\(_config.image.tag)"
						ports: [
							{
								containerPort: _config.service.targetPort
							},
						]
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
					},
				]
			}
		}
	}
}
