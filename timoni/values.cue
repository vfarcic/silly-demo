// Code generated by timoni. DO NOT EDIT.
// Note that this file must have no imports and all values must be concrete.

package main

values: {
    image: tag: "1.4.342"
    image: repository: "ghcr.io/vfarcic/silly-demo"
    replicas: 2
    autoscaling: {
        enabled:     false
        cpu:         80
        memory:      80
        maxReplicas: 6
    }
    ingress: host: "silly-demo.com"
    db: {
        enabled:  false
        provider: "google"
        type:     "postgres"
        schema:   ""
    }
    otel: {
        enabled:    false
        jaegerAddr: "http://jaeger.kube-system:4318"
    }
	frontend: {
        image: {
            tag: "0.0.3"
	        repository: "ghcr.io/vfarcic/silly-demo-frontend"
        }
        ingress: host: "silly-demo-frontend.com"
    }
    debug: enabled: false
}
