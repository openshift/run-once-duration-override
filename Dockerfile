FROM brew.registry.redhat.io/rh-osbs/openshift-golang-builder:rhel_9_1.22 as builder
WORKDIR /go/src/github.com/openshift/run-once-duration-override
COPY . .
RUN make build --warn-undefined-variables

FROM registry.redhat.io/rhel9-4-els/rhel-minimal:9.4
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/run-once-duration-override /usr/bin/
RUN mkdir /licenses
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override-operator/LICENSE /licenses/.

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
