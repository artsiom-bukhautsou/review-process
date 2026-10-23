# Review guidelines

A plan for improving code review in a team where more and more code is written by AI agents, inspired by ["How the Claude Code team does code review"](https://claude.dev/blog/solving-code-review-at-anthropic/). It has two goals: less waiting for reviews, and fewer knowledge silos, without more breakage.

## Layout

| Path | What it is |
|---|---|
| [`00-plan.md`](00-plan.md) | The plan: layers, where humans stay in the loop, review-state labels, CI friction, dependency updates, per-project docs, feedback loops, knowledge sharing, OSS. Team-neutral |
| [`01-steps.md`](01-steps.md) | Phased step list. Every step is tagged 🟢 Core / 🟡 Later / ⚪ At scale |
| [`02-team-adoption.md`](02-team-adoption.md) | Questionnaire for each adopting team; answers and adoption progress live in its team profile |
| [`03-process-tracker.md`](03-process-tracker.md) | Shared process investigations, implementation progress, and decision history |
| [`04-templates.md`](04-templates.md) | PR template, `AGENTS.md` / `REVIEW.md` templates, risk tiers, CODEOWNERS, workflow sketches |
| [`teams/`](teams/) | Team profiles: structure, pain points, team-specific practices |
| [`shared/`](shared/) | Guidelines shared by a family of repos (Go controller-runtime operators). Vendored into each repo |
| [`pilot/cass-operator/`](pilot/cass-operator/) | Pilot: baseline data, draft [`AGENTS.md`](pilot/cass-operator/AGENTS.md), [`REVIEW.md`](pilot/cass-operator/REVIEW.md), [lint proposal](pilot/cass-operator/lint-proposal.md) |
| [`tools/review-metrics/`](tools/review-metrics/) | Go CLI for the baseline and the ongoing metrics |

Start with [`00-plan.md` §0.1](00-plan.md#01-start-small-what-to-do-first) ("Start small").
