FROM kairen/golang:1.11-alpine AS build-env
LABEL maintainer="Kyle Bai <k2r2.bai@gmail.com>"

ENV GOPATH "/go"
ENV PROJECT_PATH "$GOPATH/src/github.com/kairen/kube-ldap-webhook"

COPY . $PROJECT_PATH
RUN cd $PROJECT_PATH && make && \
  mv _output/kube-ldap-webhook /tmp/kube-ldap-webhook

# serve stage
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /tmp/kube-ldap-webhook /usr/bin/kube-ldap-webhook
ENTRYPOINT ["kube-ldap-webhook"]
