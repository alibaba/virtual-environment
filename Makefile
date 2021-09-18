VERSION ?= v0.6.0
OPERATOR_IMAGE = virtualenvironment/virtual-env-operator
WEBHOOK_IMAGE = virtualenvironment/virtual-env-admission-webhook
OPERATOR_IMAGE_AND_VERSION ?= $(OPERATOR_IMAGE):$(VERSION)
WEBHOOK_IMAGE_AND_VERSION ?= $(WEBHOOK_IMAGE):$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo 'use "make build-operator" or "make build-webhook" to build images'

.PHONY: build-inspector-binary
build-inspector-binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w" \
		-o "build/_output/operator/inspector" ./cmd/inspector

.PHONY: build-operator-binary
build-operator-binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a \
		-ldflags "-s -w -X=alibaba.com/virtual-env-operator/version.BuildTime=`date +%Y-%m-%d_%H:%M` -X=alibaba.com/virtual-env-operator/version.Version=$(VERSION)" \
		-o "build/_output/operator/virtual-env-operator" ./cmd/operator

.PHONY: build-operator-image
build-operator-image:
	docker build --no-cache -t $(OPERATOR_IMAGE_AND_VERSION) -f build/Dockerfile_operator build/_output/operator/

.PHONY: build-operator
build-operator: build-operator-binary build-inspector-binary build-operator-image

.PHONY: build-webhook-binary
build-webhook-binary: $(shell find cmd/webhook -name '*.go')
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a \
		-ldflags="-s -w -X=main.buildTime=`date +%Y-%m-%d_%H:%M` -X=main.version=$(VERSION)" \
		-o "build/_output/webhook/webhook-server" ./cmd/webhook

.PHONY: build-webhook-image
build-webhook-image:
	docker build --no-cache -t $(WEBHOOK_IMAGE_AND_VERSION) -f build/Dockerfile_webhook build/_output/webhook/

.PHONY: build-webhook
build-webhook: build-webhook-binary build-webhook-image

.PHONY: push
push:
	docker push $(OPERATOR_IMAGE_AND_VERSION)
	docker push $(WEBHOOK_IMAGE_AND_VERSION)

.PHONY: clean
clean:
	rm -fr build/_output/
	rm -f examples/go/main
	rm -fr examples/springboot/target/
	rm -fr sdk/java/*/target
	for i in `docker images | grep $(OPERATOR_IMAGE) | awk '{print $$3}'`; do docker rmi -f $$i; done
	for i in `docker images | grep $(WEBHOOK_IMAGE) | awk '{print $$3}'`; do docker rmi -f $$i; done
