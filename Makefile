VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
VERSION_BUILD ?= 1
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := inwinstack
REPOPATH ?= $(ORG)/$(OWNER)/kube-ldap-webhook

GOOS ?= $(shell go env GOOS)

.PHONY: all
all: build

.PHONY: build
build: deps
	@mkdir -p _output
	GOOS=$(GOOS) go build -a -o _output/kube-ldap-webhook main.go

.PHONY: deps
deps:
	@dep ensure

.PHONY: build_image
build_image:
	docker build -t $(OWNER)/kube-ldap-webhook:$(VERSION) .

.PHONY: clean
	@rm -rf _output/
