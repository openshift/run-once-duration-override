# Architecture

## Overview

This repository is the **operand**: a mutating admission webhook that sets `activeDeadlineSeconds` on run-once pods (`RestartPolicy: Never` or `OnFailure`). It is not an operator.

Lifecycle is owned by [run-once-duration-override-operator](https://github.com/openshift/run-once-duration-override-operator), which is installed by OLM. The operator reconciles the `RunOnceDurationOverride` CR; this binary performs the admission mutation.

If a pod already has `activeDeadlineSeconds`, the webhook uses `min(existing, configured)` and never increases the deadline.

## Data flow

```text
  Pod CREATE/UPDATE
  (namespaced, labeled ns only)
              │
              ▼
  kube-apiserver  ──HTTPS──►  localhost:9448
                              (hostNetwork DaemonSet on control-plane nodes)
              │
              ▼
  generic-admission-server
  cmd/run-once-duration-override/
              │
              ▼
  Admit hook
    1. IsApplicable?  (pods, CREATE/UPDATE, no subresource)
    2. IsExempt?      (RestartPolicy not Never/OnFailure → allow)
    3. Mutate         (set/min activeDeadlineSeconds)
    4. RFC 6902 JSON patch
              │
              ▼
  AdmissionResponse (patch or allow)
```

The `MutatingWebhookConfiguration` uses `clientConfig.url: https://localhost:9448/...` because the webhook runs with `hostNetwork` on control-plane nodes and binds to `127.0.0.1:9448`. It is not a Service/ClusterIP webhook.

## Operand startup

Entry: `cmd/run-once-duration-override/main.go` → `generic-admission-server` `RunAdmissionServer`.

1. Register mutating resource `admission.runoncedurationoverride.openshift.io/v1/runoncedurationoverrides`
2. Serve TLS on `--secure-port` (9448) with `--bind-address=127.0.0.1`
3. Post-start hook: `Initialize` loads config from `CONFIGURATION_PATH`
4. Handle `Admit` until shutdown

Config is loaded **once** at initialization and cached. It is not re-read per request. In cluster, the operator writes that file from a ConfigMap and rolls the DaemonSet when the CR hash changes.

## Admission pipeline

Implemented in `pkg/runoncedurationoverride/` and invoked from `runOnceDurationOverrideHook.Admit`:

| Step | Behavior |
|------|----------|
| Not initialized | deny (`AdmissionResponse` code 500) |
| Not a Pod CREATE/UPDATE (or a pod subresource) | allow unchanged |
| Unmarshal failure | deny (`AdmissionResponse` code 400) |
| `RestartPolicy` is not `Never` or `OnFailure` | allow unchanged (exempt) |
| Mutate | `activeDeadlineSeconds = min(configured, existing)` (existing may be unset) |
| Patch | RFC 6902 JSON patch vs original object |

Default configured value is `3600` (`artifacts/configuration.yaml`).

## Scope

The webhook does **not** run for every namespace. `MutatingWebhookConfiguration` requires:

```text
runoncedurationoverrides.admission.runoncedurationoverride.openshift.io/enabled: "true"
```

Namespaces with `runlevel` in `{0,1}` are excluded. `failurePolicy: Fail`, `timeoutSeconds: 5`, `reinvocationPolicy: IfNeeded`.

Namespace: `run-once-duration-override`. Workload: DaemonSet on `node-role.kubernetes.io/master` with hostNetwork.

## Configuration

| Source | Purpose |
|--------|---------|
| `CONFIGURATION_PATH` | Path to `RunOnceDurationOverrideConfig` YAML (`spec.activeDeadlineSeconds`) |
| `--tls-cert-file` / `--tls-private-key-file` | Serving certs (operator-managed in cluster; `hack/` + `make manifests` for local deploy) |

## Design decisions

| Decision | Rationale |
|----------|-----------|
| Mutating webhook, not a controller | Must change the Pod spec at create/update before it is persisted |
| `min(existing, configured)` | Never lengthen a tighter deadline the user already set |
| Only `Never` / `OnFailure` | Those are run-once pods; long-running pods must not get a deadline injected |
| Namespace label opt-in | Avoid mutating platform namespaces by default |
| hostNetwork DaemonSet on control-plane + `https://localhost` | Webhook is local to apiserver nodes; no Service hop |
| Config cached at startup | Simple, deterministic; live file updates are ignored until process restart |
| RFC 6902 patches | Admission API expects JSON Patch, not strategic merge |

## Testing

- **Unit:** `make test-unit` — `pkg/...` and `cmd/...` (exemption, mutator, patch)
- **E2E:** `make test-e2e` — needs `KUBECONFIG`, `IMAGE_FORMAT`, `NAMESPACE`; covers default `3600`, `min(existing, configured)`, and `RestartPolicy: Always` exemption
