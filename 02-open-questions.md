# Open Questions & Investigation TODO

**Scope:** generic questions. Team-specific answers and questions live in `teams/<team>.md`. \
**How to use:** each item has an ID referenced from [`01-steps.md`](01-steps.md). Fill in **Answer / Decision** as we go; move resolved items to the bottom with a date. Status: `open` · `investigating` · `needs-discussion` · `decided`.

---

## T — Tooling

GitHub is treated as just a tool. Where a feature is missing, the check moves into our own workflow or becomes a review criterion. Nothing is shared across orgs.

| ID | Question | Why it matters | How to find out | Status | Answer / Decision |
|---|---|---|---|---|---|
| T1 | Which repos are in scope, which org is each in, and which CI does each use? | Settings and workflows are per repo, since nothing is shared across orgs | Phase 0.1–0.2 | open | |
| T2 | Which of our rules can branch rules and CODEOWNERS enforce directly in each repo, and which need a workflow check instead (e.g. "bot approval never counts", `needs-human`)? | Every rule needs an enforcer, whether the platform or our automation | Try on pilot | open | |
| T3 | How do we run the automated reviewer as an Action? Can we get inline comments, severity tags, no repeats, diff-only re-review and a round cap — or do we build the convergence logic ourselves? | Phase 5 hinges on it | Prototype on pilot repo | open | |
| T4 | What security scanning runs today (Dependabot, secret scanning, code scanning)? How does an automated security review sit next to it? | Avoid duplicate noise | Check repo Security tab per repo | open | |
| T5 | (⚪ at scale) When does churn justify a merge queue? Signal: "passed separately, broke together" incidents on main | Avoid adding a merge queue before there's a problem | Watch post-merge CI on main | open | |
| T6 | Per repo: how does it deploy, are there flags, how do we roll back, how long does it take? | Determines when lighter approval is safe | Phase 0.3 | open | |
| T7 | Does each repo have one-command test/lint/types/build? | Verification before PR | Repo inventory | open | |
| T8 | Which local review tooling do authors use (assistant review command, subagent, security review command)? | Phase 3.3 | Try on pilot | open | |
| T9 | UI repos: how do we verify visual changes? | UI bugs slip past unit tests | Per front-end stack | open | |
| T10 | Periodic review over existing code — where do findings go (issues? a report?) | Catches old bugs without flooding | Prototype | open | |
| T11 | Where do metrics live? Scheduled Action writing a markdown report to a repo is the cheapest option | Feedback loops | Decide after first script | open | |
| T12 | Can we compute knowledge spread from the GitHub API (authors + reviewers per path per period)? | Silo report (K3) | Prototype with the baseline script | open | |
| T14 | Team auto-assignment: round robin vs load balance? Can human review be requested only after the automated review completes? | Pain #1 | Try on pilot | open | |
| T15 | Stale-PR nudge: where to post (Slack/Teams/GitHub Discussions)? What chat tool do we use? | Visibility of waiting PRs | Discuss | open | |
| T17 | Detector for `needs-human`: can a workflow find, in the PR diff, (a) changed or removed assertions in existing test files (`Expect`, `g.Expect`, `require.`, `assert.` lines; deleted `It`/`Test` blocks; added `Skip`), (b) new `//nolint` directives, (c) sensitive paths, and then add the `needs-human` label that the state workflow enforces? | Makes the "humans approve these" rules enforceable instead of relying on someone noticing | Prototype on pilot | open | |
| T18 | "Reviews before your own PRs" gate: implement as an Action on `pull_request` opened/ready (check the author's review requests older than one working day → label, hold auto-assignment, comment) plus a scheduled re-check that requests reviewers once the queue is clear? How are "working days" computed across timezones? | Bootstraps review culture without blocking merges | Prototype on pilot | open | |
| T19 | Does merging base into a PR (no conflicts) dismiss existing approvals with our branch rules? | Decides how costly a branch update is | Test on pilot | open | |
| T20 | Changelog fragments: `changie` or a small command in our Go toolkit? The fragment format must render the repo's `[KIND] [#N](link) text` lines in the required order | Ends CHANGELOG conflicts | Try `changie` on the pilot | open | |
| T21 | Flaky-test detection: from CI re-runs (`run_attempt > 1` that later passed on the same SHA), or from test reports (JUnit) across runs? Where does the `flaky-test` issue get created? | Re-runs are the main slowdown | Start with re-run data from `review-metrics --ci` | open | |
| T22 | Review-state workflow ([`00-plan.md` §4.2](00-plan.md#42-knowing-when-bots-are-done-review-state-labels)): which events can trigger it in Actions (check completion, review threads), what needs the scheduled fallback, and how is the round count stored (commit status description or a hidden comment marker)? Each automated reviewer reports as its own check run `ai-review/<name>`, so any vendor or several agents work | Humans know when to start | Prototype on pilot | open | |
| T23 | Dependency bot: Renovate, Dependabot or our own scheduled agent? Which one can follow "latest patch", batching, re-scan and fragment generation in all our orgs (OSS and private)? | Removes repetitive CVE work | Try on pilot | open | |
| T16 | Public repos: how does the automated reviewer handle fork PRs without exposing secrets — label-gated `pull_request_target` reading the diff only, a `workflow_run` follow-up, or maintainer-triggered? | Security of the review pipeline | Prototype on a public pilot | open | |

## R — Risk tiers and approval policy

| ID | Question | Why it matters | How to find out | Status | Answer / Decision |
|---|---|---|---|---|---|
| R1 | What are our sensitive paths/areas, per repo? | Always-human paths | 30-min workshop with owners and non-owners together | needs-discussion | |
| R2 | Tier assignment: path rules + author label to start? Who can override a tier? | Tier must come from the change | Prototype path rules | needs-discussion | |
| R3 | Graduated ownership: what earns someone "trusted" status for an area (e.g. N authored + N reviewed PRs, no reverts)? Who decides? | Lighter paths without anyone staying an outsider by default | Discuss in retro | needs-discussion | |
| R4 | What evidence moves Low from 1 human to 0? | Concrete loosening bar | Decide after Phase 6 data | needs-discussion | |
| R5 | Circuit breaker threshold for tightening tiers? | Safety valve | Discuss | needs-discussion | |
| R6 | Written quality requirements: what are they, and are they in a form the reviewer and humans can both apply? | Unwritten bars can't be enforced by tools or by new reviewers | Mine past comments + owner interview ([`00-plan.md` §5.5](00-plan.md#55-getting-knowledge-out-of-peoples-heads)) | open | |
| R7 | Do we have compliance/regulatory requirements that mandate human review? | Could rule out zero-human merges | Ask | open | |
| R8 | Outside contributors (OSS, other teams): which approvals do they always need? | Trust and direction | Discuss | needs-discussion | |

## D — Per-project documentation

| ID | Question | Why it matters | How to find out | Status | Answer / Decision |
|---|---|---|---|---|---|
| D1 | `CLAUDE.md`, `AGENTS.md`, or both? | Tool-neutral naming | Check what our tools read | decided | `AGENTS.md`, because reviewers may be models from different vendors. A one-line pointer file is added only for a tool that requires its own name |
| D2 | Length limits? | Short docs get followed | Agree, revisit after pilot | open | |
| D3 | Doc prune owner and cadence? Suggest: area buddy pair, quarterly | Docs rot without an owner | Discuss | open | |
| D5 | How do we check the reviewer follows `REVIEW.md`? | Unfollowed rules are worse than none | Monthly sample in retro | open | |
| D6 | Developer docs: is `docs/developer/` the right home in every repo, and who reviews additions written after an "explain it back" session? | Learning goes into maintained docs, not an explainer library | Discuss | open | |
| D7 | Shared conventions for all repos: which repo is the source of truth? (Repos span several orgs, so it's copied, not referenced) | Avoid drift between copies | Discuss | open | |
| D8 | Sync tool: a command in our Go toolkit that copies the shared files (guidelines, PR template, standard workflows) into each listed repo and opens a PR. Run locally or on a schedule from one repo with a token per org? How do copies record their source version? | The reviewer only follows what it can read locally | Prototype with two repos in different orgs | open | |
| D9 | Which rules are mechanical (→ linter) vs judgement (→ `REVIEW.md`)? Re-check after each harvest | Prefer mechanical enforcement | Per harvest | open | |

## M — Measurement

| ID | Question | Why it matters | How to find out | Status | Answer / Decision |
|---|---|---|---|---|---|
| M1 | Which metrics can the GitHub API give us directly? | Baseline and dashboard | Baseline script ([`03-templates.md` §8](03-templates.md#8-metrics-tool)) | open | |
| M2 | Do we use the `Follow-up to: #N` field consistently? Also, can past follow-ups be found retroactively, e.g. from `[BUGFIX]` CHANGELOG entries that name the PR that introduced the bug? | Guardrail accuracy, since we fix forward rather than revert | Check the first month's PRs; backfill the pilot | open | |
| M3 | Is the `[wrong]/[intentional]/[not-worth-it]` reply prefix low enough friction? | Data quality | Try in pilot | open | |
| M4 | Pilot thresholds (fixed share, rounds-to-pass, survey)? | Exit criteria | Set after 2 weeks of data | open | |
| M5 | How many past follow-up fixes or bugs are usable as eval cases? | Eval effort | Search history | open | |
| M6 | How do we run the eval set — manually or as an Action? | Routine use | Depends on T3 | open | |
| M7 | Rollback time target per repo? | Phase 8 exit | Per repo | open | |

## P — Process and people

| ID | Question | Why it matters | How to find out | Status | Answer / Decision |
|---|---|---|---|---|---|
| P1 | Which repo is the pilot? Suggest: active, decent tests, review load concentrated on few people, not the most critical | Everything starts there | Phase 0.6 | open | |
| P2 | Process owner and co-owner from a different part of the team? | Needs drivers; co-ownership is itself anti-silo | Discuss | open | |
| P3 | Soft PR size limit — what number? | Small PRs are faster and easier to hand to less familiar reviewers | Look at current size distribution | open | |
| P4 | Can each product be used internally / in staging before release? | Early detection | Per repo | open | |
| P5 | How do experts feel about less familiar people reviewing their code, and vice versa? | Adoption risk; psychological safety | Team conversation | needs-discussion | |
| P6 | Who watches after merge? Is there an on-call/on-duty concept today? | Detection layer | Discuss | open | |
| P7 | Where do "should we build this?" discussions happen — draft PR, issue, chat thread? | Early direction changes | Discuss | needs-discussion | |
| P8 | First-response target, and does everyone have a daily review slot that fits their own timezone? | Pain #1 | Discuss | needs-discussion | |

## K — Knowledge sharing / anti-silo

| ID | Question | Why it matters | How to find out | Status | Answer / Decision |
|---|---|---|---|---|---|
| K1 | "Explain it back" norm: when is it required (first change in an area? any High-risk change?) and how do we notice gaps it reveals in the developer docs? | Makes learning durable | Try for one quarter | needs-discussion | |
| K2 | Walkthrough slot: fortnightly 20 min inside an existing meeting? | Context sharing | Trial for a month | needs-discussion | |
| K3 | What knowledge-spread number triggers action (e.g. area with <2 recent people)? | Make silos visible | Depends on T12 | open | |
| K4 | Area pairs: how do we split areas among people? Rotate pairs periodically? | Spread critical knowledge | Use silo map from Phase 0 | needs-discussion | |

---

## Resolved

_(Generic decisions only. Team-specific decisions go in the team profile.)_

| ID | Date | Decision |
|---|---|---|
| — | 2026-10-17 | Team-specific context and practices live in [`teams/`](teams/); pilot material in [`pilot/`](pilot/); shared stack guidance in [`shared/`](shared/) |
| — | 2026-10-17 | Per-project context file is `AGENTS.md`, vendor-neutral because reviews may be run by models from several vendors (former D1) |
| — | 2026-10-17 | No explainer library. Explainers are disposable AI aids; durable learning goes into hand-written developer docs |
| — | 2026-10-17 | Every step is tagged 🟢 Core / 🟡 Later / ⚪ At scale; small teams start with Core only |
| — | 2026-10-18 | The PR description never repeats what CI shows. The template asks only for follow-up links, manual testing, risk and revert |
| — | 2026-10-18 | The follow-up fix rate (`Follow-up to: #N`) replaces the revert rate as the main guardrail |
| — | 2026-10-18 | No separate pre-commit hooks where a single check command (e.g. `make test`) exists |
| — | 2026-10-18 | New `//nolint` directives need human approval, the same as changed test assertions |
| — | 2026-10-18 | GitHub is treated as a tool, with no plan-specific assumptions. Repos span several orgs, so shared files are vendored per repo (former T13 dropped) |
| — | 2026-10-18 | Metrics tooling is written in Go (`tools/review-metrics`), using the libraries `gh` uses |
| — | 2026-10-19 | No standalone regression tests. Bug fixes extend the normal tests; changes that don't alter behaviour add no tests. Reviewers flag test sprawl as Important |
