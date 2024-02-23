package templates

import (
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

dbProvider: "local-k8s" | "aws-official" | "azure-official" | "cnpg"

#Config: {
	metadata: metav1.#ObjectMeta
	metadata: name:         *"silly-demo" | string & =~"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" & strings.MaxRunes(63)
	metadata: namespace:    *"default" | string & strings.MaxRunes(63)
	metadata: labels:       *selectorLabels | {[ string]: string}
	metadata: annotations: {
		"description": "This is a silly demo",
		"owner": "Viktor Farcic (viktor@farcic.com)",
		"team": "dot",
		"language": "Go",
	}

	replicas:        *2 | int & >0
	selectorLabels:  *{"app.kubernetes.io/name": metadata.name} | {[ string]: string}
	podAnnotations?: {[ string]: string}
	image: {
		repository: *"ghcr.io/vfarcic/silly-demo" | string
		tag:        *"latest" | string
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
	service: {
		port:       *8080 | int & >0 & <=65535
		targetPort: *8080 | int & >0 & <=65535
	}
	autoscaling: {
		enabled:     *false | bool
		cpu:         *80 | int & >0 & <=100
		memory:      *80 | int & >0 & <=100
		minReplicas: *replicas | int
		maxReplicas: *6 | int & >=minReplicas
	}
	ingress: {
		host:    *"sillydemo.com" | string
		className?: string
	}
	db: {
		enabled: *false | bool
		provider: *"google-official" | dbProvider
		type: *"postgres" | string
	}
}

#Instance: {
	config: #Config

	objects: {
		"\(config.metadata.name)-deploy": #Deployment & {_config:     config}
		"\(config.metadata.name)-svc":    #Service & {_config:        config}
		if config.autoscaling.enabled {
			"\(config.metadata.name)-hpa": #HorizontalPodAutoscaler & {_config: config}
		}
		"\(config.metadata.name)-ingress": #Ingress & {_config: config}
		if config.db.enabled {
			if config.db.provider == "cnpg" {
				"\(config.metadata.name)-db-cnpg": #DBCNPG & {_config: config}
			}
			if config.db.provider != "cnpg" {
				"\(config.metadata.name)-db-secret": #DBSecret & {_config: config}
				"\(config.metadata.name)-db-claim": #DBClaim & {_config: config}
			}
		}
	}
}
