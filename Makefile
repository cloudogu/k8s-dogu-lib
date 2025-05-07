# Set these to the desired values
ARTIFACT_ID=k8s-dogu-operator-crd
APPEND_CRD_SUFFIX=false
VERSION=2.6.0

IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG=1.24.1
LINT_VERSION=v1.64.7
MAKEFILES_VERSION=9.9.1

ADDITIONAL_CLEAN=dist-clean

PRE_COMPILE = generate-deepcopy
CHECK_VAR_TARGETS=check-all-vars

include build/make/variables.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/mocks.mk
include build/make/k8s.mk
include build/make/k8s-crd.mk

##@ Debug

.PHONY: print-debug-info
print-debug-info: ## Generates info and the list of environment variables required to start the operator in debug mode.
	@echo "The target generates a list of env variables required to start the operator in debug mode. These can be pasted directly into the 'go build' run configuration in IntelliJ to run and debug the operator on-demand."
	@echo "STAGE=$(STAGE);LOG_LEVEL=$(LOG_LEVEL);KUBECONFIG=$(KUBECONFIG);NAMESPACE=$(NAMESPACE);DOGU_REGISTRY_ENDPOINT=$(DOGU_REGISTRY_ENDPOINT);DOGU_REGISTRY_USERNAME=$(DOGU_REGISTRY_USERNAME);DOGU_REGISTRY_PASSWORD=$(DOGU_REGISTRY_PASSWORD)"
