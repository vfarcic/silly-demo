package templates

import (
	corev1 "k8s.io/api/core/v1"
)

#Service: corev1.#Service & {
	_config:    #Config
	_port:      int
	if _config.isFrontend == false {
		_port: _config.service.port
	}
    if _config.isFrontend == true {
        _port: _config.frontend.service.port
    }
	apiVersion: "v1"
	kind:       "Service"
    metadata: {
        name:        _config.name
        namespace:   _config.metadata.namespace
        labels:      _config.metadata.labels
        annotations: _config.metadata.annotations
    }
	spec: corev1.#ServiceSpec & {
		type:     corev1.#ServiceTypeClusterIP
		selector: _config.selectorLabels
		ports: [
			{
				name:       "http"
				port:       _port
				targetPort: _port
				protocol:   "TCP"
			},
		]
	}
}
