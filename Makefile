SHELL := /bin/bash
SCRIPTS_PATH      := scripts

ifndef IMAGE_NAME
override IMAGE_NAME = labs-air-edgex
endif
ifndef IMAGE_TAG
override IMAGE_TAG = latest
endif
ifndef ECR_REGISTRY
override ECR_REGISTRY = public.ecr.aws
endif
ifndef ECR_REPO_NAME
override ECR_REPO_NAME = tibcolabs
endif
ifndef IMAGE_URL
override IMAGE_URL = "$(ECR_REGISTRY)/$(ECR_REPO_NAME)"
endif

.PHONY: build-installer
build-installer:
	@$(SCRIPTS_PATH)/build_installer.sh

.PHONY: build-push-delete-air-edgex-component
build-push-delete-air-edgex-component: build-air-edgex-component push-image delete-local-image

.PHONY: build-air-edgex-component
build-air-edgex-component:
	@$(SCRIPTS_PATH)/build_air_edgex_component.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL} ${EDGEX_COMPONENT_NAME} ${TARGET_NAME}

.PHONY: push-image
push-image:
	@$(SCRIPTS_PATH)/push_image.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL}

.PHONY: delete-local-image
delete-local-image:
	@$(SCRIPTS_PATH)/delete_local_image.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL}
