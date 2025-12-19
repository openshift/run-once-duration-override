all: build
.PHONY: all

GO=GO111MODULE=on GOFLAGS=-mod=vendor go

OUTPUT_DIR := "./_output"
ARTIFACTS := "./artifacts/manifests"
MANIFEST_DIR := "$(OUTPUT_DIR)/manifests"
CERT_FILE_PATH := "$(OUTPUT_DIR)/certs.yaml"
MANIFEST_SECRET_YAML := "$(MANIFEST_DIR)/400_secret.yaml"
MANIFEST_MUTATING_WEBHOOK_YAML := "$(MANIFEST_DIR)/600_mutating.yaml"

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/images.mk \
)

# Exclude e2e tests from unit testing 
GO_TEST_PACKAGES :=./pkg/... ./cmd/...
GO_BUILD_FLAGS :=-tags strictfipsruntime

IMAGE_REGISTRY :=registry.svc.ci.openshift.org

# This will call a macro called "build-image" which will generate image specific targets based on the parameters:
# $0 - macro name
# $1 - target name
# $2 - image ref
# $3 - Dockerfile path
# $4 - context directory for image build
$(call build-image,run-once-duration-override,$(CI_IMAGE_REGISTRY)/ocp/4.12:run-once-duration-override,./images/ci/Dockerfile,.)

test-e2e: GO_TEST_PACKAGES :=./test/e2e
test-e2e: GO_TEST_FLAGS :=-v
test-e2e: test-unit
.PHONY: test-e2e

# generate manifests for installing on a dev cluster.
manifests:
	rm -rf $(MANIFEST_DIR)
	mkdir -p $(MANIFEST_DIR)
	cp -r $(ARTIFACTS)/* $(MANIFEST_DIR)/

	# generate certs
	./hack/generate-cert.sh "$(CERT_FILE_PATH)"

	# load the certs into the manifest yaml.
	./hack/load-cert-into-manifest.sh "$(CERT_FILE_PATH)" "$(MANIFEST_SECRET_YAML)" "$(MANIFEST_MUTATING_WEBHOOK_YAML)"

clean:
	$(RM) -r ./apiserver.local.config
	$(RM) -r ./_output
.PHONY: clean
