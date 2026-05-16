.DEFAULT_GOAL := help

ifneq (,$(wildcard .env))
include .env
export
endif

GOOSE ?= goose
GOOSE_DRIVER ?= postgres
DEFAULT_DATABASE_URL ?= postgres://zennify:zennify@localhost:5432/zennify?sslmode=disable

.PHONY: help fmt lint lint-fix test buf-lint buf migrate-up migrate-down create-mig migrate-status

help:
	@echo "make migrate-up    SERVICE=user [STEPS=1] [TO=version]  — применить миграции"
	@echo "make migrate-down  SERVICE=user [STEPS=1] [TO=version]  — откат (по умолчанию 1 шаг)"
	@echo "make create-mig    SERVICE=user NAME=add_email          — создать SQL-миграцию (goose timestamp)"
	@echo "make migrate-status SERVICE=user                        — статус миграций"
	@echo ""
	@echo "Сервисы: каталог internal/<SERVICE>/store/migrations (сейчас: user)"
	@echo "БД: <SERVICE>_DATABASE_URL в .env, иначе DEFAULT_DATABASE_URL"

fmt:
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./...

buf-lint:
	buf lint

buf:
	buf lint & buf generate

migrate-up:
	@$(MAKE) _goose-up

migrate-down:
	@$(MAKE) _goose-down

create-mig:
	@$(MAKE) _goose-create

migrate-status:
	@$(MAKE) _goose-status

_goose-up:
	@test -n "$(SERVICE)" || (echo "migrate-up: укажите SERVICE=..., например SERVICE=user"; exit 1)
	@test -d "internal/$(SERVICE)/store/migrations" || (echo "migrate-up: нет internal/$(SERVICE)/store/migrations"; exit 1)
	@set -a; [ -f .env ] && . ./.env; set +a; \
	SVC_UPPER=$$(echo "$(SERVICE)" | tr '[:lower:]' '[:upper:]'); \
	eval "DB_URL=\$${$${SVC_UPPER}_DATABASE_URL:-$(DEFAULT_DATABASE_URL)}"; \
	MIG_DIR="internal/$(SERVICE)/store/migrations"; \
	export GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$$DB_URL"; \
	if [ -n "$(TO)" ]; then \
		$(GOOSE) -dir "$$MIG_DIR" up-to "$(TO)"; \
	elif [ -n "$(STEPS)" ]; then \
		i=0; while [ "$$i" -lt "$(STEPS)" ]; do \
			$(GOOSE) -dir "$$MIG_DIR" up-by-one || exit $$?; \
			i=$$((i + 1)); \
		done; \
	else \
		$(GOOSE) -dir "$$MIG_DIR" up; \
	fi

_goose-down:
	@test -n "$(SERVICE)" || (echo "migrate-down: укажите SERVICE=..., например SERVICE=user"; exit 1)
	@test -d "internal/$(SERVICE)/store/migrations" || (echo "migrate-down: нет internal/$(SERVICE)/store/migrations"; exit 1)
	@set -a; [ -f .env ] && . ./.env; set +a; \
	SVC_UPPER=$$(echo "$(SERVICE)" | tr '[:lower:]' '[:upper:]'); \
	eval "DB_URL=\$${$${SVC_UPPER}_DATABASE_URL:-$(DEFAULT_DATABASE_URL)}"; \
	MIG_DIR="internal/$(SERVICE)/store/migrations"; \
	export GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$$DB_URL"; \
	if [ -n "$(TO)" ]; then \
		$(GOOSE) -dir "$$MIG_DIR" down-to "$(TO)"; \
	elif [ -n "$(STEPS)" ]; then \
		i=0; while [ "$$i" -lt "$(STEPS)" ]; do \
			$(GOOSE) -dir "$$MIG_DIR" down || exit $$?; \
			i=$$((i + 1)); \
		done; \
	else \
		$(GOOSE) -dir "$$MIG_DIR" down; \
	fi

_goose-create:
	@test -n "$(SERVICE)" || (echo "create-mig: укажите SERVICE=..., например SERVICE=user"; exit 1)
	@test -n "$(NAME)" || (echo "create-mig: укажите NAME=..., например NAME=add_avatar_column"; exit 1)
	@mkdir -p "internal/$(SERVICE)/store/migrations"
	@$(GOOSE) -dir "internal/$(SERVICE)/store/migrations" create "$(NAME)" sql

_goose-status:
	@test -n "$(SERVICE)" || (echo "migrate-status: укажите SERVICE=..., например SERVICE=user"; exit 1)
	@test -d "internal/$(SERVICE)/store/migrations" || (echo "migrate-status: нет internal/$(SERVICE)/store/migrations"; exit 1)
	@set -a; [ -f .env ] && . ./.env; set +a; \
	SVC_UPPER=$$(echo "$(SERVICE)" | tr '[:lower:]' '[:upper:]'); \
	eval "DB_URL=\$${$${SVC_UPPER}_DATABASE_URL:-$(DEFAULT_DATABASE_URL)}"; \
	MIG_DIR="internal/$(SERVICE)/store/migrations"; \
	export GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$$DB_URL"; \
	$(GOOSE) -dir "$$MIG_DIR" status
