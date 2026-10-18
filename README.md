# Review guidelines

A plan for improving code review in a team where more and more code is written by AI agents, inspired by ["How the Claude Code team does code review"](https://claude.dev/blog/solving-code-review-at-anthropic/). It has two goals: less waiting for reviews, and fewer knowledge silos, without more breakage.

## Layout

| Path | What it is |
|---|---|
| `00-plan.md` | The plan: layers, where humans stay in the loop, review-state labels, CI friction, dependency updates, per-project docs, feedback loops, knowledge sharing, OSS. Team-neutral |
| `01-steps.md` | Phased step list. Every step is tagged 🟢 Core / 🟡 Later / ⚪ At scale |
| `02-open-questions.md` | What still needs investigating or deciding, plus decisions made so far |
| `03-templates.md` | PR template, `AGENTS.md` / `REVIEW.md` templates, risk tiers, CODEOWNERS, workflow sketches |
| `teams/` | Team profiles: structure, pain points, team-specific practices |
| `shared/` | Guidelines shared by a family of repos (Go controller-runtime operators). Vendored into each repo |
| `pilot/cass-operator/` | Pilot: baseline data, draft `AGENTS.md`, `REVIEW.md`, lint proposal |
| `tools/review-metrics/` | Go CLI for the baseline and the ongoing metrics |

Start with `00-plan.md` §0.1 ("Start small").
