VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
VERSION_BUILD ?= 1
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := kairen
REPOPATH ?= $(ORG)/$(OWNER)/kube-ldap-webhook

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build: dep
	@mkdir -p _output
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	  -ldflags="-s -w" \
	  -a -o _output/kube-ldap-webhook main.go

.PHONY: dep
dep:
	@dep ensure

.PHONY: build_image
build_image:
	docker build -t $(OWNER)/kube-ldap-webhook:$(VERSION) .

.PHONY: clean
	@rm -rf _output/
