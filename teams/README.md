# Team profiles

The main plan ([`00-plan.md`](../00-plan.md) – [`04-templates.md`](../04-templates.md)) is team-neutral. Each team that adopts it writes a short profile here, `teams/<team-name>.md`, covering only what differs for them.

Use [`02-team-adoption.md`](../02-team-adoption.md) as the adoption questionnaire. Copy relevant questions into this profile's **Team-specific open questions** table, preserving their IDs; link answers already recorded in the context or decisions instead of duplicating them. Track status here and record agreed answers with a date in **Team decisions**. Skip questions for rollout tiers the team is not adopting yet.

Shared implementation questions belong in [`03-process-tracker.md`](../03-process-tracker.md). Link those dependencies from the profile when needed; their status is maintained by the process authors.

## Template

```markdown
# Team profile — <team name>

**Last updated:** YYYY-MM-DD   **Process owner:** <name>

## Context
| Item | Value | What it means for the plan |
|---|---|---|
| Team size and structure | | |
| Code hosting / CI | | |
| Repos (count, public/private, mono/multi) | | |
| Compliance / quality requirements | | |
| Pain #1 | | |
| Pain #2 | | |

## Diagnosis
<!-- One paragraph: what causes the pain points? Use the Phase 0 baseline numbers. -->

## How the ownership model (00-plan §7.1) maps onto this team
<!-- Who are primaries/secondaries? Who reviews whom? What does graduated ownership look like here? -->

## Team-specific adjustments to the steps
<!-- Only deviations from 01-steps.md -->

## Team decisions
| Date | Decision | Why |
|---|---|---|

## Team-specific open questions
| ID | Question | Status | Answer |
|---|---|---|---|
```
