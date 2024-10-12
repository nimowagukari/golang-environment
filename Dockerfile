FROM golang:1.21.7-bookworm

ARG UID=1000
ARG GID=1000

# 一般ユーザの作成
RUN groupadd -g ${GID} app && \
    useradd -m -s /bin/bash -g app app

# OS パッケージのインストール
RUN apt-get update && apt-get install -y \
        less \
        vim \
        zip && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /workspaces
USER app

# デバックに必要な golang パッケージのインストール
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest
