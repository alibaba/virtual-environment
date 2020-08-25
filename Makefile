VERSION ?= v0.3.2
OPERATOR_IMAGE ?= virtualenvironment/virtual-env-operator
ADMISSION_IMAGE ?= virtualenvironment/virtual-env-admission-webhook

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo 'use "make build-operator" or "make build-admission" to build images'

.PHONY: build-operator
build-operator:
	operator-sdk build \
		--go-build-args "-ldflags -X=alibaba.com/virtual-env-operator/version.BuildTime=`date +%Y-%m-%d_%H:%M` -o build/_output/operator/virtual-env-operator" \
		--image-build-args "--no-cache" $(OPERATOR_IMAGE):$(VERSION)

.PHONY: build-admission
build-admission: $(shell find cmd/webhook -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o "build/_output/admission/webhook-server" ./cmd/webhook
	docker build -t $(ADMISSION_IMAGE):$(VERSION) -f build/Dockerfile_webhook build/_output/admission/

.PHONY: clean
clean:
	rm -fr build/_output/
	docker rmi -f $(OPERATOR_IMAGE):$(VERSION)
	docker rmi -f $(ADMISSION_IMAGE):$(VERSION)

.PHONY: clean-image
clean-image:
	for i in `docker images | grep $(OPERATOR_IMAGE) | awk '{print $$3}'`; do docker rmi -f $$i; done
	for i in `docker images | grep $(ADMISSION_IMAGE) | awk '{print $$3}'`; do docker rmi -f $$i; done
