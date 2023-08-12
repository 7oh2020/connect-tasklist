FROM mcr.microsoft.com/devcontainers/go:0-1.20-bullseye

# モック生成のためにvektra/mockeryをインストール
RUN curl -L https://github.com/vektra/mockery/releases/download/v2.25.0/mockery_2.25.0_Linux_arm64.tar.gz | tar xvz && mv ./mockery /usr/bin/mockery

# データベースのマイグレーションのためにgolang-migrate/migrateをインストール
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-arm64.tar.gz | tar xvz && mv ./migrate /usr/bin/migrate

# サーバーのライブリロードのためにcosmtrek/airをインストール
RUN go install github.com/cosmtrek/air@latest

# データベースアクセスのためにkyleconroy/sqlcをインストール
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

# gRPCとREST APIのコード生成のためにconnect-goとbufをインストール
RUN go install github.com/bufbuild/buf/cmd/buf@latest \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
  && go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# モジュールへ実行権限を与える(vscodeユーザーでもアクセス可能にする)
RUN sudo chmod -R a+rwX /go/pkg
