package templates

import (
	netv1 "k8s.io/api/networking/v1"
)

#Ingress: netv1.#Ingress & {
	_config:      #Config
	_ingressHost: string
	_port:		  int
	if _config.isFrontend == false {
		_ingressHost: _config.ingress.host
		_port: _config.service.port
	}
	if _config.isFrontend == true {
		_ingressHost: _config.frontend.ingress.host
		_port: _config.frontend.service.port
	}
	apiVersion: "networking.k8s.io/v1"
	kind:       "Ingress"
    metadata: {
        name:        _config.name
        namespace:   _config.metadata.namespace
        labels:      _config.metadata.labels
        annotations: _config.metadata.annotations
    }
	if _config.ingress.annotations != _|_ {
		metadata: annotations: _config.ingress.annotations
	}
	spec: netv1.#IngressSpec & {
		rules: [{
			host: _ingressHost
			http: {
				paths: [{
					pathType: "ImplementationSpecific"
					path:     "/"
					backend: service: {
						name: _config.name
						port: number: _port
					}
				}]
			}
		}]
		if _config.ingress.className != _|_ {
			ingressClassName: _config.ingress.className
		}
	}
}
