# Review guidance — cass-operator

<!-- Draft harvested from past maintainer review comments. Owner to correct. Shared rules live in shared/go-operator-guidelines.md; this file covers only what is specific to this repo. -->

Owners: @k8ssandra/maintainers

## What failure looks like here

A healthy Cassandra cluster gets disrupted. This includes an unexpected rolling restart after an operator upgrade, a scale-down or decommission that gets stuck or cleans up too early, data or nodes lost, or a reconcile that loops forever or floods the API server.

## Focus on

1. **Upgrade safety.** Any change to the generated StatefulSet, pod template, Services or labels can restart or change every existing cluster when the operator is upgraded. Flag it unless the PR says it is intended and the CHANGELOG mentions it.
2. **Server and version gating.** Behaviour differs across `cassandra`, `dse` and `hcd`, and across major versions. For example, MCAC exists only on older versions. Check that version-specific code is gated correctly and that the gate is right.
3. **Rack and StatefulSet sequencing.** Check rack balance, ordering, and acting on StatefulSets before their controller has caught up. Ask where in the sequence the edge case is handled.
4. **Long operations as state machines.** Scale, decommission, replace and tasks must be resumable after a restart and must not finish early when a Management API call fails. Check what happens when a request fails halfway.
5. **Transient vs terminal errors and requeues** ([shared guidelines §1](../../shared/go-operator-guidelines.md#1-reconcile-loop)). Look especially at code that polls Cassandra or the Management API.
6. **Per-reconcile cost.** New reads of Secrets, PVCs, Endpoints or config on every reconcile are expensive in cluster-wide deployments.
7. **Names and labels.** Kubernetes object names follow `metadata.name`, never the datacenter name override. Label value rules are not the same as annotation value rules.
8. **Webhook coverage.** Invalid transitions, such as scaling down while scaling up, are blocked in the webhook, not handled late in reconcile.
9. **Compatibility with k8ssandra-operator.** Changes to defaults or fields that k8ssandra-operator sets or reads.
10. **Tests prove the behaviour, at the right level.** Fake Management API features match what the test claims. Controller and webhook changes come with envtest coverage. Operational changes come with an end-to-end suite, or an explanation of why not.

## Repo hygiene (flag as Important)

Items marked *(CI)* are mechanical. They stay here only until the matching CI check exists, then the line is deleted ([plan §3.1](../../00-plan.md#31-deterministic-checks-are-code-not-review)).

- No changelog entry (a fragment, once 1.6b is in place; until then under `## unreleased` in `CHANGELOG.md`), a wrong category, or no linked issue. *(CI for a missing entry or issue; the category is judgement)*
- `apis/` changed but generated files not regenerated, or regenerated with a different controller-gen version. *(CI)*
- Image or version references hardcoded in Go instead of `config/imageconfig`.
- Unrelated changes: reformatting, untouched files in the diff, test data edited without the test needing it.
- **New `//nolint` directives** *(CI detector, T17)*: always flag these. The PR needs human approval whatever its risk tier, and the directive must name the linter and give a reason.
- **User-facing change without docs**: a change users can see (fields, annotations, behaviour, metrics) with no docs note or docs link, even when the changelog kind isn't `FEATURE` ([plan §5.6](../../00-plan.md#56-user-facing-documentation-docs-sync)).
- **Test sprawl** ([shared §6, items 7–8](../../shared/go-operator-guidelines.md#6-tests)): standalone regression tests, or new tests for a change that doesn't alter behaviour. Report it once per PR and list every such test. For each one, say which existing test the case belongs in, or that it can be deleted.
- **Changed or removed assertions in existing tests** *(CI detector, T17)*: different expected values, relaxed matchers, deleted checks, skipped tests. Always flag these, and the PR needs human approval whatever its risk tier. The PR description must explain why the expectation changed.

## Don't flag

- Style the linter enforces (`.golangci.yml`).
- Comments the shared guidelines say to delete. Flag that as a single Nit per PR, not one finding per comment.

## Skip these paths

- `**/zz_generated.deepcopy.go`
- `config/crd/bases/**`
- `go.sum`
- `tests/testdata/**`, unless the PR's purpose is test data

## Evidence

Every finding cites a file and line and describes the concrete failure: what breaks, for which cluster or version, and when. Don't guess from a name. Read the code it calls.

## Limits

- At most 3 Nits per review.
- Tag each finding Important, Nit or Pre-existing.
- Don't repeat findings already posted on this PR.
- If the same problem appears in several places, report it once and list every location.

## Summary format

Start with a count, for example "N important, M nits, K pre-existing." If there is nothing to fix, say so plainly.
