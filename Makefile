PROJECT := "go-effective-mobile"
USER := rptl8sr
EMAIL := $(USER)@gmail.com
LOCAL_BIN:=$(CURDIR)/bin
GO_VERSION?=1.22.5
GO := go


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
