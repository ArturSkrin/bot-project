REGISTRY   ?= ghcr.io
REPOSITORY ?= arturskrin/bot-project
VERSION    ?= v1.0.0
OS         ?= linux
ARCH       ?= amd64
SHA        := $(shell git rev-parse --short HEAD)
TAG        := $(VERSION)-$(SHA)
IMAGE      := $(REGISTRY)/$(REPOSITORY):$(TAG)-$(OS)-$(ARCH)

.PHONY: build push helm-set image

image:
	@echo $(IMAGE)

build:
	docker build --platform $(OS)/$(ARCH) -t $(IMAGE) .

push:
	docker build --platform $(OS)/$(ARCH) -t $(IMAGE) . 
	docker push $(IMAGE)

helm-set:
	yq -i '.image.registry   = "$(REGISTRY)"'   bot-project/values.yaml
	yq -i '.image.repository  = "$(REPOSITORY)"' bot-project/values.yaml
	yq -i '.image.tag         = "$(TAG)"'        bot-project/values.yaml
	yq -i '.image.os          = "$(OS)"'         bot-project/values.yaml
	yq -i '.image.arch        = "$(ARCH)"'       bot-project/values.yaml