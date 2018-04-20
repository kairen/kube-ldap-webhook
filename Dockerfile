FROM golang:1.9-alpine AS build-env
MAINTAINER Kyle Bai(k2r2bai@gmail.com)

ENV GOPATH "/go"

RUN apk add --no-cache git && go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/app
RUN cd /go/src/app && dep ensure && go build -o k8s-ldap

# serve stage
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/app/k8s-ldap /usr/bin/k8s-ldap
ENTRYPOINT ["k8s-ldap"]
