# Review guidance — cass-operator

<!-- Draft harvested from past maintainer review comments. Owner to correct. Shared rules live in shared/go-operator-guidelines.md; this file covers only what is specific to this repo. -->

Owners: @k8ssandra/maintainers

## What failure looks like here

A healthy Cassandra cluster gets disrupted. This includes an unexpected rolling restart after an operator upgrade, a scale-down or decommission that gets stuck or cleans up too early, data or nodes lost, or a reconcile that loops forever or floods the API server.

## Focus on

1. **Upgrade safety.** After an operator upgrade, existing clusters must behave as before unless the PR says otherwise.
   - Any change to the generated StatefulSet, pod template, Services or labels can restart or change every existing cluster. Flag it unless the PR says it is intended and the CHANGELOG mentions it.
   - Also look for changed behaviour on existing code paths: Management API client and transport settings, timeouts, defaults.
   - When the impact depends on the environment and you can't show it, ask whether the change is intended instead of describing a failure.
2. **Server and version gating.** Behaviour differs across `cassandra`, `dse` and `hcd`, and across major versions. For example, MCAC exists only on older versions. Check that version-specific code is gated correctly and that the gate is right.
3. **Rack and StatefulSet sequencing.** Check rack balance, ordering, and acting on StatefulSets before their controller has caught up. Ask where in the sequence the edge case is handled.
4. **Long operations as state machines.** Scale, decommission, replace and tasks must be resumable after a restart and must not finish early when a Management API call fails. Check what happens when a request fails halfway.
5. **Transient vs terminal errors and requeues** ([shared guidelines §1](../../shared/go-operator-guidelines.md#1-reconcile-loop)). Look especially at code that polls Cassandra or the Management API.
6. **Per-reconcile cost.** New reads of Secrets, PVCs, Endpoints or config on every reconcile are expensive in cluster-wide deployments.
7. **Names and labels.** Kubernetes object names follow `metadata.name`, never the datacenter name override. Label value rules are not the same as annotation value rules.
8. **Webhook coverage.** Invalid transitions, such as scaling down while scaling up, are blocked in the webhook, not handled late in reconcile.
9. **Compatibility with k8ssandra-operator and k8ssandra-client.**
   - Watch for changes to defaults or fields that k8ssandra-operator sets or reads.
   - Both projects import `pkg/httphelper` as a library, so a change to its exported behaviour is an API change.
   - k8ssandra-client passes its own transport for port-forwarded connections.
10. **Tests prove the behaviour, at the right level.** Fake Management API features match what the test claims. Controller and webhook changes come with envtest coverage. Operational changes come with an end-to-end suite, or an explanation of why not.

## Repo hygiene

Flag these as Important unless the item says otherwise. Items marked *(CI)* are mechanical. They stay here only until the matching CI check exists, then the line is deleted ([plan §3.1](../../00-plan.md#31-deterministic-checks-are-code-not-review)).

- **No changelog entry** *(CI, once the changelog check exists; there is none yet)*. Until changelog fragments (1.6b) land, the entry goes under `## unreleased` in `CHANGELOG.md`. Every change outside docs needs one, including a PR titled as a test fix that changes production code. Don't flag the category: which category fits is a judgement call.
- Image or version references hardcoded in Go instead of `config/imageconfig`.
- Unrelated changes: reformatting, untouched files in the diff, test data edited without the test needing it.
- **User-facing change without docs**: a change users can see (fields, annotations, behaviour, metrics) with no docs note or docs link, even when the changelog kind isn't `FEATURE` ([plan §5.6](../../00-plan.md#56-user-facing-documentation-docs-sync)).
- **Test sprawl** ([shared §6, items 7–8](../../shared/go-operator-guidelines.md#6-tests)): standalone regression tests, new tests for a change that doesn't alter behaviour, or new tests that only repeat assertions an existing test already makes. Report it once per PR and list every such test. For each one, say which existing test the case belongs in, or that it can be deleted.
  - **Important** when an existing test already owns the behaviour and the new test bypasses it.
  - **Nit** when no test owns the behaviour yet. Then the new test becomes the owning test, and renaming it after the behaviour or function is the whole fix.

## Needs a human (a marker, not a finding)

Some changes need human approval whatever the risk tier, even when they are correct. Don't report them as findings. List them on one line at the end of the summary: `Needs human: <what> in <file:line, ...>. Explained in the description: yes/no.`

- **New `//nolint` directives** *(CI detector, T17)*. The directive must name the linter and give a reason.
- **Changed or removed assertions in existing tests** *(CI detector, T17)*: different expected values, relaxed matchers, deleted checks, skipped tests.

Report an Important finding only when the description doesn't explain the change, or the explanation doesn't match the diff.

## Intentional trade-offs

If the PR description states a trade-off or a known limitation, don't flag it. Raise it only if it leads to one of the failures under "What failure looks like here". Then ask the author about it as a question rather than reporting it as a finding.

## Don't flag

- Style the linter enforces (`.golangci.yml`).
- Comments the shared guidelines say to delete. Flag that as a single Nit per PR, not one finding per comment.
- The wording of error and log messages that use the terms users know. For example, "the Cassandra version does not support it" is fine when the exact cause is the Management API image.

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
- A transient error that surfaces as a Warning event or an error log but converges is a Nit, unless it loops.
- **Pre-existing** findings don't count toward the limits.
  - List them in the summary only, not as inline comments.
  - Give each one a one-line proposed `tech-debt` issue. Check first whether an issue already exists.
- Don't repeat findings already posted on this PR.
- If the same problem appears in several places, report it once and list every location.

## Summary format

Start with a count, for example "N important, M nits, K pre-existing." If there is nothing to fix, say so plainly. End with the `Needs human:` line when it applies.
