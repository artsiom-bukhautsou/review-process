# Team Adoption Questionnaire

**Audience:** each team adopting the review process. \
**Purpose:** identify the team's repos, constraints, owners, targets and local choices before and during rollout. This is a reusable questionnaire; adoption progress belongs to each team.

## How to use

1. Start a profile using [`teams/README.md`](teams/README.md). Copy the relevant questions into its **Team-specific open questions** table, keeping their IDs. Answers already captured in the profile can be linked rather than repeated.
2. Track each copied question with a status (`open` · `investigating` · `needs-discussion` · `decided`), an answer and, when decided, a date. Record per-repo details in the profile or the corresponding `pilot/<repo>/` material.
3. Work through questions when the related phase in [`01-steps.md`](01-steps.md) needs them. Skip Later and At scale questions until those parts are adopted. A question is resolved for a team when its answer is recorded and agreed by the relevant owners.

The tables below have no shared status or answer: one team's answer does not resolve a question for another team. Shared implementation work and process decisions are tracked in [`03-process-tracker.md`](03-process-tracker.md). Original question IDs are preserved; `-A` identifies the adoption part of a question split between the two files.

GitHub is treated as just a tool. Where a feature is missing, the check moves into our own workflow or becomes a review criterion. Nothing is shared across orgs.

## T — Tooling

