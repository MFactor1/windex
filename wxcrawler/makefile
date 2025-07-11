TOP_LEVEL_DIR:=$(shell git rev-parse --show-toplevel)
include $(TOP_LEVEL_DIR)/make_include/dirs.mk

DOCKER:=$(shell which docker)
GO:=$(shell which go)
GET_DEPS:=get -u
GO_CACHE:=/go/cache
GO_MOD_CACHE:=/go/pkg/mod
INSTALL:=install
BUILD:=buildx build
BUILD_FLAGS:=
RUN:=run
RUN_FLAGS:=--rm
PROJECT_NAME:=wxcrawler
PROJECT_DIR:=$(TOP_LEVEL_DIR)/$(PROJECT_NAME)
BUILD_TAG:=$(PROJECT_NAME)
UID:=$(shell id -u)
USERNAME:=$(shell whoami)

BUILD_DEPS+=github.com/gocolly/colly/v2

.PHONY: build
build:
	$(DOCKER) $(BUILD) $(BUILD_FLAGS) \
		--tag $(BUILD_TAG) \
		--build-arg UID=$(UID) \
		--build-arg USERNAME=$(USERNAME) \
		--build-arg PROJECT_MOUNT_DIR=$(REPO_MOUNT_DIR)/$(PROJECT_NAME) \
		--build-arg GO_CACHE=$(GO_CACHE) \
		--build-arg GO_MOD_CACHE=$(GO_MOD_CACHE) \
		$(PROJECT_DIR);

	$(DOCKER) $(RUN) $(RUN_FLAGS) \
		--volume $(TOP_LEVEL_DIR):$(REPO_MOUNT_DIR) \
		--volume go-mod-cache:$(GO_MOD_CACHE) \
		--volume go-cache:$(GO_CACHE) \
		--env GOMODCACHE=$(GO_MOD_CACHE) \
		--env GOCACHE=$(GO_CACHE) \
		$(BUILD_TAG);

.PHONY: run
run: build
	$(PROJECT_DIR)/build/$(PROJECT_NAME)

.PHONY: clean
clean:
	rm -rf $(PROJECT_DIR)/build

.PHONY: docker-build
docker-build:
	$(GO) $(GET_DEPS) $(BUILD_DEPS);

	$(GO) $(INSTALL);
