[![Build Status](https://travis-ci.org/kairen/kube-ldap-webhook.svg?branch=master)](https://travis-ci.org/kairen/kube-ldap-webhook) [![Docker Build Statu](https://img.shields.io/docker/build/inwinstack/kube-ldap-webhook.svg)](https://hub.docker.com/r/inwinstack/kube-ldap-webhook/)
# Kubernetes LDAP Webhook Authentication
Kubernetes LDAP authentication service written in Go. This repo is reimplement from [kube-ldap-authn](https://github.com/torchbox/kube-ldap-authn).

According to the project documentation we have the following schema:
```sh
# kubernetesToken.schema
attributeType ( 1.3.6.1.4.1.18171.2.1.8
        NAME 'kubernetesToken'
        DESC 'Kubernetes authentication token'
        EQUALITY caseExactIA5Match
        SUBSTR caseExactIA5SubstringsMatch
        SYNTAX 1.3.6.1.4.1.1466.115.121.1.26 SINGLE-VALUE )

objectClass ( 1.3.6.1.4.1.18171.2.3
        NAME 'kubernetesAuthenticationObject'
        DESC 'Object that may authenticate to a Kubernetes cluster'
        AUXILIARY
        MUST kubernetesToken )
```

## Quick Start
In this first, modified the `ldap-auth.conf` file to match our LDAP, and then upload to Kubernetes:
```sh
$ kubectl -n kube-system create secret generic ldap-auth-config --from-file=./ldap-auth.conf
secret "ldap-auth-config" created
```

Now create the `ldap-auth-webhook` daemonset:
```sh
$ kubectl create -f deploy
```

Finally start the API server with the following flags:
```yaml
...
spec:
  containers:
  - command:
    ...
    - --runtime-config=authentication.k8s.io/v1beta1=true
    - --authentication-token-webhook-config-file=/srv/kubernetes/webhook-auth
    - --authentication-token-webhook-cache-ttl=5m
    volumeMounts:
      ...
    - mountPath: /srv/kubernetes/webhook-auth
      name: webhook-authn
      readOnly: true
  volumes:
    ...
  - hostPath:
      path: /srv/kubernetes/webhook-auth
      type: File
    name: webhook-authn
```

## How to test
In this first, we need to run gin as below command:
```sh
$ go run main.go -config=ldap-auth.conf
...
[GIN-debug] GET    /healthz                  --> main.healthz (3 handlers)
[GIN-debug] POST   /auth                     --> main.auth (3 handlers)
[GIN-debug] Listening and serving HTTP on :8087
```

To test API using cURL tool:
```sh
# failed
$ curl -X POST -H 'content-type: application/json' -d '{"apiVersion": "authentication.k8s.io/v1beta1", "kind": "TokenReview", "spec": {"token": "pSex7npm80w5Y293BNl80DeyvZL"}}' http://localhost:8087/auth
{"apiVersion":"authentication.k8s.io/v1beta1","kind":"TokenReview","status":{"authenticated":false}}

# success
$ curl -X POST -H 'content-type: application/json' -d '{"apiVersion": "authentication.k8s.io/v1beta1", "kind": "TokenReview", "spec": {"token": "pSex7npm80w5Y293BNl80DeyvZLy8Iz0"}}' http://localhost:8087/auth
{"apiVersion":"authentication.k8s.io/v1beta1","kind":"TokenReview","status":{"authenticated":true,"user":{"groups":["kubernetes","kubernetes-admin"],"uid":"user2","username":"user2"}}}
```
