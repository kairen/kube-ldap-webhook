FROM golang:1.9-alpine AS build-env
LABEL maintainer="Kyle Bai <kyle.b@inwinstack.com>"

ENV GOPATH "/go"

RUN apk add --no-cache git make g++ && \
  go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/app
RUN cd /go/src/app && make

# serve stage
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/app/_output/kube-ldap-webhook /usr/bin/kube-ldap-webhook
ENTRYPOINT ["kube-ldap-webhook"]
