TOP_LEVEL_DIR:=$(shell git rev-parse --show-toplevel)
DOCKER:=$(shell which docker)
BUILD:=buildx build
BUILD_FLAGS:=
RUN:=run
RUN_FLAGS:=--rm
PROJECT_NAME:=wxcrawler
PROJECT_DIR:=$(TOP_LEVEL_DIR)/$(PROJECT_NAME)
BUILD_TAG:=$(PROJECT_NAME)
UID:=$(shell id -u)
USERNAME:=$(shell whoami)

.PHONY: build
build:
	$(DOCKER) $(BUILD) $(BUILD_FLAGS) \
		--tag $(BUILD_TAG) \
		--build-arg UID=$(UID) \
		--build-arg USERNAME=$(USERNAME) \
		$(PROJECT_DIR);

	$(DOCKER) $(RUN) $(RUN_FLAGS) \
		--volume $(PROJECT_DIR):/mnt/repo/$(PROJECT_NAME) \
		$(BUILD_TAG);

.PHONY: clean
clean:
	rm -rf $(PROJECT_DIR)/build
