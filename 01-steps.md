# Step List — Phased Rollout

**Status:** Draft v0.3 — generic. Team-specific choices live in `teams/<team>.md`; pilot-specific material in `pilot/<repo>/`.
**How to use:** each phase has steps (checkboxes) and **exit criteria**. A phase is done when its exit criteria are met, not when every box is ticked. Steps marked 🔎 depend on an open question in `02-open-questions.md` (ID in brackets).


**Start small.** Every step is tagged with when it's worth doing:

- 🟢 **Core.** Start here. This is enough for a small team, and it covers review wait, the knowledge silo and basic automated review.
- 🟡 **Later.** Add once the core runs smoothly, usually after month 2–3, or when a metric shows you need it.
- ⚪ **At scale.** Only for high PR volume, many repos or several teams. Skip until there's evidence you need it.

A small team can run the 🟢 steps alone and stop there. See `00-plan.md` §0.1 for the reasoning.

Phases are ordered by dependency, but some run in parallel. We start on **one pilot repo**, then widen. **Phase 1 includes quick wins** for review wait that need no new tooling.

```
Phase 0  Discover & baseline                       (~1–2 weeks)
   │
Phase 1  Foundations + review-wait quick wins      (~2 weeks)
   │
   ├── Phase 2  Per-project docs (AGENTS.md / REVIEW.md)    ┐
   ├── Phase 3  Verification before the PR                  ├ in parallel
   └── Phase 4  Knowledge sharing & ownership model ┘
   │
Phase 5  Automated review on PRs (pilot repo)
   │
Phase 6  Feedback loops & eval set
   │
Phase 7  Risk-tiered approval
   │
Phase 8  Merge queue, flags, rollback
   │
Phase 9  Widen to all repos (ongoing)
```

> Knowledge sharing runs early: it attacks silos directly and also reduces review wait (more people able to review).
> Risk-tiered approval stays after the feedback loops: we only lighten human review once data shows the automated layers catch what they should.

---

## Phase 0 — Discover and baseline

**Goal:** know what we have today, and measure where we start.

- [ ] 0.1 🟢 List all repos: language, size, activity, criticality, who knows it best (be honest — this is the silo map) 🔎 [T1]
- [ ] 0.1b 🟢 Write the team profile (`teams/README.md`): structure, pain points, constraints
- [ ] 0.2 🟡 Tooling inventory per repo: which org it's in, CI, existing branch rules, CODEOWNERS, security and dependency scanning, merge strategy (squash or merge) 🔎 [T1, T4]
- [ ] 0.3 🟡 Release inventory per repo: how it deploys, flags, how to roll back, how long it takes 🔎 [T6]
- [ ] 0.4 🟢 Run `tools/review-metrics` for the last 3–6 months (`03-templates.md` §8): 🔎 [M1, M2]
  - Time to first human review (median, p95)
  - Time from open to merge
  - Approvals and review comments per PR
  - **Who reviewed** and **who authored** (distribution across people)
  - **Who authored/reviewed per repo** (knowledge spread)
  - Follow-up fixes and reverts
  - Branch updates per PR (merges from base, rebases)
  - CI re-run rate (flaky-test signal)
- [ ] 0.5 🟢 15-minute team conversation: where does waiting hurt, which areas does only one or two people understand?
- [ ] 0.6 🟢 Choose the **pilot repo** 🔎 [P1]
- [ ] 0.7 🟢 Name a process owner and a co-owner from a different part of the team 🔎 [P2]

**Exit criteria:**
- Baseline numbers for the pilot repo exist (rough is fine)
- Silo map exists (repo/area → people who know it)
- Pilot repo agreed

---

## Phase 1 — Foundations + review-wait quick wins

**Goal:** basics later layers depend on, plus immediate relief for review wait. Mostly configuration, plus a few small workflows.

