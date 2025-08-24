# Copyright 2023 The run-once-duration-operator Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

.PHONY: verify-boilerplates
verify-boilerplates:
	go tool boilersuite --author "run-once-duration-operator" .

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
