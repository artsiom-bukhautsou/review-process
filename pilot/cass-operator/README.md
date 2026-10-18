# Pilot — k8ssandra/cass-operator

**Status:** Phase 0 (baseline) done from public data; Phase 2 drafts ready for owner review.

| File | What it is | Next action |
|---|---|---|
| `README.md` | Baseline, findings, pilot plan | Review findings |
| `AGENTS.md` | Draft project context file for the repo | Owner confirms marked lines; a non-owner opens the PR |
| `REVIEW.md` | Draft review guidance harvested from past review comments | Owner corrects; then commit |
| `lint-proposal.md` | golangci-lint changes that replace recurring review comments | Agree; enable and fix the codebase in per-linter PRs |
| `../../shared/go-operator-guidelines.md` | Rules shared by all our Go operators | Owner corrects; decide where it lives (D8) |

## Why this repo

It is active, central to our work, and its patterns (controller-runtime reconcile loop, webhooks, envtest, Ginkgo e2e on kind) transfer directly to our other operators. It is public, so the pilot also tests the OSS considerations in `00-plan.md` §8.

## Baseline (merged PRs, Oct 2024 – Oct 2026)

Source: GitHub REST API, 98 merged PRs, their reviews, review comments and changed files. Bot reviews and self-reviews excluded. Reproduce with `tools/review-metrics k8ssandra/cass-operator --since 2024-10-01`.

| Metric | Value |
|---|---|
| Merged PRs | 98 (about 4 a month) |
| Time to first human review, median / p90 | **~89 h / ~616 h** (about 3.7 days / 26 days) |
| Open → merge, median / p90 | ~122 h / ~978 h |
| Merged with no human review by someone other than the author | 15 (15%) |
| Human approvals per PR | 0: 17, 1: 79, 2: 2 |
| Authorship | **One maintainer authored 83 of 98 (85%)**. The rest came from 6 other people |
| Median first-review wait, main maintainer's PRs | ~90 h |
| Median first-review wait, everyone else's PRs | ~65 h |
| PR size (lines changed), median / p90 | 124 / 561 |
| **PRs updated from base at least once** | **45%** (mean 0.69 updates per PR, max 7). Mostly rebases: 53 force-pushes vs 15 merge updates |
| **CI re-runs** (Jul–Oct 2026, 21 PRs) | **24% of PR head commits needed a re-run**; 4 of 5 were *Kind Integration Tests* |
| Recent trend | Since 2026-07 the median first-review wait rose to ~236 h |
| Revert PRs | 0. The team fixes forward, so breakage shows up as later `[BUGFIX]` PRs. These will be tracked through `Follow-up to: #N` from now on (M2) |

### Knowledge spread by area

People who authored or reviewed in each area, over the same period:

| Area | PRs with main maintainer involved | People involved | Note |
|---|---|---|---|
| `pkg/reconciliation` | ~57 | 11 | Next most involved person: ~26 PRs |
| `internal/controllers` | ~27 | 7 | |
| `apis/cassandra` | ~19 | 6 | API changes: highest upgrade risk |
| `pkg/httphelper` | ~18 | 7 | |
| `.github/workflows` | ~16 | 6 | Sensitive path, should always need a human |
| `config/crd` | ~16 | 4 | Generated, but review needs to know why it changed |
| `internal/webhooks` | ~8 | 4 | Thinnest coverage of a critical area |
| `apis/control` | ~8 | 4 | |

"People involved" overstates real knowledge: in every area most activity comes from the main maintainer plus one or two regular approvers.

## Findings

1. **Review wait is the dominant cost.** A median of nearly four days to the first human review on the main maintainer's PRs means most of a PR's life is spent waiting. The automated layers plus team auto-assignment target exactly this.
2. **The silo is authorship, not only review.** Reviews are spread over a few approvers, but 85% of the code comes from one person, so approvers often review code in areas they don't author in. The ownership model (`00-plan.md` §7.1) needs other people *writing* code in `pkg/reconciliation` and `internal/webhooks`, not only reviewing it.
3. **The quality bar is consistent but unwritten.** About 160 of the main maintainer's top-level review comments fall into a small number of recurring themes: reconcile-loop correctness, Kubernetes semantics, upgrade safety, test meaningfulness, Go idiom and change hygiene. They are now drafted in `REVIEW.md` and the shared guidelines.
4. **Several recurring comments are lintable.** `revive` and `govet` are already configured but disabled (`lint-proposal.md`).
5. **15% of PRs merged without a second person.** This is fine for genuinely low-risk changes, but today it happens without a written rule. Risk tiers (Phase 7) make it explicit.
6. **The PR template lacks risk, rollback and follow-up links.** It has a checklist (tests, docs, CHANGELOG, CLA) but no **Risk**, **How to revert** or **Follow-up to** field. Without that field, fix-forward PRs can't be linked to the PR that caused the problem.
7. **The OSS setup is already partly there.** A linked-issue check, OSV scanner and govulncheck run in CI. CODEOWNERS is a single team, with no path-specific owners.
8. **Friction is large.** Nearly half of all PRs needed a branch update, mostly rebases, which also break review threads. A quarter needed CI re-runs, almost all in the Kind integration tests. Changelog fragments (1.6b), the update-only-on-conflict norm (1.6d) and flaky-test hunting in the Kind suites (1.6c) target exactly this.

## Pilot plan (deviations from `01-steps.md` only)

| Step | Pilot-specific action |
|---|---|
| 1.1 | Extend the existing PR template: keep the CLA and issue lines and the "manually tested" flag (as a **Manual testing** section); add **Follow-up to**, **Risk**, **How to revert**. Drop the "automated tests added" checkbox, since CI and the diff show it |
| 1.6b | Changelog fragments: `CHANGELOG.md` is the main source of conflicts and branch updates here |
| 1.3 | Enable conversation resolution and stale-approval dismissal on `master` |
| 1.4 | CODEOWNERS: `.github/workflows/`, `apis/`, `internal/webhooks/` → named owners; the rest → maintainers team |
| 2.x | Owner corrects `REVIEW.md` and the shared guidelines; a non-owner commits `AGENTS.md` |
| 2.4 | Enable `revive`, `govet` and the additions, fixing the code in the same PRs (`lint-proposal.md`) |
| 4.x | Pick `internal/webhooks` and one part of `pkg/reconciliation` (e.g. decommission) as the first areas where a second person **authors** changes, with the "explain it back" norm; gaps found go into `docs/developer/` (starting with replacing the outdated `kind.md` with `hack/cluster.sh` usage) |
| 5.1b | Decide fork-PR handling (T16) before the reviewer workflow goes live, since outside contributors open PRs here |
| 6.2 | Eval set: use past `[BUGFIX]` entries in the CHANGELOG that link to the PR that introduced the bug, plus bugs found in review (the harvested comments show where), confirmed by the main maintainer |

## How the harvest was done (repeatable for other repos)

1. Fetch all PR review comments: `GET /repos/{owner}/{repo}/pulls/comments` (paginate with `page=`).
2. Keep top-level comments (no `in_reply_to_id`) from the expert, longer than a few words.
3. Group them into themes, then sort each theme into **lint** (mechanical) or **REVIEW.md / shared guidelines** (judgement).
4. The expert corrects the draft. That takes about 30 minutes instead of writing from scratch.
