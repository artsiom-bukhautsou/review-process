# Team profile — our team

**Process owner:** TBD

## Context

| Item | Value | What it means for the plan |
|---|---|---|
| Team size and structure | 6 devs: 3 senior, 3 junior | Process must be light. Fold it into existing rituals; no new meetings. |
| Code hosting / CI | GitHub, GitHub Actions; repos spread over at least 4 orgs (fully OSS, mixed OSS and private, fully private) | No cross-org sharing apart from ticket references. Shared files are vendored per repo with a sync tool. GitHub is a tool; missing features are covered by our own workflows |
| Tooling language | Go (an existing toolkit of Go CLIs) | Our own tools are Go CLIs using the libraries `gh` uses |
| Revert culture | Bugs are fixed forward with a new PR, rarely reverted | Guardrail is the follow-up fix rate (`Follow-up to: #N`) |
| Merge churn | Low; no merge queue. `CHANGELOG.md` conflicts constantly | Changelog fragments; no "branch must be up to date" rule |
| Repos | Multiple repos, public and private; some work is OSS | OSS section ([`00-plan.md` §8](../00-plan.md#8-open-source-repositories)) applies to public repos: fork PRs, outside contributors, CLA |
| Stack | Many Go controller-runtime operators | Shared guidance in [`shared/go-operator-guidelines.md`](../shared/go-operator-guidelines.md), referenced from each repo |
| Compliance | None, only quality requirements | Zero-human merges for Low risk are allowed in principle, with evidence |
| Quality requirements | Mostly unwritten, held by the main maintainer; linting exists but doesn't capture "Effective Go"-type judgement | Harvest rules from past review comments ([`00-plan.md` §5.5](../00-plan.md#55-getting-knowledge-out-of-peoples-heads)) |
| Pain #1 | **Waiting for reviews** | Primary target |
| Pain #2 | **Knowledge concentrated in a few people** (not always one) | Secondary target; same root cause as #1 |
| Pilot repo | `k8ssandra/cass-operator` | See [`pilot/cass-operator/`](../pilot/cass-operator/) |

## Diagnosis

The senior developers are both the review bottleneck and the knowledge holders. The pilot baseline makes this concrete (see [`pilot/cass-operator/README.md`](../pilot/cass-operator/README.md)): in the last two years one maintainer authored about 85% of merged PRs and left the largest share of review comments, and the median wait for a first human review on that maintainer's PRs was close to four days. Reviews arrive slowly because few people feel able to review, and few people feel able to review because the knowledge sits with the people writing most of the code.

So the plan for this team has two halves that support each other:

1. **Unload the seniors.** Automated layers handle routine findings; juniors take the human review on Low/Medium changes; seniors see only what truly needs them.
2. **Grow the juniors on purpose.** Review and explaining back are learning tools; juniors earn ownership of areas step by step. Every area a junior can own is one less area waiting on a senior.

## How the ownership model (00-plan §7.1) maps onto this team

| Generic practice | In our team |
|---|---|
| Area pairs | Each repo/area: **senior primary + junior secondary** in CODEOWNERS. Three seniors, three juniors → each senior–junior pair covers a set of repos |
| Review to learn | **Juniors review seniors' Low/Medium PRs.** Every junior review leaves at least a summary or a question |
| Coaching, not gatekeeping | **Seniors review juniors' PRs** for approach and design; mechanical comments go into lint or `REVIEW.md` |
| Graduated ownership | A junior becomes "trusted" for an area after N authored + N reviewed PRs there with no reverts; agreed in the retro and logged |
| Explain it back before the first change | Junior reads up (developer docs, code, an AI explainer), then explains the area to the senior primary in their own words in 15 minutes. Gaps go into `docs/developer/`, written by the junior and reviewed by the senior |
| Pairing on High-risk work | Senior + junior pair on High-risk changes; the second High-risk approver rotates through the juniors |
| Process ownership | One senior owns the process, one junior co-owns it |

## Team-specific adjustments to the steps

- **Phase 2.2:** a junior drafts the pilot `AGENTS.md` from the draft in [`pilot/cass-operator/`](../pilot/cass-operator/); a senior reviews it.
- **Phase 2.3:** the main maintainer corrects the harvested `REVIEW.md` draft; a junior reviews it for clarity ("could I apply this rule?").
- **Phase 6.2:** juniors dig up eval cases from past reverts and bug fixes; a senior confirms each one.
- **Phase 6.6:** the monthly review-process check happens inside the existing team retro.
- **Phase 8.5:** the on-duty rotation after merge includes the juniors.
- **Pilot is OSS:** Phase 5 starts with T16 (fork PRs) answered before the reviewer workflow goes live.

## Team decisions

| Date | Decision | Why |
|---|---|---|
| 2026-10-17 | Scope: one team of 6 (3 senior, 3 junior) | Start small |
| 2026-10-17 | GitHub; multiple public and private repos, no monorepo | Per-repo docs + shared guidance for the operator family |
| 2026-10-17 | No compliance constraints, only quality requirements | Lighter approval paths possible with evidence |
| 2026-10-17 | Priorities: (1) review wait, (2) knowledge concentration | Team input |
| 2026-10-17 | Pilot repo: `k8ssandra/cass-operator` | Active, central, patterns transfer to our other controller-runtime operators |
| 2026-10-17 | One review slot per person per day, at a time that fits their own timezone | Team spans timezones; standup is at end of day for some |
| 2026-10-17 | `AGENTS.md`, not `CLAUDE.md` | Reviews may be run by models from several vendors |
| 2026-10-17 | Start with 🟢 Core steps only | Small team; add more only on evidence |
| 2026-10-18 | Dependency/CVE PRs automated (latest patch, re-scan, batch); auto-merge for Low after one clean month | Scanner churn is a large share of routine work |
| 2026-10-18 | `make test` is the local check; no extra pre-commit hooks | Already runs fmt, vet and lint |

## Team-specific open questions

| ID | Question | Status | Answer |
|---|---|---|---|
| OT1 | Which of the cass-operator reviewers and committers are part of this team, and which are outside maintainers or contributors? | open | [k8ssandra maintainers team](https://github.com/orgs/k8ssandra/teams/maintainers) |
| OT2 | How do seniors feel about juniors reviewing their code, and juniors about approving senior code? | needs-discussion | It’s common practice: juniors learn, offer fresh perspectives, and provide a sanity check. |
| OT3 | How do we split repos/areas into three senior–junior pairs? Rotate yearly? | needs-discussion | Assign senior and junior reviewers per PR to spread knowledge and responsibility, using [GitHub auto-assignment](https://docs.github.com/en/organizations/organizing-members-into-teams/managing-code-review-settings-for-your-team), with direct requests when needed. |
| OT4 | What number of authored + reviewed PRs makes a junior "trusted" for an area? | needs-discussion | 20 substantive PRs authored and merged, plus as many reviewed. |
| OT5 | Which other operators come right after the pilot? | open | K8ssandra operator: its similarity to cass-operator makes it a quick way to test the review plan’s portability. |
| OT6 | First-response target, and each person's review slot in their own timezone? | needs-discussion | |
| OT7 | Do we run the "reviews before your own PRs" gate for the first 1–2 months? | needs-discussion | Leaning yes, as a bootstrapping measure |
