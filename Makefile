SHELL := /bin/bash
SCRIPTS_PATH      := scripts

.PHONY: build-installer
build-installer:
	@$(SCRIPTS_PATH)/build_installer.sh

.PHONY: build-push-delete-air-edgex-component
build-push-delete-air-edgex-component: build-air-edgex-component push-image delete-local-image

.PHONY: build-air-edgex-component
build-air-edgex-component:
	@$(SCRIPTS_PATH)/build_air_edgex_component.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL} ${EDGEX_COMPONENT_NAME} ${TARGET_NAME}
