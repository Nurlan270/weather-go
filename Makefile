include .env
export

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

$(shell [ -d bin ] || mkdir -p $(PROJECT_BIN))

.PHONY: install-goose
install-goose:
	@[ -f $(PROJECT_BIN)/goose ] || \
 	curl -fsSL \
 		https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
 		GOOSE_INSTALL=$(PROJECT_DIR) sh -s

goose-status: install-goose
	$(PROJECT_BIN)/goose status

goose-up: install-goose
	$(PROJECT_BIN)/goose up

goose-down: install-goose
	$(PROJECT_BIN)/goose down

goose-reset: install-goose
	$(PROJECT_BIN)/goose reset

goose-create: install-goose
	$(PROJECT_BIN)/goose create $(ARGS)
