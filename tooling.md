# Tooling

The [plan](00-plan.md) says *what* has to be checked, measured or automated. This page lists *which* tools and agents do it. Tool choices that are still open are tracked in [`03-process-tracker.md`](03-process-tracker.md); scheduled maintenance agents are described in [`05-maintenance-agents.md`](05-maintenance-agents.md).

| Need | Tool | Status | Plan |
|---|---|---|---|
| Baseline and ongoing metrics, CI re-run rate and CI time | [`tools/review-metrics`](tools/review-metrics/) (Go CLI on the libraries `gh` uses), `--ci` for the CI numbers. Usage in [`04-templates.md` §8](04-templates.md#8-metrics-tool) | Exists | [§2](00-plan.md#2-goals-and-how-well-measure-them) |
| Copying shared files into each repo across orgs | `sync-review-files`, a command in our Go toolkit that copies the current version into each listed repo and opens a PR | Planned (D8) | [§5.4](00-plan.md#54-shared-guidance-across-a-family-of-repos) |
| Re-scan after a dependency update | govulncheck, OSV scanner | Already in CI on the pilot | [§4.4](00-plan.md#44-dependency-and-cve-updates) |
| Dependency and CVE update PRs | Renovate, Dependabot or our own scheduled agent | Open (T23) | [§4.4](00-plan.md#44-dependency-and-cve-updates) |
| Combined coverage map (unit, envtest, end-to-end) | `go build -cover` and `GOCOVERDIR` for end-to-end binaries, merged with `go tool covdata` | Open (T25) | [§3.2](00-plan.md#32-repository-readiness) |
| CRD field descriptions | controller-gen turns doc comments into the descriptions `kubectl explain` shows | Exists in operator repos | [§5.6](00-plan.md#56-user-facing-documentation-docs-sync) |
| Review-state workflow: is a review thread resolved? | GitHub GraphQL `reviewThreads.isResolved` (REST doesn't expose it) | Open (T22) | [§4.2](00-plan.md#42-knowing-when-bots-are-done-review-state-labels) |
| Automated code and security review | To be chosen | Open (T3, T4) | [§3](00-plan.md#3-the-layered-model) |
| Local review before the PR | To be chosen | Open (T8) | [§3](00-plan.md#3-the-layered-model) |
| Maintenance agents (flaky-test fixer, agent-context tuner, CI time tracker, …) | See [`05-maintenance-agents.md`](05-maintenance-agents.md) | Planned | [§3.2](00-plan.md#32-repository-readiness) |