Foundations:
- [ ] 1.1 🟢 PR template in the pilot repo (**Intent / Follow-up to / Manual testing / Risk / How to revert**). It only asks what CI can't show (`03-templates.md` §1)
- [ ] 1.2 🟢 Required status checks on the pilot repo's main branch (lint, types, tests, build)
- [ ] 1.3 🟢 Branch rules: no direct pushes, required checks, 1 approval, **require conversation resolution**, **dismiss stale approvals on new commits**. **Don't** require the branch to be up to date before merging (`00-plan.md` §4.3). If a rule can't be enforced by the platform, enforce it with a workflow check instead 🔎 [T2, T19]
- [ ] 1.4 🟢 `CODEOWNERS` on pilot repo: sensitive paths and `.github/workflows/**` → their owners; everything else → whole team 🔎 [R1]
- [ ] 1.5 🟡 Rollback runbook for the pilot repo, tested once 🔎 [T6]
- [ ] 1.6 🟡 Start `CHANGELOG-review-process.md`
- [ ] 1.6b 🟢 **Changelog fragments** instead of editing `CHANGELOG.md` directly: one file per PR, assembled at release, with a CI check that the fragment exists (`00-plan.md` §4.3) 🔎 [T20]
- [ ] 1.6c 🟢 **Flaky-test policy:** re-run detection opens a `flaky-test` issue, flaky tests are fixed with priority or quarantined (with human approval), and open flaky issues are reviewed in the monthly check 🔎 [T21]
- [ ] 1.6d 🟢 **Branch update norm:** update only on conflict or stale CI. Merge base into the branch rather than rebasing once human review has started. Agents follow the same rule (`00-plan.md` §4.3)

Review-wait quick wins:
- [ ] 1.7 🟢 Create a GitHub team for the devs; enable **auto-assignment** (load balance) so reviews spread across everyone, not just the experts 🔎 [T14]
- [ ] 1.8 🟢 Agree a **first-response target** (e.g. one working day) and **one daily review slot per person**, at a time that fits their own timezone 🔎 [P8]
- [ ] 1.9 🟢 Soft PR size limit (e.g. ~400 changed lines) with "split or explain" as the escape hatch 🔎 [P3]
- [ ] 1.10 🟢 Draft-PR norm for direction questions ("should we do this?") 🔎 [P7]
- [ ] 1.11 🟡 Scheduled workflow that posts PRs waiting longer than the target to the team channel 🔎 [T15]
- [ ] 1.12 🟢 *(Optional, time-boxed to 1–2 months)* "Reviews before your own PRs" gate: a new PR's reviewers are only requested once the author has no review requests older than one working day (`00-plan.md` §4.1) 🔎 [T18]

**Exit criteria:**
- Template + rules active on pilot repo for 2 weeks
- Time to first human review measurably lower than baseline
- Review load starting to spread across more people

---

## Phase 2 — Per-project documentation

**Goal:** project knowledge the reviewer and humans need, written down and short. *Parallel with 3 and 4.*

- [ ] 2.1 🟢 Use `AGENTS.md` (vendor-neutral; decided). Add a one-line pointer file only for tools that require their own name
- [ ] 2.2 🟢 Pilot repo `AGENTS.md` — **drafted by someone who is not the main owner, reviewed by the owner**. Writing it is itself a learning exercise
- [ ] 2.3 🟢 Pilot repo `REVIEW.md` — owner corrects a draft extracted from their past review comments (`00-plan.md` §5.5); secondary owner reviews it
- [ ] 2.4 🟢 Mine past human review comments; move recurring mechanical ones into lint rules, judgement ones into `REVIEW.md`
- [ ] 2.4b 🟡 If several repos share a stack: extract common rules into `shared/<family>-guidelines.md` and vendor a copy into each repo with the sync tool (no cross-org references) 🔎 [D8]
- [ ] 2.5 🟡 Add the doc maintenance rules (`00-plan.md` §5.3) to each repo's contributing guide (vendored with the other shared files)
- [ ] 2.6 🟡 Agree length limits and a quarterly prune slot 🔎 [D2, D3]

**Exit criteria:**
- Both files exist on pilot repo and are reviewed by both owners
- Recurring human comments have a home

---

## Phase 3 — Verification before the PR

