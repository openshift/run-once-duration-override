# AI Agent Guide for Run Once Duration Override

This repo is the **operand**: a mutating admission webhook that sets `activeDeadlineSeconds` on run-once pods. Lifecycle is owned by [run-once-duration-override-operator](https://github.com/openshift/run-once-duration-override-operator) (OLM). This is not an operator or controller.

If a pod already has `activeDeadlineSeconds`, the webhook uses `min(existing, configured)` and **never increases** the deadline. Only pods with `RestartPolicy: Never` or `OnFailure` are mutated.

Design and admission flow: [ARCHITECTURE.md](ARCHITECTURE.md). PR process: [CONTRIBUTING.md](CONTRIBUTING.md).

## Build and Test

```bash
make build        # Build the binary
make test-unit    # Unit tests
make test-e2e     # Requires KUBECONFIG, IMAGE_FORMAT, NAMESPACE
make manifests    # TLS certs + manifests into _output/manifests/
make verify       # Formatting and vet
```

See `go.mod` for the Go version.

## Layout

| Path | Purpose |
|------|---------|
| `cmd/run-once-duration-override/` | Admission server entry (`generic-admission-server` hook) |
| `pkg/runoncedurationoverride/` | Config, exemption, mutation, JSON patch |
| `pkg/api/` | Admission API group/version constants |
| `pkg/response/` | Admission response helpers |
| `artifacts/` | Default config YAML and deploy manifests |
| `test/e2e/` | Cluster e2e (default mutation, min value, RestartPolicy exemption) |
| `hack/` | Local TLS cert generation for manifests |
| `vendor/` | Vendored deps — do not edit; use `go mod tidy && go mod vendor` |

## Admission flow

```
Pod CREATE/UPDATE → IsApplicable → RestartPolicy exempt check → mutate from cached config → RFC 6902 JSON patch
```

- Webhook applies only in namespaces labeled `runoncedurationoverrides.admission.runoncedurationoverride.openshift.io/enabled: "true"`.
- Config comes from `CONFIGURATION_PATH` and is loaded **once at startup**, not per request.
- Local `hack/` cert scripts are for `make manifests` only. In cluster, the operator owns certs.

## Do not

- Change exemption or min-deadline semantics without an explicit request.
- Log TLS material, CA bundles, or tokens.
- Edit `vendor/` or `OWNERS` without maintainer direction.
- Skip `make verify` and `make test-unit`.
