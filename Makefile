COMPOSE ?= docker compose

.PHONY: help
help:
	@echo "Available targets:"
	@awk 'BEGIN {FS=":.*## "}; /^[[:alnum:]_.-]+:.*## / {printf "  %-24s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: install build up up-rebuild down restart ps logs \
		shell

install: ## 初回セットアップ
	cp .env.example .env
	@make build
	@make up

build: ## すべてのサービスをビルド
	$(COMPOSE) build

up: ## バックグラウンドで起動
	$(COMPOSE) up -d

up-rebuild: ## 再ビルドして起動（Dockerfile/依存更新があるとき）
	$(COMPOSE) up -d --build

down: ## 停止（ボリュームは残す）
	$(COMPOSE) down

restart: ## 再起動
	$(COMPOSE) down && $(COMPOSE) up -d

ps: ## 稼働状況を表示
	$(COMPOSE) ps

app-logs: ## app サービスのログ
	$(COMPOSE) logs -f app

shell: ## app コンテナに Bash で入る
	$(COMPOSE) exec app bash