FROM golang:1.26.6-trixie

ARG UID=1000
ARG GID=1000

# 一般ユーザの作成
RUN groupadd -g ${GID} golang && \
    useradd -m -s /bin/bash -g golang golang

# OS パッケージのインストール
RUN apt-get update && apt-get install -y \
        jq \
        less \
        locales \
        vim \
        zip && \
    rm -rf /var/lib/apt/lists/*

# 日本語フォントの設定
RUN sed -ri -e "s/^# ja_JP.UTF-8/ja_JP.UTF-8/g" /etc/locale.gen && \
    locale-gen && \
    update-locale LANG=ja_JP.UTF-8 && \
    echo 'export LANG=ja_JP.utf8' >> ~/.bashrc && \
    echo 'export LANG=ja_JP.utf8' >> ~golang/.bashrc

# 作業ディレクトリの作成
RUN mkdir -p /workspaces/golang-environment && \
    chown golang:golang /workspaces/golang-environment
WORKDIR /workspaces/golang-environment



USER golang

# デバックに必要な golang パッケージのインストール
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/google/yamlfmt/cmd/yamlfmt@latest

RUN mkdir -p /home/golang/.claude \
    && curl -fsSL https://claude.ai/install.sh | bash
