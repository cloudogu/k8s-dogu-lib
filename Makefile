# Set these to the desired values
PROJECT_NAME=k8s-dogu-lib
ARTIFACT_ID=k8s-dogu-operator-crd
APPEND_CRD_SUFFIX=false
VERSION=3.0.0

IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG=1.26.0
MOCKERY_VERSION=v2.53.5
LINT_VERSION=v2.10.1
MAKEFILES_VERSION=10.11.0
CONTROLLER_GEN_VERSION=v0.19.0

ADDITIONAL_CLEAN=dist-clean

CRD_DOGU_SOURCE = ${HELM_CRD_SOURCE_DIR}/templates/k8s.cloudogu.com_dogus.yaml
CRD_POST_MANIFEST_TARGETS = crd-add-labels crd-add-conversion-webhook crd-copy-for-go-embedding

# Cross-repo contract with k8s-dogu-operator, which must provide a matching webhook Service and
# cert-manager Certificate in the same namespace as this CRD chart's release.
DOGU_CONVERSION_WEBHOOK_SERVICE_NAME ?= k8s-dogu-operator-webhook
DOGU_CONVERSION_WEBHOOK_SERVICE_PORT ?= 443
DOGU_CONVERSION_WEBHOOK_CERT_REF ?= k8s-dogu-operator-webhook-cert

PRE_COMPILE = generate-deepcopy
IMAGE_IMPORT_TARGET=image-import
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
include build/make/k8s-controller.mk
include build/make/release.mk

.PHONY: crd-copy-for-go-embedding
crd-copy-for-go-embedding:
	@echo "Copy CRD to crds/ for embedding"
	@cp ${CRD_DOGU_SOURCE} crds/

# controller-gen does not populate spec.conversion for multi-version CRDs on its own (that's
# normally done via a kustomize patch, which this repo's raw generation pipeline doesn't use).
# {{ .Release.Namespace }} is a Helm placeholder: templates/*.yaml is itself a Helm chart
# template, and yq only ever treats it as an opaque string, so it survives untouched here and is
# rendered later by Helm at chart-install time.
.PHONY: crd-add-conversion-webhook
crd-add-conversion-webhook: $(BINARY_YQ)
	@echo "Adding conversion webhook block to Dogu CRD..."
	@$(BINARY_YQ) -i e '.metadata.annotations."cert-manager.io/inject-ca-from" = "{{ .Release.Namespace }}/${DOGU_CONVERSION_WEBHOOK_CERT_REF}"' ${CRD_DOGU_SOURCE}
	@$(BINARY_YQ) -i e '.spec.conversion.strategy = "Webhook"' ${CRD_DOGU_SOURCE}
	@$(BINARY_YQ) -i e '.spec.conversion.webhook.conversionReviewVersions = ["v1"]' ${CRD_DOGU_SOURCE}
	@$(BINARY_YQ) -i e '.spec.conversion.webhook.clientConfig.service.name = "${DOGU_CONVERSION_WEBHOOK_SERVICE_NAME}"' ${CRD_DOGU_SOURCE}
	@$(BINARY_YQ) -i e '.spec.conversion.webhook.clientConfig.service.namespace = "{{ .Release.Namespace }}"' ${CRD_DOGU_SOURCE}
	@$(BINARY_YQ) -i e '.spec.conversion.webhook.clientConfig.service.port = ${DOGU_CONVERSION_WEBHOOK_SERVICE_PORT}' ${CRD_DOGU_SOURCE}
	@$(BINARY_YQ) -i e '.spec.conversion.webhook.clientConfig.service.path = "/convert"' ${CRD_DOGU_SOURCE}

# Override make target to use k8s-dogu-lib as label
.PHONY: crd-add-labels
crd-add-labels: $(BINARY_YQ)
	@echo "Adding labels to CRD..."
	@for file in ${HELM_CRD_SOURCE_DIR}/templates/*.yaml ; do \
		$(BINARY_YQ) -i e ".metadata.labels.app = \"ces\"" $${file} ;\
		$(BINARY_YQ) -i e ".metadata.labels.\"app.kubernetes.io/name\" = \"${PROJECT_NAME}\"" $${file} ;\
	done
