TOP_LEVEL_DIR:=$(shell git rev-parse --show-toplevel)
include $(TOP_LEVEL_DIR)/make_include/dirs.mk

PROJECT_DIR := $(TOP_LEVEL_DIR)/wxindexer/redis
REDIS_DATA := $(PROJECT_DIR)/localdata
CONTAINER_NAME := wxredis
REDIS_PORT := 6380
REDIS_IMAGE := redis:8

.PHONY: redis-start redis-stop redis-clean redis-status

start:
	@if [ $$(docker ps -a -q -f name=$(CONTAINER_NAME)) ]; then \
		echo "Container '$(CONTAINER_NAME)' exists. Starting..."; \
		docker start $(CONTAINER_NAME); \
	else \
		echo "Container '$(CONTAINER_NAME)' does not exist. Creating and starting..."; \
		mkdir -p $(REDIS_DATA); \
		docker run -d \
			--name $(CONTAINER_NAME) \
			-p $(REDIS_PORT):6379 \
			-v $(REDIS_DATA):/data \
			$(REDIS_IMAGE) \
			redis-server --appendonly yes; \
	fi

stop:
	docker stop $(CONTAINER_NAME) || true

clean: stop
	docker rm $(CONTAINER_NAME) || true
	sudo rm -rf $(REDIS_DATA)

status:
	docker ps -f name=$(CONTAINER_NAME)

