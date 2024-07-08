include .env
PROJECT := "go-effective-mobile"
USER := rptl8sr
EMAIL := $(USER)@gmail.com
LOCAL_BIN:=$(CURDIR)/bin
MIGRATIONS_DIR=$(CURDIR)/migrations
GO_VERSION?=1.22.5
GO := go
DSN="host=$(PG_HOST) port=$(PG_PORT) dbname=$(PG_DATABASE) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=$(PG_SSL_MODE)"

.PHONY: git-init
git-init:
	gh repo create $(PROJECT) --private
	git init
	git config user.name "$(USER)"
	git config user.email "$(EMAIL)"
	git add go.mod Makefile
	git commit -m "Init commit"
	git remote add origin git@github.com:$(USER)/$(PROJECT).git
	git remote -v
	git push -u origin master
	mkdir -p $(LOCAL_BIN)
	mkdir -p $(MIGRATIONS_DIR)


BN ?= dev
# make git-checkout BN=dev
.PHONY: git-checkout
git-checkout:
	git checkout -b $(BN)


.PHONY: golangci-lint-install
golangci-lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1


.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run ./...


.PHONY: goose-get
get-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.21.1


.PHONY: goose-make-migrations
make-goose:
ifndef MN
	$(error MN is undefined)
endif
	$(LOCAL_BIN)/goose -dir=$(MIGRATIONS_DIR) create $(MN) sql


.PHONY: goose-migrate-status
goose-migrate-status:
	$(LOCAL_BIN)/goose -dir $(MIGRATIONS_DIR) postgres $(DSN) status -v

.PHONY: goose-migrate-up
goose-migrate-up:
	$(LOCAL_BIN)/goose -dir $(MIGRATIONS_DIR) postgres $(DSN) up -v


.PHONY: goose-migrate-down
goose-migrate-down:
	$(LOCAL_BIN)/goose -dir $(MIGRATIONS_DIR) postgres $(DSN) down -v