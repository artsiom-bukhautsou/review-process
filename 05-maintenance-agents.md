# Maintenance agents

**Scope:** generic and team-neutral. Not published on the site yet (excluded in [`site/hugo.toml`](site/hugo.toml)). \
**Related:** [`00-plan.md` §3.1](00-plan.md#31-deterministic-checks-are-code-not-review) (deterministic checks are code), [§3.2](00-plan.md#32-repository-readiness) (repository readiness), [§5.6](00-plan.md#56-user-facing-documentation-docs-sync) (user docs), [§4.3](00-plan.md#43-ci-friction-flaky-tests-changelog-conflicts-and-branch-updates) (flaky tests), [§4.4](00-plan.md#44-dependency-and-cve-updates) (dependency updates), [§5.3](00-plan.md#53-how-these-docs-improve-over-time) (docs improve over time), [§6](00-plan.md#6-feedback-loops) (feedback loops); investigation [T10](03-process-tracker.md#t--tooling) and decision [T21](03-process-tracker.md#decision-log) in the process tracker. Each team chooses the flaky-issue destination in [T21-A](02-team-adoption.md#t--tooling).

The rest of the plan makes reviewing a PR cheaper. This page is about the work *before* and *around* the PR. Part of every review is spent on problems that existed before the PR was opened:

- unused code that the reviewer still reads, and that agents copy patterns from,
- a flaky test that sends the PR back for another CI round,
- an authoring agent that makes the same mistake in every PR, which a reviewer then has to catch,
- old experiments that double the code paths a reviewer has to think about,
- CI that takes so long that every review round costs an hour.

None of these are review problems, but each one makes review slower or longer. Small scheduled agents can remove them at the source.

"Agent" here means a **scheduled job**, for example a GitHub Actions workflow on a cron schedule. It runs a deterministic tool to find candidates, then optionally an AI agent in headless mode for the judgement part, and produces a report, an issue or a PR. No agent platform is needed.

Agents don't replace the one-off readiness work ([`01-steps.md` Phase 3](01-steps.md#phase-3--repository-readiness-and-verification-before-the-pr)): lint fixes, the flaky-test burn-down, CI quick wins and test gaps are done by hand first. The agents keep that state from slipping back.

The dependency and CVE update bot ([§4.4](00-plan.md#44-dependency-and-cve-updates)) is already the first agent of this kind. The rules it follows apply to every agent on this page.

---

## 1. Rules every agent follows

Agent output goes through the same review process as everyone else's. An agent that produces more review work than it removes makes the problem worse, so every agent follows these rules:

| Rule | Why |
|---|---|
| **Deterministic first, model second.** A Go tool finds the candidates (unreachable code, CI failures and their logs). The model only does the judgement ("is this a real flake, and how do we fix it?") and writes the change | Cheaper and repeatable, and the candidate list can be checked. The same rule as for review: deterministic checks are code ([`00-plan.md` §3.1](00-plan.md#31-deterministic-checks-are-code-not-review)) |
| **Small, single-purpose PRs** with the evidence in the description: tool output, before/after test runs. One test, one package or one switch per PR | Easy to review, easy to revert. One large agent PR costs more review than it saves |
| **Review budget.** At most 1–2 open PRs per agent per repo. While those are open, new findings go into the agent's report instead of new PRs | An agent that floods the queue increases review wait, which is the opposite of the goal |
| **No duplicates.** Output is labelled `agent/<name>`. Before creating anything, the agent checks open issues and PRs, and ones closed as "won't fix" | Rejected suggestions don't come back every week |
| **Same rules as for people.** Changed assertions, new `//nolint` directives, workflow files and sensitive paths still need a human ([§4](00-plan.md#4-where-humans-stay-in-the-loop)). The risk tier comes from the change, not from who wrote it | Agents get no shortcut past the review layers |
| **A named owner.** Each agent has an owner who answers for its PRs, tunes it and turns it off. Ownership rotates like area pairs | A scheduled job has no "directing person" ([`01-steps.md` 7.6](01-steps.md#phase-7--risk-tiered-human-approval)) unless we name one |
| **Kill switch.** A repo variable (`AGENT_<NAME>_ENABLED`) turns the agent off without a PR | Stopping a noisy agent must be instant |
| **Measured.** Each agent reports PRs opened, merged and closed unmerged, and findings confirmed or rejected. Below roughly 50 % acceptance over a month, it's retuned or switched off | Same as Loop A for the reviewer ([§6](00-plan.md#6-feedback-loops)) |
| **Report-only start.** A new agent first runs for a few weeks writing only its report. It opens PRs once its owner trusts the findings | Start strict, loosen with evidence |
| **Least privilege.** A GitHub App per org with write access to contents and pull requests. It pushes only to `agent/*` branches, never to the default branch | Repos span several orgs, and nothing is shared across them |
| **Vendored.** The workflow and the agent's prompt file are shared files copied by the sync tool (D8) | Same behaviour in every repo and org |

**Output kinds.** Each agent keeps **one report issue per repo**, updated in place on every run (current findings, trend, what it did). Individual findings become issues or PRs only when they're actionable.

**Workflow skeleton:**

```yaml
# .github/workflows/agent-<name>.yml (vendored, see D8)
name: agent-<name>
on:
  schedule:
    - cron: "0 4 * * 1"          # weekly, early Monday UTC
  workflow_dispatch: {}
concurrency:
  group: agent-<name>
  cancel-in-progress: false
permissions:
  contents: read
jobs:
  run:
    if: vars.AGENT_<NAME>_ENABLED == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
      - id: token
        uses: actions/create-github-app-token@v2
        with:
          app-id: ${{ vars.AGENT_APP_ID }}
          private-key: ${{ secrets.AGENT_APP_KEY }}
      # 1. Deterministic: find candidates, write JSON (envelope in §6)
      - run: go run <toolkit>/cmd/<name> find --out candidates.json
        env:
          GH_TOKEN: ${{ steps.token.outputs.token }}
      # 2. Judgement: the agent reads the candidates and its prompt, then
      #    updates the report issue and opens at most the budgeted PRs
      - run: <agent-cli> --prompt-file .github/agents/<name>.md --input candidates.json
        env:
          GH_TOKEN: ${{ steps.token.outputs.token }}
```

---

## 2. The agents

| # | Agent | Tier | Cadence | Output | What it removes from review |
|---|---|---|---|---|---|
| 2.1 | Flaky-test fixer | 🟢 Core | Weekly | Report, issues, fix PRs | Re-runs, extra review rounds, distrust of red CI |
| 2.2 | Agent-context tuner | 🟢 Core | Monthly | PR to `AGENTS.md` or lint config | Repeated author mistakes that reviewers catch |
| 2.3 | CI time and cost tracker | 🟡 Later | Weekly report, alert on regression | Report, issues | Slow review rounds, runner shortages |
| 2.4 | Dead-code remover | 🟡 Later | Monthly | Small removal PRs | Code nobody uses but everyone reads |
| 2.5 | Stale-experiment finder | 🟡 Later | Monthly or quarterly | Report; PRs after a human decision | Untested code paths behind switches |
| 2.6 | Docs-sync | 🟡 Later | Monthly, and before each release | Report, docs-note drafts, small docs fix PRs | User-facing features nobody documented; docs that describe what no longer exists |

### 2.1 Flaky-test fixer 🟢

This is the first agent to build. Every flaky failure costs a re-run, another wait and a bit more trust in red CI. [§4.3](00-plan.md#43-ci-friction-flaky-tests-changelog-conflicts-and-branch-updates) already defines how flaky tests are detected and quarantined. The one-off burn-down of today's flaky tests is readiness work ([`01-steps.md` 3.7](01-steps.md#phase-3--repository-readiness-and-verification-before-the-pr)). This agent keeps the re-run rate down after that.

**Weekly run:**

1. **Collect.** All CI runs from the last week, for PRs and the default branch, in every repo in scope. A failure counts as flaky when the same SHA later passed on a re-run, or when the same test fails on unrelated PRs. Failures are read per test from the failed jobs' logs (`gh run view --log-failed`), not per job (T21). Logs are kept for 90 days by default, which is enough for a weekly run, so no test-report artifacts are needed.
2. **Group** by test name and failure signature: the assertion message, panic or timeout, with names, IDs and timestamps stripped. Rank the groups by occurrences × CI time lost.
3. **Classify** each group:

   | Class | Typical signs in our operators | Action |
   |---|---|---|
   | Infrastructure | Lost runner, image pull or registry rate limit, kind cluster creation timeout, external download failing, disk full | One issue per cause for whoever owns CI. Fix it in CI: cached or pre-pulled images, a mirror, a retry of the *setup* step. Never retry the test itself |
   | Test bug | `Eventually` timeout too short for CI runners; assertion on a cached read right after a write; reliance on list order; shared state between specs; leftover objects from an earlier spec; name or port collisions in parallel runs; `time.Sleep`; plain `Expect` inside `Eventually` | Fix PR |
   | Product bug | Reconcile not idempotent; missing requeue; status update conflict not retried; data race found with `-race`; behaviour that depends on how fresh the informer cache is | Issue with the analysis and a proposed fix. PR only if it's small. Always `needs-human`, because it changes behaviour |
   | Unknown | Not enough data | Keep counting. If it recurs, a PR that adds logging to the test so the next failure explains itself |

4. **Reproduce and verify.** For unit and envtest tests, run the test before and after the fix with `go test -run '<name>' -count=<N> -race`, also under CPU pressure (`-cpu 1,2,4`), and put both failure rates in the PR. End-to-end suites on kind can't be run hundreds of times, so the PR says "fixed by analysis" and the agent watches for recurrence.
5. **Close the loop.** The `flaky-test` issue closes only after a few weeks without a recurrence (start with 2). A recurrence reopens it, and the next fix PR names the earlier one in `Follow-up to:`.

**Rules for the fix:**

- Never weaken an assertion to make a test pass. That's a changed assertion and needs a human anyway.
- A longer timeout needs a reason in the PR: what the step really waits for, and how long it takes in CI. Timeouts stay realistic ([shared guidelines §6](shared/go-operator-guidelines.md#6-tests)).
- Quarantining (skipping the test with a link to the issue) stays a human decision.
- A fix corrects the existing test. It doesn't add a new one ([shared guidelines §6](shared/go-operator-guidelines.md#6-tests), no standalone regression tests).

**Measure:** CI re-run rate (already in `review-metrics --ci`), CI minutes lost to flaky failures, number and age of open `flaky-test` issues.

### 2.2 Agent-context tuner 🟢

[§5.3](00-plan.md#53-how-these-docs-improve-over-time) improves the docs the *reviewer* reads. This agent looks at the same data from the *author's* side: which mistakes do authoring agents (and people) make again and again, so that a reviewer has to catch them? That's not a review problem. It's an implementation problem that makes every review longer than it needs to be.

**Monthly run:**

1. **Collect** signs of avoidable rework from PRs merged in the last month:
   - automated review findings resolved as `fixed` (Loop A). These are real mistakes the author made;
   - human review comments that led to a change;
   - CI failures on a PR's first push that the repo's check command would have caught locally: lint, a `make manifests generate` diff, a missing changelog fragment;
   - commits pushed after review started, and what they changed.
2. **Cluster** them by kind of mistake. A cluster needs at least 3 PRs to count.
3. **Pick the cheapest place that prevents it**, in the order of [§5.3](00-plan.md#53-how-these-docs-improve-over-time): a lint rule or CI check, then a step in the `AGENTS.md` "How to verify" section, then a convention line in `AGENTS.md`, then a worked example. If the mistake comes from a wrong or missing line in `AGENTS.md`, fix that line.
4. **Open one PR** with the change and the evidence, for example "seen in #12, #15 and #19: CRDs not regenerated after a change in `apis/`". Code owners approve it, because these docs are owner-written.

**Deterministic checks on every run** (no model needed):

- every command and `make` target named in `AGENTS.md` exists and runs;
- every path named in `AGENTS.md` exists;
- `AGENTS.md` stays under the agreed length (D2).

**Measure:** first-push CI failure rate, and `fixed` findings per PR for each category that got a fix. Both should fall.

### 2.3 CI time and cost tracker 🟡

Not a weekly fixer. This is mostly a metric with an alert, plus improvement work now and then. Faster CI means faster review rounds, and the tracker also shows when we start running out of runner capacity. The first measurement and the cheap wins are readiness work ([`01-steps.md` 3.8](01-steps.md#phase-3--repository-readiness-and-verification-before-the-pr)); the tracker keeps watching after that.

**Tracked per repo, per workflow and per job** (default branch and PRs separately):

| Metric | Why |
|---|---|
| Push to all required checks done (p50, p90) | What the author and reviewer actually wait for |
| Queue time (job created → started) | Rising queue time means too few runners: the "running out of resources" signal |
| Duration per job | Where the time goes |
| Runner minutes (billable where it applies) | Cost, and limits on hosted runners |
| Slowest tests (durations from the job logs) | Targets for speed-up |
| Cache hit rate (Go module and build cache, image layers) | A drop is usually a cheap fix |
| Superseded runs that weren't cancelled | Wasted minutes |

**Source:** the Actions API (runs and jobs, with created, started and completed times). It's an extension of `review-metrics --ci` or a sibling command in the same tool, with the same JSON output ([§6](#6-future-one-place-for-all-reports)).

**Alert:** the weekly report opens an issue when the p50 time-to-green or queue time rises more than 20 % over the previous 4 weeks. The agent compares job durations commit by commit on the default branch and names the first commit where the job became slower.

**Improvement backlog** (the first items are the readiness quick wins; the rest are done by people or with an agent when the numbers say so):

- `concurrency` with `cancel-in-progress` on PR workflows, so a new push cancels the old run.
- Go module and build caches. Build images once per run and share them between jobs.
- Split long end-to-end suites across a matrix. Run a scenario only when relevant paths change, but note that a workflow skipped by a `paths` filter leaves its required check pending forever. Use one always-running gate job that decides which suites run and reports the combined result.
- Reuse a kind cluster across suites where isolation allows; pre-pull images.
- Move checks that don't need a real cluster from end-to-end to envtest ([shared guidelines §6](shared/go-operator-guidelines.md#6-tests) prefers envtest anyway).

### 2.4 Dead-code remover 🟡

Unused code costs review time: reviewers and agents read it, refactors have to keep it compiling, and agents copy patterns from it.

**Tools:**

| Tool | Finds | Notes |
|---|---|---|
| `deadcode` (`golang.org/x/tools/cmd/deadcode`) | Functions unreachable from any `main`, across the whole program (`-test` counts tests as roots) | Also reports exported functions. See public packages below |
| `staticcheck` U1000 / golangci-lint `unused` | Unused unexported identifiers | Safe to act on |
| `unparam` (enabled in the pilot) | Unused parameters and results | |
| `go mod tidy` | Unused modules | |

**Rules:**

- **`internal/`:** nothing outside the module can import it, so unreachable code there is removed. One PR per package. Low risk when the check command passes.
- **Exported identifiers in public packages** (`pkg/`, `apis/`): other modules may use them. k8ssandra-operator imports cass-operator packages, for example. Before suggesting removal, the agent looks for external users with GitHub code search across our orgs and with pkg.go.dev "Imported by". Even with no users found, removal follows Go's deprecation path: a `// Deprecated:` comment in one release, removal in a later one, both in the changelog. Always a human decision.
- **Never touched:** CRD types in `apis/` (they are the API; changes there always need a human), generated files (`zz_generated.*`, `config/crd/bases`), code under build tags the analysis didn't include.
- **Known false positives:** code reached through registration or reflection (scheme registration, webhooks registered in `init`), code built only with end-to-end build tags, code used by `go:generate` or scripts. Run the analysis with the repo's build tags, and keep confirmed false positives in an allowlist file (`.deadcode-ignore`) so they aren't raised again.

**Report only, same agent:** environment variables and command-line flags that nothing reads, and `hack/` scripts referenced by nothing (Makefile, workflows, docs).

### 2.5 Stale-experiment finder 🟡

Features added behind a switch, never promoted or removed, and no longer tested. Each one doubles the code paths a reviewer must keep in mind, and an untested path breaks silently.

**What counts as a switch** in our operators: feature gates, environment variables that change behaviour, annotations on a custom resource that turn on alternative behaviour, alpha or optional CRD fields that only gate behaviour, operator config fields, and version gates for server versions that are no longer supported.

**Monthly or quarterly report per repo** (no PRs by default):

| Switch | Kind | Introduced | Default | Tested on / off | Last changed | Documented | Suggestion |
|---|---|---|---|---|---|---|---|
| `example.io/some-annotation` | Annotation | 2023-04, #412 | off | ✗ / ✓ | 2023-05 | no | remove or add test |

**How:**

- **Inventory** by deterministic search (annotation constants, `os.Getenv` / `os.LookupEnv`, feature-gate registry, CRD fields), plus the model reading the code for less regular patterns. An owner-provided list per repo seeds it (A6).
- **"Tested"** comes from coverage: run unit and envtest suites with `-coverprofile`, and end-to-end suites with a coverage-instrumented operator binary (`go build -cover`, `GOCOVERDIR`). Check whether the branches behind each switch run with the switch on and with it off.
- **Age and activity** from `git log`.

**Suggestions:** *promote* (make it the default and remove the switch), *remove* (delete the path), *add a test*, or *keep* (with a reason, recorded so the agent doesn't raise it again for a few months). Promoting or removing is a product decision ("should we build this?", [§4](00-plan.md#4-where-humans-stay-in-the-loop)), so a human decides. After that, the agent can open the PR.

**Also in the same report:**

- skipped and pending tests (`t.Skip`, `Skip(`, `PIt`, `XIt`, `Pending`), especially quarantined flaky tests whose issue is closed or has gone quiet;
- `TODO` and `FIXME` comments without an issue link ([shared guidelines §8](shared/go-operator-guidelines.md#8-change-hygiene): follow-ups get a ticket);
- blocks of commented-out code.

### 2.6 Docs-sync 🟡

The PR-level docs check ([`00-plan.md` §5.6](00-plan.md#56-user-facing-documentation-docs-sync)) makes sure new user-facing changes carry a docs note. This agent catches what that check can't, and gives the docs team one place to start from.

- **Missed features.** Compares the user-facing surface between the last release and the default branch: CRD fields (from the generated CRD schemas), annotations, operator flags, environment variables, status conditions, events and metrics. Anything new without a docs note, a docs link or a mention in the docs goes in the report. The diff is deterministic; the model only decides whether existing docs already cover an item, and drafts a docs note from the code and the PR. The PR's author confirms the draft.
- **Docs drift.** Docs and examples that mention fields, flags, `make` targets or behaviour that no longer exist. Example manifests are checked mechanically against a kind cluster (`kubectl apply --dry-run=server`); the agent handles the prose. The per-PR reviewer flags stale lines in a diff; this catches what slipped in earlier.
- **Release handoff.** Before a release, one issue for the docs team with every docs note since the previous release, plus the missed features above.

**Output:** the report issue; docs-note drafts for missed features; small fix PRs for drift in in-repo docs. What gets a full docs page is the docs team's decision.

---

## 3. Other candidates

Ideas that come from the automation and development side rather than from review. Add one only when a metric or a retro points at it.

- **Lint backlog burn-down** 🟢. When a new linter or rule is added and the existing violations are too many for one PR ([lint proposal](pilot/cass-operator/lint-proposal.md)), an agent fixes them per linter and per package in small mechanical PRs. Same agent: remove `//nolint` directives that no longer suppress anything (`nolintlint` reports them).
- **Toolchain and scaffolding drift** 🟡. Across a family of operators, the Go version, controller-gen, kustomize, golangci-lint, the envtest Kubernetes version and kubebuilder scaffolding drift apart. When a feature PR also bumps a tool, its diff fills with regenerated files and the reviewer has to separate the two. An agent bumps tools in their own PRs, with the rules of [§4.4](00-plan.md#44-dependency-and-cve-updates), keeps versions aligned across repos, and leaves feature PRs clean.
- **Kubernetes version support** 🟡. Adds new Kubernetes minor versions to the end-to-end matrix when kind images exist, drops end-of-life ones according to the support policy, and scans code and manifests for deprecated or removed API versions. A scheduled PR is cheaper than a user's bug report or a surprise failure on an unrelated PR.
- **Broken default branch** 🟡. When CI on the default branch turns red, the agent finds the commit that broke it (re-run, then bisect) and comments on that PR, mentioning its author, who owns the fix ([§4](00-plan.md#4-where-humans-stay-in-the-loop) accountability). Supports P6, "who watches after merge".
- **Test-suite hygiene** 🟡. Tests that duplicate each other, tests whose assertions can't fail, standalone regression tests that belong in the normal tests ([shared guidelines §6](shared/go-operator-guidelines.md#6-tests)), unit tests that are slow. Consolidation PRs, kept rare and small, because changes to existing tests always need a human.
- **Tech-debt fixer** 🟡. The automated reviewer lists pre-existing problems in its summary, each with a proposed `tech-debt` issue. A human triages the issue. When they label it for agents, the agent opens a small fix PR, which needs human approval like any other. Examples: an error that retries forever and should be terminal, or a missing status condition.
- **Stale PRs and branches** ⚪. Closes abandoned drafts and agent PRs after a notice, deletes merged branches. Little value for a small team; useful once several agents run.
- **Issue preparation** ⚪. For a new bug report: find the relevant code, try to reproduce it on kind, link recent changes in that area. A well-prepared issue leads to a smaller, more focused fix PR.

---

## 4. Rollout order

| Step | What | Depends on |
|---|---|---|
| 0 🟢 | Repository readiness by hand ([`01-steps.md` Phase 3](01-steps.md#phase-3--repository-readiness-and-verification-before-the-pr)): lint, flaky-test burn-down, CI quick wins, test gaps, deterministic checks | — |
| 1 🟢 | Flaky-test fixer: report-only for two runs, then fix PRs | 0, T21 |
| 2 🟢 | Agent-context tuner: deterministic `AGENTS.md` checks first, then the monthly PR | Finding outcomes recorded (Phase 5, Loop A) |
| 3 🟡 | CI time and cost tracker in `review-metrics`, weekly report | 0 (baseline) |
| 4 🟡 | Docs-sync | PR-level docs check ([`01-steps.md` 3.11](01-steps.md#phase-3--repository-readiness-and-verification-before-the-pr)), D10 |
| 5 🟡 | Dead-code remover, `internal/` only | — |
| 6 🟡 | Stale-experiment report | Coverage map ([`01-steps.md` 3.9](01-steps.md#phase-3--repository-readiness-and-verification-before-the-pr)), A6 |
| 7 🟡 | Others from [§3](#3-other-candidates), when a metric points at them | — |

---

## 5. Open questions

| ID | Question | Why it matters | Status | Answer / Decision |
|---|---|---|---|---|
| A1 | Where does the scheduled agent's model run, and how is cost capped per run (token or time budget)? Vendor choice is out of scope ([§9](00-plan.md#9-explicitly-out-of-scope-this-round)), a budget isn't | Weekly jobs across many repos add up | open | |
| A2 | One GitHub App per org for all agents, or one per agent? Should it be the same identity as the automated reviewer? | Permissions, and telling agent PRs apart in metrics | open | |
| A3 | Test reports: JUnit or `go test -json`? Is artifact retention long enough for trends, or does a weekly collector copy them into the metrics store? | 2.1 and 2.3 both depend on per-test history | decided | No test-report artifacts. 2.1 and 2.3 read per-test failures and durations from the job logs (kept 90 days by default) |
| A4 | Flaky end-to-end tests that can't be reproduced locally: how many quiet weeks close the issue? | Avoids closing on luck | open | |
| A5 | Public packages: is code search in our orgs plus pkg.go.dev "Imported by" enough evidence, given users we can't see? How long is the deprecation period? | Removing exported API breaks others | open | |
| A6 | Which switches count as experiments in each repo (annotations, env vars, gates)? Needs an owner's list for the pilot | Seeds the 2.5 inventory | open | |
| A7 | Review budget per agent: how many open agent PRs per repo is acceptable for the team? | Agents must not become the review bottleneck | open | |
| A8 | Where can reports from private orgs be aggregated without breaking the "nothing shared across orgs" rule? | Needed before the dashboard ([§6](#6-future-one-place-for-all-reports)) | open | |

---

## 6. Future: one place for all reports

Today each tool has its own output: `review-metrics` (text, JSON, markdown), each agent's report issue, `flaky-test` issues, and the CI numbers. Eventually we want **one running UI** that shows them together, per repo and per team:

- review metrics ([§2](00-plan.md#2-goals-and-how-well-measure-them)): wait times, rounds, review distribution, knowledge spread;
- CI: time-to-green, queue time, minutes, re-run rate, open flaky tests;
- agents: runs, PRs opened, merged and rejected, acceptance rate, cost per run;
- inventories: stale experiments, dead-code candidates, skipped tests;
- trends and thresholds that feed the monthly check (Loop D).

**Decide now, before any UI:** every tool and agent writes the same JSON envelope. The UI is then only a reader, and can come later without changing the tools. `review-metrics --format json` moves to it too.

```json
{
  "tool": "flaky-tests",
  "toolVersion": "<git sha>",
  "repo": "k8ssandra/cass-operator",
  "period": { "from": "2026-10-12", "to": "2026-10-19" },
  "generatedAt": "2026-10-20T04:00:00Z",
  "metrics": { "rerunRate": 0.18, "ciMinutesLost": 340 },
  "findings": [
    { "id": "flaky/TestScaleUp/eventually-timeout", "kind": "flaky-test", "status": "open",
      "title": "TestScaleUp times out waiting for Ready", "links": ["https://github.com/..."] }
  ]
}
```

**Options, cheapest first:**

| Option | How | Fits when |
|---|---|---|
| Markdown reports in a repo | A scheduled Action commits each week's JSON and a rendered markdown summary into a `reports` repo in each org | Now. No infrastructure |
| Terminal dashboard | A `gh` extension built with charm (bubbletea, lipgloss) that reads the JSON history | Fits our Go toolkit; good for whoever runs the monthly check |
| Static site | A Go command renders the JSON history into static HTML with charts, published to an internal host | When trends across many repos matter |
| Small service | A Go service with SQLite that ingests the JSON from every run, serves the UI and sends alerts | Several teams, many repos, alerts wanted |

**Constraint:** repos are spread across several orgs and private data must stay private. The aggregation point has to be somewhere every org may send its reports to (A8). The public review-guidelines site isn't that place.
