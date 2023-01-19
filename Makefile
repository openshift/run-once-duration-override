all: build
.PHONY: all

GO=GO111MODULE=on GOFLAGS=-mod=vendor go

OUTPUT_DIR := "./_output"

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/images.mk \
)

IMAGE_REGISTRY :=registry.svc.ci.openshift.org

# This will call a macro called "build-image" which will generate image specific targets based on the parameters:
# $0 - macro name
# $1 - target name
# $2 - image ref
# $3 - Dockerfile path
# $4 - context directory for image build
$(call build-image,run-once-duration-override,$(CI_IMAGE_REGISTRY)/ocp/4.12:run-once-duration-override,./images/ci/Dockerfile,.)
