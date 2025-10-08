# Set these to the desired values
PROJECT_NAME=k8s-dogu-lib
ARTIFACT_ID=k8s-dogu-operator-crd
APPEND_CRD_SUFFIX=false
VERSION=2.10.0

IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG=1.25.1
MOCKERY_VERSION=v2.53.5
LINT_VERSION=v2.4.0
MAKEFILES_VERSION=10.2.1
CONTROLLER_GEN_VERSION=v0.19.0

ADDITIONAL_CLEAN=dist-clean

CRD_DOGU_SOURCE = ${HELM_CRD_SOURCE_DIR}/templates/k8s.cloudogu.com_dogus.yaml
CRD_POST_MANIFEST_TARGETS = crd-add-labels crd-copy-for-go-embedding

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
	@echo "Copy CRD to api/v2/"
	@cp ${CRD_DOGU_SOURCE} api/v2/

# Override make target to use k8s-dogu-lib as label
.PHONY: crd-add-labels
crd-add-labels: $(BINARY_YQ)
	@echo "Adding labels to CRD..."
	@for file in ${HELM_CRD_SOURCE_DIR}/templates/*.yaml ; do \
		$(BINARY_YQ) -i e ".metadata.labels.app = \"ces\"" $${file} ;\
		$(BINARY_YQ) -i e ".metadata.labels.\"app.kubernetes.io/name\" = \"${PROJECT_NAME}\"" $${file} ;\
	done
