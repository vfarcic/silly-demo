package templates

import (
	corev1 "k8s.io/api/core/v1"
)

#DBSecret: corev1.#Secret & {
    _config:    #Config
	apiVersion: "v1"
	kind:       "Secret"
	metadata: {
		name: _config.metadata.name + "-password"
		namespace: _config.metadata.namespace
		labels: _config.metadata.labels
		annotations: {
			"argocd.argoproj.io/sync-wave": "-10"
		}
	}
	data: {
		password: 'cG9zdGdyZXM='
	}
}

#DBClaim: {
	_config:    #Config
	apiVersion: "devopstoolkitseries.com/v1alpha1"
	kind: 		"SQLClaim"
	metadata: {
		name: _config.metadata.name
		namespace: _config.metadata.namespace
		labels: _config.metadata.labels
		annotations: {
			"argocd.argoproj.io/sync-wave": "-10"
		}
	}
	spec: {
		id: _config.metadata.name
  		compositionSelector: {
			matchLabels: {
				provider: _config.db.provider
				db: "postgres"
			}
		}
  		parameters: {
    		version: "13"
    		size: "small"
		}	
	}
}
