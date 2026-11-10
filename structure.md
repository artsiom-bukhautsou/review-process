# Structure

The documents are **team-neutral**: they describe the process, the layers and the feedback loops that apply to any team. Things that depend on a particular team or repo live in separate files.

| File | Contains | Reuse by another team |
|---|---|---|
| [`00-plan.md`](00-plan.md) | The plan: layers, where humans stay in the loop, review-state labels, CI friction, dependency updates, per-project docs, feedback loops, knowledge sharing, OSS | Use as-is |
| [`01-steps.md`](01-steps.md) | Phased step list with exit criteria. Every step is tagged 🟢 Core / 🟡 Later / ⚪ At scale, so a small team knows where to start | Use as-is |
| [`02-open-questions.md`](02-open-questions.md) | What still needs investigating or deciding, plus decisions made so far | Use as-is |
| [`03-templates.md`](03-templates.md) | PR template, `AGENTS.md` / `REVIEW.md` templates, risk tiers, CODEOWNERS, workflow sketches | Use as-is |
| [`tooling.md`](tooling.md) | The tools and agents the plan needs, and which exist yet | Use as-is |
| [`teams/<team>.md`](teams/) | Team context, pain points, team structure, how review load and ownership are split in *that* team, team-specific decisions | Write your own from [`teams/README.md`](teams/README.md) |
| [`pilot/<repo>/`](pilot/) | Baseline data and drafts of `AGENTS.md`, `REVIEW.md` and lint changes for a pilot repo | Use as a worked example |
| [`shared/`](shared/) | Guidance that applies to a family of repos (e.g. all Go controller-runtime operators). Vendored into each repo | Use if your repos share the stack |
| [`tools/review-metrics/`](tools/review-metrics/) | Go CLI that produces the baseline and the ongoing metrics | Use as-is |

**Platform assumptions.** GitHub is treated as just a tool. The plan uses common features (branch rules, CODEOWNERS, labels, checks, Actions). If a feature is missing somewhere, that check moves into our own automation or becomes a review criterion. The plan doesn't depend on any particular product tier. Repos may also be spread over **several orgs** with a mix of OSS and private, so the plan assumes **nothing is shared across orgs** except ticket references. Shared material is copied into each repo ([plan §5.4](00-plan.md#54-shared-guidance-across-a-family-of-repos)).
