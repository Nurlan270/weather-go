include .env
export

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

GOOSE = $(PROJECT_BIN)/goose
LINTER = $(PROJECT_BIN)/golangci-lint

$(shell [ -d bin ] || mkdir -p $(PROJECT_BIN))

.PHONY: install-goose
install-goose:
	@[ -f $(GOOSE) ] || \
 	curl -sSfL \
 		https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
 		GOOSE_INSTALL=$(PROJECT_DIR) sh -s

goose-status: install-goose
	$(GOOSE) status

goose-up: install-goose
	$(GOOSE) up

goose-down: install-goose
	$(GOOSE) down

goose-reset: install-goose
	$(GOOSE) reset

goose-create: install-goose
	$(GOOSE) create $(ARGS)

.PHONY: install-linter
install-linter:
	@[ -f $(LINTER) ] || \
	curl -sSfL \
		https://golangci-lint.run/install.sh | sh -s v2.12.2

lint: install-linter
	$(LINTER) run

lint-file: install-linter
	$(LINTER) run $(FILE)

lint-fmt: install-linter
	$(LINTER) fmt