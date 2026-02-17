FROM brew.registry.redhat.io/rh-osbs/openshift-golang-builder:rhel_9_1.24 as builder
WORKDIR /go/src/github.com/openshift/run-once-duration-override
COPY . .
RUN make build --warn-undefined-variables

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:c7d44146f826037f6873d99da479299b889473492d3c1ab8af86f08af04ec8a0
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/run-once-duration-override /usr/bin/
RUN mkdir /licenses
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/LICENSE /licenses/.

LABEL io.k8s.display-name="Run Once Duration Override mutating admission webhook based on RHEL 9" \
      io.k8s.description="This is a component of OpenShift for the Run Once Duration Override mutating admission webhook based on RHEL 9" \
      distribution-scope="public" \
      com.redhat.component="run-once-duration-override-container" \
      name="run-once-duration-override-operator/run-once-duration-override-rhel9" \
      cpe="cpe:/a:redhat:run_once_duration_override_operator:1.4::el9" \
      release="1.4.0" \
      version="1.4.0" \
      url="https://github.com/openshift/run-once-duration-override" \
      vendor="Red Hat, Inc." \
      summary="run-once-duration-override" \
      io.openshift.expose-services="" \
      io.openshift.tags="openshift,run-once-duration-override" \
      description="run-once-duration-override-container" \
      maintainer="AOS workloads team, <aos-workloads@redhat.com>"

USER nobody
