VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

.PHONY: all
all: build

build:
	mkdir -p _output
	go build -o _output/k8s-ldap

build_images:
	docker build -t kairen/k8s-ldap:$(VERSION) .