**Goal:** routine mistakes caught before anyone else looks. *Parallel with 2 and 4.*

- [ ] 3.1 🟢 One documented check command that runs format, vet, lint and tests (in cass-operator: `make test`), listed in `AGENTS.md`. No separate pre-commit hooks are needed when this exists; add them only to repos without such a command 🔎 [T7]
- [ ] 3.3 🟢 Local review step: fresh-context review of the diff + local security review 🔎 [T8]
- [ ] 3.4 🟢 Automated test results are read from CI, not copied into the PR. The description only covers **manual testing** that CI can't show
- [ ] 3.5 🟡 UI repos: define visual verification 🔎 [T9]

**Exit criteria:**
- Authors and agents run the check command before opening a PR (CI failures on first push trend down)
- Local review step used by the whole pilot team

---

## Phase 4 — Knowledge sharing and ownership model

**Goal:** knowledge spreads beyond the current experts; more people carry real review load. *Parallel with 2 and 3.* See `00-plan.md` §7.1; the team profile says how it maps onto your team.

- [ ] 4.1 🟢 Assign an **area pair** (primary + secondary owner) for each repo/area; encode in CODEOWNERS 🔎 [K4]
- [ ] 4.2 🟢 Route Low/Medium PRs from the experts to less familiar reviewers; every such review leaves a summary or a question
- [ ] 4.3 🟢 "Explain it back before your first change" norm for anyone entering an unfamiliar area (an AI explainer can help, but the explaining back is done by the person); whatever the developer docs were missing gets written into `docs/developer/` by hand 🔎 [K1]
- [ ] 4.4 🟡 Pair on High-risk changes
- [ ] 4.5 🟡 Define the **graduated ownership** criteria (non-owner → trusted for an area) 🔎 [R3]
- [ ] 4.6 🟡 Fortnightly 20-minute walkthrough in an existing meeting (rotating presenter, everyone takes a turn) 🔎 [K2]
- [ ] 4.7 🟡 Monthly knowledge-spread report from the GitHub API (who authored/reviewed per repo) 🔎 [K3, T12]

**Exit criteria (after ~6 weeks):**
- Every pilot area has at least the team's target number of people who authored or reviewed there recently
- Review and authorship distribution clearly more even than baseline

---

## Phase 5 — Automated review on every PR (pilot)

**Goal:** an automated reviewer comments on every pilot PR; findings are high-signal and must be answered.

