BIN := $(CURDIR)/bin
KUSTOMIZE := $(BIN)/kustomize
KUBECTL := $(BIN)/kubectl
NAMESPACE := wioctl
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

prometheus-operator:									## Deploy Prometheus Operator
	@$(info Deploying Prometheus Operator)
	@$(KUBECTL) create namespace monitoring --dry-run -o yaml | $(KUBECTL) apply -f -
	@$(KUBECTL) apply --wait -n default -f https://raw.githubusercontent.com/coreos/prometheus-operator/v0.34.0/bundle.yaml --all
	@sleep 5
	@$(KUBECTL) wait -n default --for condition=established crds --all --timeout=60s

prometheus:												## Deploy Prometheus
	$(info Deploying Prometheus)
	$(KUBECTL) kustomize deploy/prometheus | $(KUBECTL) apply -n $(NAMESPACE) -f -

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

prometheus:
	$(KUBECTL) port-forward svc/prometheus-operated 9090:9090