BIN := $(CURDIR)/bin
KUSTOMIZE := $(BIN)/kustomize
KUBECTL := $(BIN)/kubectl
SUBDIRS = bin
.DEFAULT_GOAL = help

##########################################################
##@ GO
##########################################################

doc:
	$(info http://localhost:6060/pkg/$(MODULE))
	godoc

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser

fmt:
	go fmt ./...

##########################################################
##@ BOOTSTRAP
##########################################################
.PHONY: bootstrap

bootstrap:
	make -C $(BIN)

##########################################################
##@ DEPLOY
##########################################################
.PHONY: hide reveal deploy

hide:
	git secret hide -d

reveal:
	@git secret reveal -f

deploy: bootstrap reveal
	$(KUSTOMIZE) build deploy | $(KUBECTL) apply -f -
