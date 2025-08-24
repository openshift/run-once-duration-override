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

# Used for Konflux pipeline

FROM brew.registry.redhat.io/rh-osbs/openshift-golang-builder:rhel_9_1.24 as builder
WORKDIR /go/src/github.com/openshift/run-once-duration-override
COPY . .
RUN make build --warn-undefined-variables

FROM registry.redhat.io/rhel9-4-els/rhel-minimal:9.4
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/run-once-duration-override /usr/bin/
RUN mkdir /licenses
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/LICENSE /licenses/.

LABEL io.k8s.display-name="Run Once Duration Override mutating admission webhook based on RHEL 9" \
      io.k8s.description="This is a component of OpenShift for the Run Once Duration Override mutating admission webhook based on RHEL 9" \
      com.redhat.component="run-once-duration-override-container" \
      name="run-once-duration-override-rhel-9" \
      summary="run-once-duration-override" \
      io.openshift.expose-services="" \
      io.openshift.tags="openshift,run-once-duration-override" \
      description="run-once-duration-override-container" \
      maintainer="AOS workloads team, <aos-workloads@redhat.com>"

USER nobody
