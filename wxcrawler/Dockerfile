# syntax=docker/dockerfile:1.4
FROM ubuntu:24.04

ARG UID
ARG GID
ARG USERNAME
ARG PROJECT_MOUNT_DIR
ARG GO_CACHE
ARG GO_MOD_CACHE

ENV USERNAME=$USERNAME

RUN useradd --create-home --shell /bin/bash "$USERNAME" --uid "$UID" --non-unique

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
	&& apt-get install -y \
	ca-certificates \
	curl \
	make \
	git \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION=1.24.4

RUN curl -sSL https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

RUN mkdir -p ${GO_CACHE} && chown -R ${USERNAME} ${GO_CACHE} \
	&& mkdir -p ${GO_MOD_CACHE} && chown -R ${USERNAME} ${GO_MOD_CACHE}

WORKDIR ${PROJECT_MOUNT_DIR}

ENV GOBIN=${PROJECT_MOUNT_DIR}/build

USER ${USERNAME}

CMD ["make", "docker-build"]
