VERSION ?= v0.4.0
OPERATOR_IMAGE = virtualenvironment/virtual-env-operator
WEBHOOK_IMAGE = virtualenvironment/virtual-env-admission-webhook
OPERATOR_IMAGE_AND_VERSION ?= $(OPERATOR_IMAGE):$(VERSION)
WEBHOOK_IMAGE_AND_VERSION ?= $(WEBHOOK_IMAGE):$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo 'use "make build-operator" or "make build-webhook" to build images'

.PHONY: build-operator
build-operator:
	operator-sdk build \
		--go-build-args "-ldflags -X=alibaba.com/virtual-env-operator/version.BuildTime=`date +%Y-%m-%d_%H:%M` -o build/_output/operator/virtual-env-operator" \
		--image-build-args "--no-cache" $(OPERATOR_IMAGE_AND_VERSION)

.PHONY: build-webhook
build-webhook: $(shell find cmd/webhook -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X=main.buildTime=`date +%Y-%m-%d_%H:%M`" \
		-o "build/_output/webhook/webhook-server" ./cmd/webhook
	docker build -t $(WEBHOOK_IMAGE_AND_VERSION) -f build/Dockerfile_webhook build/_output/webhook/

.PHONY: clean
clean:
	rm -fr build/_output/
	for i in `docker images | grep $(OPERATOR_IMAGE) | awk '{print $$3}'`; do docker rmi -f $$i; done
	for i in `docker images | grep $(WEBHOOK_IMAGE) | awk '{print $$3}'`; do docker rmi -f $$i; done
