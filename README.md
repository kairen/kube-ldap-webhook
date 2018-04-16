# Kubernetes LDAP Webhook Authentication
Kubernetes LDAP authentication service written in Go.
(TBD..)

# Test
In this first, we need to run gin as below command:
```sh
$ go run main.go
...
[GIN-debug] GET    /healthz                  --> main.healthz (3 handlers)
[GIN-debug] POST   /auth                     --> main.auth (3 handlers)
[GIN-debug] Listening and serving HTTP on :8087
```

To test API using cURL tool:
```sh
# failed
$ curl -X POST -H 'content-type: application/json' -d '{"apiVersion": "authentication.k8s.io/v1beta1", "kind": "TokenReview", "spec": {"token": "test1"}}' http://localhost:8087/auth
{"apiVersion":"authentication.k8s.io/v1beta1","kind":"TokenReview","status":{"authenticated":false}}

# success
$ curl -X POST -H 'content-type: application/json' -d '{"apiVersion": "authentication.k8s.io/v1beta1", "kind": "TokenReview", "spec": {"token": "test"}}' http://localhost:8087/auth
{"apiVersion":"authentication.k8s.io/v1beta1","kind":"TokenReview","status":{"authenticated":true,"user":{"groups":"kubernetes-admin","uid":"test","username":"test"}}}
```
