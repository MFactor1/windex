IMAGE_NAME=mongo:7.0
CONTAINER_NAME=mongodb_server
DATA_DIR=$(shell pwd)/data

MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=example

run:
	@mkdir -p $(DATA_DIR)
	docker run -d --rm \
		--name $(CONTAINER_NAME) \
		-p 27017:27017 \
		-v $(DATA_DIR):/data/db \
		-e MONGO_INITDB_ROOT_USERNAME=$(MONGO_INITDB_ROOT_USERNAME) \
		-e MONGO_INITDB_ROOT_PASSWORD=$(MONGO_INITDB_ROOT_PASSWORD) \
		$(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME)

logs:
	docker logs -f $(CONTAINER_NAME)

status:
	docker ps -f name=$(CONTAINER_NAME)

