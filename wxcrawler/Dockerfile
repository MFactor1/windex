FROM ubuntu:24.04

ARG UID
ARG GID
ARG USERNAME

ENV USERNAME=$USERNAME

RUN useradd --create-home --shell /bin/bash "$USERNAME" --uid "$UID" --non-unique

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
	&& apt-get install -y \
	ca-certificates \
	curl \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION=1.24.4

RUN curl -sSL https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

WORKDIR /mnt/repo/wxcrawler

ENV GOBIN=/mnt/repo/wxcrawler/build

USER ${USERNAME}

CMD ["go", "install"]