| ID | Question | Why it matters | How to find out |
|---|---|---|---|
| T1 | Which repos are in scope, which org is each in, and which CI does each use? | Settings and workflows are per repo, since nothing is shared across orgs | Phase 0.1–0.2 |
| T2 | Which of our rules can branch rules and CODEOWNERS enforce directly in each repo, and which need a workflow check instead (e.g. "bot approval never counts", `needs-human`)? | Every rule needs an enforcer, whether the platform or our automation | Try on pilot |
| T4 | What security scanning runs today in each repo (Dependabot, secret scanning, code scanning)? Which results must the automated security review complement? See [T4-P](03-process-tracker.md#t--tooling). | Avoid duplicate noise | Check repo Security tab per repo |
| T5 | (⚪ at scale) When does churn justify a merge queue? Signal: "passed separately, broke together" incidents on main | Avoid adding a merge queue before there's a problem | Watch post-merge CI on main |
| T6 | Per repo: how does it deploy, are there flags, how do we roll back, how long does it take? | Determines when lighter approval is safe | Phase 0.3 |
| T7 | Does each repo have one-command test/lint/types/build? | Verification before PR | Repo inventory |
| T8 | Which local review tooling do authors use (assistant review command, subagent, security review command)? | Phase 3.3 | Try on pilot |
| T9 | UI repos: how do we verify visual changes? | UI bugs slip past unit tests | Per front-end stack |
| T11 | Where do metrics live? Scheduled Action writing a markdown report to a repo is the cheapest option | Feedback loops | Decide after first script |
| T14-A | Which review assignment policy should we use: round robin or load balance? Which team should receive requests? Coordinate with [T14](03-process-tracker.md#t--tooling). | Spreads reviews according to the team's needs | Agree with the team before enabling auto-assignment |
| T15 | Stale-PR nudges: where should reminders be sent? | Visibility of waiting PRs | Slack DM to the relevant team member, or a message in the team’s Slack channel for group reminders. |
| T19 | Does merging base into a PR (no conflicts) dismiss existing approvals with our branch rules? | Decides how costly a branch update is | Test on pilot |
| T20 | Changelog fragments: `changie` or a small command in our Go toolkit? The fragment format must render the repo's `[KIND] [#N](link) text` lines in the required order | Ends CHANGELOG conflicts | Try `changie` on the pilot |
| T21-A | In which repo should flaky-test issues be created? Use the detection approach recorded in [T21](03-process-tracker.md#decision-log). | Detection needs a destination | Agree during the pilot flaky-test burn-down |
| T16 | Public repos: how does the automated reviewer handle fork PRs without exposing secrets — label-gated `pull_request_target` reading the diff only, a `workflow_run` follow-up, or maintainer-triggered? | Security of the review pipeline | Prototype on a public pilot |
| T24-A | What are our CI time targets and action thresholds? Does any repo need more or larger runners? Use the measurements from [T24](03-process-tracker.md#t--tooling). | Turns CI measurements into local capacity and performance decisions | Set targets from the pilot baseline |

## R — Risk tiers and approval policy

| ID | Question | Why it matters | How to find out |
|---|---|---|---|
| R1 | What are our sensitive paths/areas, per repo? | Always-human paths | 30-min workshop with owners and non-owners together |
| R2-A | Who can override a risk tier? Use the tier-assignment mechanism from [R2](03-process-tracker.md#r--risk-tiers-and-approval-policy). | Defines who may override the assigned tier | Agree with repo owners before enabling tier assignment |
| R3 | Graduated ownership: what earns someone "trusted" status for an area (e.g. N authored + N reviewed PRs, no reverts)? Who decides? | Lighter paths without anyone staying an outsider by default | Discuss in retro |
| R4 | What evidence moves Low from 1 human to 0? | Concrete loosening bar | Decide after Phase 6 data |
| R5 | Circuit breaker threshold for tightening tiers? | Safety valve | Discuss |
| R6 | Written quality requirements: what are they, and are they in a form the reviewer and humans can both apply? | Unwritten bars can't be enforced by tools or by new reviewers | Mine past comments + owner interview ([`00-plan.md` §5.5](00-plan.md#55-getting-knowledge-out-of-peoples-heads)) |
| R7 | Do we have compliance/regulatory requirements that mandate human review? | Could rule out zero-human merges | Ask |
| R8 | Outside contributors (OSS, other teams): which approvals do they always need? | Trust and direction | Trusted team member approval |

## D — Per-project documentation

| ID | Question | Why it matters | How to find out |
|---|---|---|---|
| D2 | Length limits? | Short docs get followed | Agree, revisit after pilot |
| D3 | Doc prune owner and cadence? Suggest: area buddy pair, quarterly | Docs rot without an owner | Discuss |
| D6 | Developer docs: is `docs/developer/` the right home in every repo, and who reviews additions written after an "explain it back" session? | Learning goes into maintained docs, not an explainer library | Discuss |
| D7 | Shared conventions for all repos: which repo is the source of truth? (Repos span several orgs, so it's copied, not referenced) | Avoid drift between copies | Discuss |
| D9 | Which rules are mechanical (→ linter) vs judgement (→ `REVIEW.md`)? Re-check after each harvest | Prefer mechanical enforcement | Per harvest |
| D10 | User docs, per repo: what counts as user-facing, where docs go (in-repo `docs/`, a separate docs repo, the docs team's tracker), which format (full page or docs note), and who in the docs team picks them up? Does the docs note travel in the changelog fragment? | Every user-facing feature reaches the docs team | Ask the docs team; decide for the pilot |

## M — Measurement

| ID | Question | Why it matters | How to find out |
|---|---|---|---|
| M2 | Do we use the `Follow-up to: #N` field consistently? Also, can past follow-ups be found retroactively, e.g. from `[BUGFIX]` CHANGELOG entries that name the PR that introduced the bug? | Guardrail accuracy, since we fix forward rather than revert | Check the first month's PRs; backfill the pilot |
| M4 | Pilot thresholds (fixed share, rounds-to-pass, survey)? | Exit criteria | Set after 2 weeks of data |
| M5 | How many past follow-up fixes or bugs are usable as eval cases? | Eval effort | Search history |
| M7 | Rollback time target per repo? | Phase 8 exit | Per repo |
| M8 | Which CI re-run rate counts as "ready" (e.g. under 5 % of PR head commits)? | Phase 3 exit | Set from the pilot baseline |

## P — Process and people

| ID | Question | Why it matters | How to find out |
|---|---|---|---|
| P1 | Which repo is the pilot? Suggest: active, decent tests, review load concentrated on few people, not the most critical | Everything starts there | Phase 0.6 |
| P2 | Process owner and co-owner from a different part of the team? | Needs drivers; co-ownership is itself anti-silo | Discuss |
| P3 | Soft PR size limit — what number? | Small PRs are faster and easier to hand to less familiar reviewers | Look at current size distribution |
| P4 | Can each product be used internally / in staging before release? | Early detection | Per repo |
| P5 | How do experts feel about less familiar people reviewing their code, and vice versa? | Adoption risk; psychological safety | Team conversation |
| P6 | Who watches after merge? Is there an on-call/on-duty concept today? | Detection layer | Discuss |
| P7 | Where do "should we build this?" discussions happen — draft PR, issue, chat thread? | Early direction changes | Discuss |
| P8 | First-response target, and does everyone have a daily review slot that fits their own timezone? | Pain #1 | Discuss |

## K — Knowledge sharing / anti-silo

| ID | Question | Why it matters | How to find out |
|---|---|---|---|
| K1 | "Explain it back" norm: when is it required (first change in an area? any High-risk change?) and how do we notice gaps it reveals in the developer docs? | Makes learning durable | Try for one quarter |
| K2 | Walkthrough slot: fortnightly 20 min inside an existing meeting? | Context sharing | Trial for a month |
| K3 | What knowledge-spread number triggers action (e.g. area with <2 recent people)? | Make silos visible | Depends on T12 |
| K4 | Area pairs: how do we split areas among people? Rotate pairs periodically? | Spread critical knowledge | Use silo map from Phase 0 |
