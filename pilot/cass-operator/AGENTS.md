# cass-operator

<!-- Draft for k8ssandra/cass-operator. Lines marked "confirm" need an owner's check before committing. -->

## What this is

A Kubernetes operator that deploys and operates Apache Cassandra, DSE and HCD clusters. Its main resources are `CassandraDatacenter` (`cassandra.datastax.com/v1beta1`) and `CassandraTask` (`control.k8ssandra.io/v1alpha1`). It talks to each node through the Management API sidecar. It runs standalone and as part of k8ssandra-operator, so changes must work for both.

## Layout

| Path | What lives there |
|---|---|
| `apis/` | CRD types and API groups (`cassandra`, `control`, `config`). `zz_generated.deepcopy.go` is generated |
| `internal/controllers/cassandra`, `internal/controllers/control` | Controllers for CassandraDatacenter and CassandraTask |
| `internal/webhooks/` | Validating webhook for CassandraDatacenter, which blocks invalid changes and transitions |
| `internal/result/` | Reconcile result helpers: `Continue()`, `Done()`, `RequeueSoon()`, `Error()` |
| `pkg/reconciliation/` | Core reconcile logic. The entry point is `CalculateReconciliationActions` (`handler.go`); the rack and StatefulSet lifecycle is in `ReconcileAllRacks` (`reconcile_racks.go`) |
| `pkg/httphelper/` | Management API client. `server_test_utils.go` holds the fake server used in tests |
| `pkg/images/`, `config/imageconfig/` | Image configuration. Image versions go in config, never in Go code |
| `tests/` | End-to-end Ginkgo/Gomega suites, one directory per scenario, run against a kind cluster. `tests/testdata/` holds fixtures |
| `hack/cluster.sh` | Creates the kind cluster for end-to-end tests (`ipv4`, `ipv6` or `dual`) |
| `config/` | Kustomize manifests. `config/crd/bases` is generated |
| `docs/developer/` | Developer docs: release process, workflow details (`kind.md` is outdated; use `hack/cluster.sh`) |

## How to verify changes

| What | Command | Notes |
|---|---|---|
| Format, vet, lint | `make fmt vet lint` | golangci-lint config in `.golangci.yml` |
| Generate code and CRDs | `make manifests generate` | Required after any change in `apis/`. Commit the generated files. Uses the Makefile's pinned controller-gen |
| Unit and envtest | `make test` | Runs generate, fmt, vet and lint first. envtest suites start a real API server |
| One end-to-end suite | `hack/cluster.sh` then `make integ-test M_INTEG_DIR=<dir>` | Needs a kind cluster created with `hack/cluster.sh` |

Run `make test` before opening a PR. CI results appear on the PR, so don't paste them into the description. Describe only manual testing there.

## Conventions

- Follow [`shared/go-operator-guidelines.md`](../../shared/go-operator-guidelines.md) <!-- confirm final location/link -->. It covers the reconcile loop, Kubernetes semantics, API design, Go style, comments and tests.
- **CHANGELOG:** add an entry under `## unreleased` in the order `[CHANGE]`, `[FEATURE]`, `[ENHANCEMENT]`, `[BUGFIX]`. Every entry links its GitHub issue. Docs-only PRs need no entry. <!-- replace with the changelog-fragment instructions once step 1.6b lands -->
- **Branch updates:** merge `master` into your branch only when there's a conflict. Don't rebase or force-push once human review has started.
- **Lint suppressions:** avoid `//nolint`. If one is unavoidable, name the linter and give a reason; the PR will need human approval.
- **Every PR links an issue.** A workflow checks this.
- **API changes** must stay compatible with existing `CassandraDatacenter` resources and with k8ssandra-operator.
- **Server differences matter.** `serverType` is `cassandra`, `dse` or `hcd`, and behaviour often depends on version, for example 3.11, 4.x or 5.x. Gate version-specific code explicitly.
- **Three test levels** (details in the [shared guidelines §6](../../shared/go-operator-guidelines.md#6-tests)):
  - **Unit tests** sit next to the code and use testify, with `require` for error checks.
  - **envtest** suites live in the controller and webhook packages (`internal/controllers/*`, `internal/webhooks/*`), use Ginkgo/Gomega, and run against a real API server. Prefer envtest for controller behaviour.
  - **End-to-end** suites live in `tests/`, use Ginkgo/Gomega, and run on kind.

## Goals

<!-- confirm -->
- Operate stateful clusters safely. An operator upgrade or a reconcile must never disrupt a healthy cluster unexpectedly.
- Scale to large, cluster-wide deployments. Keep per-reconcile API-server load and memory low.
- Converge. Every operation (scale, decommission, replace, upgrade, task) reaches a clear end state or a clear failure.

## Non-goals

<!-- confirm -->
- A second way to configure something that already has one. For example, ports are configured through `containerPort` in the pod template, not through `cassandra-yaml`.

## Further reading

- Review guidance: [`REVIEW.md`](REVIEW.md)
- Developer docs: `docs/developer/`
