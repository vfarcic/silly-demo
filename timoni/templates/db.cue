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
				db: "postgresql"
			}
		}
  		parameters: {
    		version: "13"
    		size: "small"
		}	
	}
}


#DBCNPG: {
	_config:    #Config
	apiVersion: "postgresql.cnpg.io/v1"
	kind: 		"Cluster"
	metadata: {
		name: _config.metadata.name
		namespace: _config.metadata.namespace
		labels: _config.metadata.labels
	}
	spec: {
		instances: 1
  		storage: { size: "1Gi" }
	}
}

#DBSchema: {
    _config:    #Config
    apiVersion: "db.atlasgo.io/v1alpha1"
    kind: "AtlasSchema"
    metadata: {
        name: _config.metadata.name + "-videos"
		namespace: _config.metadata.namespace
		labels: _config.metadata.labels
    }
    spec: {
        credentials: {
            scheme: "postgres"
            host: _config.metadata.name + "-rw" + "." + _config.metadata.namespace
            port: 5432
            user: "app"
            passwordFrom: {
                secretKeyRef: {
                    key: "password"
                    name: _config.metadata.name + "-app"
                }
            }
            database: "app"
            parameters: {
                sslmode: "disable"
            }
        }
        schema: {
            sql: _config.db.schema
        }
    }
}

