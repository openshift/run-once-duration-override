FROM registry.ci.openshift.org/openshift/release:golang-1.19 AS builder
WORKDIR /go/src/github.com/openshift/run-once-duration-override
COPY . .

RUN make build
FROM registry.ci.openshift.org/ocp/4.12:base
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/run-once-duration-override /usr/bin/
COPY --from=builder /go/src/github.com/openshift/run-once-duration-override/artifacts/configuration.yaml /etc/runoncedurationoverride/config/override.yaml

ENV CONFIGURATION_PATH=/etc/runoncedurationoverride/config/override.yaml