- [ ] 5.1 🟢 Integrate automated review as a GitHub Actions workflow on `pull_request` 🔎 [T3]
- [ ] 5.1b 🟢 Public repos: decide how fork PRs get reviewed without exposing secrets (`00-plan.md` §8) 🔎 [T16]
- [ ] 5.2 🟢 Inline comments tagged **Important / Nit / Pre-existing**
- [ ] 5.3 🟡 Separate security review pass 🔎 [T4]
- [ ] 5.4 🟢 Drop style-only findings (lint's job)
- [ ] 5.5 🟢 Every finding fixed or resolved with a reason; conversation resolution blocks merge (from 1.3)
- [ ] 5.6 🟢 Resolution category captured (`[wrong]` / `[intentional]` / `[not-worth-it]` prefixes, see templates) 🔎 [M3]
- [ ] 5.7 🟡 Convergence guardrails: round cap 3, diff-only re-review, no repeats, flag the bug class 🔎 [T3]
- [ ] 5.8 🟢 **Review-state workflow and labels** (`review/bots` → `review/author` → `review/human`, plus the `review/bots-done` commit status). Human reviewers are requested only at `review/human` (`00-plan.md` §4.2) 🔎 [T14, T22]
- [ ] 5.8b 🟢 **Detector** for `needs-human`: changed or removed test assertions, new `//nolint`, sensitive paths 🔎 [T17]
- [ ] 5.9 🟢 Two weeks non-blocking, then blocking
- [ ] 5.10 🟢 Weekly `REVIEW.md` tuning during pilot

**Exit criteria (~4–6 weeks):**
- `fixed` share above agreed threshold; mean rounds-to-pass below threshold 🔎 [M4]
- Team says it's net helpful
- No rise in follow-up fixes

---

## Phase 6 — Feedback loops and eval set

**Goal:** the process improves itself, measurably.

- [ ] 6.1 🟡 Metrics dashboard (scheduled Action → markdown report or simple page) 🔎 [M1, T11]
- [ ] 6.2 🟡 Eval set from past follow-up fixes, reverts and bugs (the PR named in `Follow-up to:` is the eval case), each confirmed by someone who knows the code — **non-owners do the digging, an owner confirms** (good way to learn failure modes) 🔎 [M5]
- [ ] 6.3 ⚪ Known-clean PR set for false alarms
- [ ] 6.4 ⚪ Reviewer config changes scored against both sets before rollout 🔎 [M6]
- [ ] 6.5 🟢 Incident/bug review asks "which layer should have caught this?"
- [ ] 6.6 🟢 Monthly 30-minute retro **inside the existing team retro**, not a new meeting
- [ ] 6.7 ⚪ Periodic review of existing code 🔎 [T10]

**Exit criteria:**
- Eval set used on every reviewer change
- Two retros held, decisions logged

---

## Phase 7 — Risk-tiered human approval

**Goal:** human approvals scale with risk. **Start strict, loosen with evidence.**

- [ ] 7.1 🟡 Write tiers for the pilot repo (template in `03-templates.md`) 🔎 [R1]
- [ ] 7.2 🟢 Always-human paths via CODEOWNERS (sensitive + workflows)
- [ ] 7.3 🟡 Trusted list **per area**, using graduated ownership from 4.5 🔎 [R3]
- [ ] 7.4 ⚪ Tier assignment: start with **path rules + labels** in an Action; classifier later if needed 🔎 [R2]
- [ ] 7.5 🟡 Bot approval never counts as human. Enforce it in the review-state workflow if branch rules can't express it 🔎 [T2]
- [ ] 7.6 🟢 Agent-written PRs: directing person is accountable and counts as a reviewer
- [ ] 7.7 🟡 Low starts at 1 human; move to 0 only with ≥1 month of supporting data 🔎 [R4]
- [ ] 7.8 ⚪ Circuit breaker defined 🔎 [R5]
- [ ] 7.9 🟡 **Dependency/CVE update PRs** by a bot or agent following the rules in `00-plan.md` §4.4: latest patch, re-scan, batch, dependencies only. Low ones need one human for the first month, then auto-merge if follow-up fixes stay at zero 🔎 [T23]

**Exit criteria:**
- Tiers enforced on pilot; wait time down for Low/Medium; no guardrail regression for ≥1 month

---

## Phase 8 — Merge queue, flags, rollback

**Goal:** mistakes are cheap to catch and undo — this is what makes lighter approval safe.

- [ ] 8.1 ⚪ Merge queue, only when churn makes "passed separately, broke together" a real problem 🔎 [T5]
- [ ] 8.2 ⚪ Feature flags for new behaviour where the product allows 🔎 [T6]
- [ ] 8.3 ⚪ Internal use / staging before release 🔎 [P4]
- [ ] 8.4 🟡 Rollback time measured and within target 🔎 [M7]
- [ ] 8.5 🟡 Who watches after merge: author first, then a rotating "on-duty" dev (everyone in the rotation) 🔎 [P6]

**Exit criteria:**
- Rollback (or fix-forward) within target

---

## Phase 9 — Widen and sustain (ongoing)

- [ ] 9.1 🟡 Short adoption guide from pilot learnings
- [ ] 9.2 🟡 Roll to remaining repos in order of pain (most waiting / most siloed first). Shared pieces are vendored into each repo with the sync tool; each repo adds its own `AGENTS.md`, `REVIEW.md`, CODEOWNERS
- [ ] 9.3 🟡 Quarterly: doc prune (repo and shared), risk-tier review, graduated-ownership review, knowledge-spread check
- [ ] 9.4 🟡 Revisit out-of-scope items (tooling/vendor choices) once gaps are clear
