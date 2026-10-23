# Review Process Development Tracker

**Audience:** maintainers of this review-process project. \
**Purpose:** track shared investigations, implementation work and decisions needed to make the process usable by adopting teams. Team-specific inventories, targets and policy choices belong in the [adoption questionnaire](02-team-adoption.md) and each team's profile.

## How to use

- Assign an owner when work starts; `TBD` means unassigned. Use **Next step / validation** to identify the next concrete action, and link results, code or documentation in **Answer / Decision**.
- Status: `open` · `investigating` · `needs-discussion` · `decided` · `done`. `decided` means an approach is agreed; `done` means the required implementation or documentation is delivered and validation evidence is linked. A sketch or a decision alone does not establish that a workflow is working.
- When an investigation is decided, add a dated entry to the **Decision log** and update the affected plan, steps or templates. Keep implementation work in the table until it is done. The historical log below records decisions, not verified delivery.
- Use [`01-steps.md`](01-steps.md) for rollout dependencies. Original question IDs are preserved; `-P` identifies the process part of a question split between the two files. Maintenance-agent implementation questions (A-series) remain in [`05-maintenance-agents.md`](05-maintenance-agents.md).

## T — Tooling

| ID | Question | Why it matters | Next step / validation | Status | Owner | Answer / Decision |
|---|---|---|---|---|---|---|
| T3 | How do we run the automated reviewer as an Action? Can we get inline comments, severity tags, no repeats, diff-only re-review and a round cap — or do we build the convergence logic ourselves? | Phase 5 hinges on it | Prototype on pilot repo | open | TBD |  |
| T4-P | How should the automated security review complement existing scanners without duplicating their findings? Use the repo inventory from [T4](02-team-adoption.md#t--tooling). | Avoid duplicate noise | Compare reviewer output with existing scanner findings on the pilot | open | TBD |  |
| T10 | Periodic review over existing code — where do findings go (issues? a report?) | Catches old bugs without flooding | Prototype | open | TBD |  |
| T12 | Can we compute knowledge spread from the GitHub API (authors + reviewers per path per period)? | Silo report (K3) | Prototype with the baseline script | open | TBD |  |
| T14 | How do we implement team auto-assignment and request human review only after automated review completes? Support the assignment policy chosen by each team in [T14-A](02-team-adoption.md#t--tooling). | Pain #1 | Try on pilot | open | TBD |  |
| T17 | Detector for `needs-human`: can a workflow find, in the PR diff, (a) changed or removed assertions in existing test files (`Expect`, `g.Expect`, `require.`, `assert.` lines; deleted `It`/`Test` blocks; added `Skip`), (b) new `//nolint` directives, (c) sensitive paths, and then add the `needs-human` label that the state workflow enforces? | Makes the "humans approve these" rules enforceable instead of relying on someone noticing | Prototype on pilot | open | TBD |  |
| T18 | "Reviews before your own PRs" gate: implement as an Action on `pull_request` opened/ready (check the author's review requests older than one working day → label, hold auto-assignment, comment) plus a scheduled re-check that requests reviewers once the queue is clear? How are "working days" computed across timezones? | Bootstraps review culture without blocking merges | Prototype on pilot | open | TBD |  |
| T22 | Review-state workflow ([`00-plan.md` §4.2](00-plan.md#42-knowing-when-bots-are-done-review-state-labels)): which events can trigger it in Actions (check completion, review threads), what needs the scheduled fallback, and how is the round count stored (commit status description or a hidden comment marker)? Each automated reviewer reports as its own check run `ai-review/<name>`, so any vendor or several agents work | Humans know when to start | Prototype on pilot | open | TBD |  |
| T23 | Dependency bot: Renovate, Dependabot or our own scheduled agent? Which one can follow "latest patch", batching, re-scan and fragment generation in all our orgs (OSS and private)? | Removes repetitive CVE work | Try on pilot | open | TBD |  |
| T24 | How do we compute CI time to green, queue time, job duration and runner minutes from the Actions API in `review-metrics --ci`? Team targets and runner capacity are covered by [T24-A](02-team-adoption.md#t--tooling). | Every review round waits for CI; queue time shows when runners run out | Extend `review-metrics --ci` on the pilot | open | TBD |  |
| T25 | Coverage map: how do we merge unit, envtest and end-to-end coverage (`go build -cover`, `GOCOVERDIR`, `go tool covdata`) and report it per area? How do we show changed-line coverage on a PR without an external service? | Readiness step 3.9 and reviewer information 3.12 | Prototype on the pilot | open | TBD |  |

## R — Risk tiers and approval policy

| ID | Question | Why it matters | Next step / validation | Status | Owner | Answer / Decision |
|---|---|---|---|---|---|---|
| R2 | How do we implement tier assignment using path rules and author labels? Each team defines who may override a tier in [R2-A](02-team-adoption.md#r--risk-tiers-and-approval-policy). | Tier must come from the change | Prototype path rules | needs-discussion | TBD |  |

## D — Per-project documentation

| ID | Question | Why it matters | Next step / validation | Status | Owner | Answer / Decision |
|---|---|---|---|---|---|---|
| D5 | How do we check the reviewer follows `REVIEW.md`? | Unfollowed rules are worse than none | Monthly sample in retro | open | TBD |  |
| D8 | Sync tool: a command in our Go toolkit that copies the shared files (guidelines, PR template, standard workflows) into each listed repo and opens a PR. Run locally or on a schedule from one repo with a token per org? How do copies record their source version? | The reviewer only follows what it can read locally | Prototype with two repos in different orgs | open | TBD |  |

## M — Measurement

| ID | Question | Why it matters | Next step / validation | Status | Owner | Answer / Decision |
|---|---|---|---|---|---|---|
| M1 | Which metrics can the GitHub API give us directly? | Baseline and dashboard | Baseline script ([`04-templates.md` §8](04-templates.md#8-metrics-tool)) | open | TBD |  |
| M3 | Is the `[wrong]/[intentional]/[not-worth-it]` reply prefix low enough friction? | Data quality | Try in pilot | open | TBD |  |
| M6 | How do we run the eval set — manually or as an Action? | Routine use | Depends on T3 | open | TBD |  |

## Decision log

| ID | Date | Decision |
|---|---|---|
| — | 2026-10-17 | Team-specific context and practices live in [`teams/`](teams/); pilot material in [`pilot/`](pilot/); shared stack guidance in [`shared/`](shared/) |
| D1 | 2026-10-17 | Per-project context file is `AGENTS.md`, vendor-neutral because reviews may be run by models from several vendors. A one-line pointer file is added only for a tool that requires its own name |
| — | 2026-10-17 | No explainer library. Explainers are disposable AI aids; durable learning goes into hand-written developer docs |
| — | 2026-10-17 | Every step is tagged 🟢 Core / 🟡 Later / ⚪ At scale; small teams start with Core only |
| — | 2026-10-18 | The PR description never repeats what CI shows. The template asks only for follow-up links, manual testing, risk and revert |
| — | 2026-10-18 | The follow-up fix rate (`Follow-up to: #N`) replaces the revert rate as the main guardrail |
| — | 2026-10-18 | No separate pre-commit hooks where a single check command (e.g. `make test`) exists |
| — | 2026-10-18 | New `//nolint` directives need human approval, the same as changed test assertions |
| — | 2026-10-18 | GitHub is treated as a tool, with no plan-specific assumptions. Repos span several orgs, so shared files are vendored per repo (former T13 dropped) |
| — | 2026-10-18 | Metrics tooling is written in Go (`tools/review-metrics`), using the libraries `gh` uses |
| — | 2026-10-19 | No standalone regression tests. Bug fixes extend the normal tests; changes that don't alter behaviour add no tests. Reviewers flag test sprawl: Important when an existing test owning the behaviour was bypassed, Nit when no test owned it yet |
| — | 2026-10-20 | Repetitive, deterministic checks are scripts or programs that run in CI and in the local check command, never instructions to the LLM reviewer. `REVIEW.md` holds judgement only |
| — | 2026-10-20 | Repository readiness (lint, flaky tests, CI time, test gaps in critical areas, deterministic checks) is hand-made work done before automated review starts on a repo. Test coverage is a per-area map, with no percentage gate |
| — | 2026-10-20 | User-facing changes carry user docs or a docs note; where and in which format is set per repo in `AGENTS.md` |
| — | 2026-10-21 | New `//nolint` directives and changed test assertions are reported as a `Needs human:` marker in the review summary, not as findings. A finding is raised only when the PR description doesn't explain the change |
| — | 2026-10-21 | Trade-offs stated in the PR description are not flagged. The PR template gets an **Intentional trade-offs** field |
| — | 2026-10-21 | Pre-existing findings go in the review summary with a proposed `tech-debt` issue. An agent fixes one only after a human labels the issue for it, and its PR needs human approval |
| — | 2026-10-21 | The changelog category is not a review finding: it is a matter of opinion |
| T21 | 2026-10-21 | Flaky tests are detected from re-runs (`run_attempt > 1` that later passed on the same SHA), with per-test detail from failed jobs' logs. Logs are kept 90 days by default, enough for a weekly run, so no test-report artifacts (also A3). Issue destination is a team choice in [T21-A](02-team-adoption.md#t--tooling) |
