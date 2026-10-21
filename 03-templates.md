# Starter Templates

**How to use:** adapt per repo. These are things we can create today, without new tooling.

---

## 1. Pull request template

Save as `.github/pull_request_template.md` **in each repo**. Our repos are spread over several orgs, so there's no shared default. The copy is kept in sync from the source-of-truth repo (see [`00-plan.md` §5.4](00-plan.md#54-shared-guidance-across-a-family-of-repos)).

The template only asks for what GitHub can't show. Automated test results, which test files changed, and whether CI passed are all visible on the PR, and failing CI blocks the merge anyway, so none of that is repeated in the description.

```markdown
## Intent
<!-- What problem does this solve, and why this approach? -->
Fixes #

## Follow-up to
<!-- Does this fix a problem introduced by an earlier PR? List it: "Follow-up to: #123". Leave empty otherwise.
     Used to measure how often changes need fixing after merge, and to build the reviewer's eval set. -->
Follow-up to:

## Manual testing
<!-- Only what CI can't show: what you tried by hand, on which environment, and what you saw. "None" is a valid answer. -->

## Risk
<!-- Low / Medium / High (see RISK-TIERS.md). What could break? What's the blast radius? -->
Tier:
What could go wrong:

## Intentional trade-offs
<!-- Known limitations, or cases this PR deliberately doesn't handle, and why. The automated reviewer doesn't flag what's stated here. "None" is a valid answer. -->

## How to revert
<!-- Plain revert? Needs a follow-up fix instead (data/CRD migrations)? Anything that can't be undone? -->

## Accountability
<!-- If an agent wrote this, who directed it? That person is accountable and counts as a reviewer. -->
Directed by:
```

Repos keep their own obligations on top, such as a CLA checkbox.

---

## 2. `AGENTS.md` (project context)

Vendor-neutral, so every reviewing agent reads the same file. If a tool only reads its own file name, add that file with a single line: `See AGENTS.md`.

Keep it short. Point to other docs rather than copying them.

```markdown
# <Project name>

## What this is
One paragraph: what the project does and who uses it.

## Stack
- Language / framework / key libraries
- Where it runs

## How to verify changes
- Test:        `<command>`
- Lint:        `<command>`
- Type check:  `<command>`
- Build:       `<command>`
- UI check:    `<how to run and screenshot>`  (if applicable)

## Conventions
- <naming, structure, error handling, logging, etc. — only the ones that matter>

## Goals
- <what we optimise for>

## Non-goals
- <what we deliberately don't do — prevents "helpful" scope creep>

## User documentation
- User-facing here means: <CRD fields, annotations, flags, behaviour changes, metrics …>
- Where docs go: <`docs/` path | docs repo | docs team tracker>
- Format: <full page | docs note in the changelog fragment (`docs:` field)>
- Picked up by: <docs team contact>

## Further reading
- Architecture: `docs/architecture.md`
- Decisions: `docs/adr/`

<!-- If a change makes any line here out of date, update it in the same PR. -->
```

---

## 3. `REVIEW.md` (owner-written review guidance)

```markdown
# Review guidance for <repo / area>

Owners: <team or people>

<!-- Judgement only. If a rule here can be checked by a script, write the check and delete the line. -->

## What failure looks like here
<!-- One or two lines. e.g. "Config: it misroutes traffic." "Docs: it lies." "Billing: it charges the wrong amount." -->

## Focus on
- <the 3–6 things that actually cause bugs in this code>

## Don't flag
- Style issues the linter already enforces
- Anything a CI check already enforces (changelog fragment, generated files, linked issue)
- <patterns that look wrong but are intentional, with a short why>
- Trade-offs the PR description states, unless they lead to a failure under "What failure looks like here". Then ask, don't report a finding

## Needs a human (a marker, not a finding)
New `//nolint` directives and changed or removed assertions in existing tests need human approval even when correct.
List them on one line at the end of the summary: `Needs human: <what> in <file:line>. Explained in the description: yes/no.`
Report a finding only when the description doesn't explain the change.

## Skip these paths
- `generated/**`
- `vendor/**`
- `**/fixtures/**`

## Evidence
Every finding must cite a file and line, and explain the concrete failure. No guesses based on names alone.

## Limits
- At most <3> nits per review.
- Tag each finding: Important / Nit / Pre-existing.
- Pre-existing findings don't count toward the limits. List them in the summary only, each with a one-line proposed `tech-debt` issue.
- Don't repeat findings already posted on this PR.
- If a bug pattern appears more than once, flag the pattern and list all instances.

## Summary format
Start with a count: "N important, M nits, K pre-existing."
If there is nothing to fix, say so plainly.
```

---

## 4. Risk tiers (`RISK-TIERS.md`)

Start strict. Loosen only with evidence from the feedback loops.

```markdown
# Risk tiers — <repo>

The change sets the tier, not the author.

| Tier   | Examples                                                         | Human approvals (start) | Human approvals (target, with evidence) |
|--------|------------------------------------------------------------------|-------------------------|-----------------------------------------|
| Low    | Docs, copy, internal-only, additive behind a flag                | 1                       | 0 (trusted authors, all checks green)   |
| Medium | Ordinary product changes; sensitive-but-not-critical areas       | 1                       | 1                                       |
| High   | <critical paths listed below>                                    | 2                       | 2                                       |

## Always need a human (regardless of tier)
- CI/CD workflow files
- Auth, billing/payments, installers/updaters, data migrations
- Changed or removed assertions in existing tests (they can mask a newly introduced bug)
- New `//nolint` directives
- Dependency PRs with a major version bump, a new or removed dependency, or a `replace` directive
- <add>

## Rules
- A bot approval never counts as a required human approval.
- Approval applies to a specific commit; new commits dismiss it.
- On agent-written PRs, the directing person is one of the required reviewers.
- Trusted status is **per area** (graduated ownership): a person becomes trusted for an area after the agreed number of authored + reviewed PRs there without reverts.
- Authors not trusted for the area always need a second person — ideally that area's secondary owner.

## Circuit breaker
If the revert or incident rate exceeds <X> over <period>, all tiers move up one level until reviewed.
```

---

## 5. Finding resolution convention

Lightweight way to capture Loop A data until tooling does it automatically.

When resolving an automated finding without a code change, start the reply with one of:

- `[wrong]` — reviewer is incorrect (say why in one line)
- `[intentional]` — behaviour is deliberate (say why)
- `[not-worth-it]` — real but not worth fixing now (link an issue if it matters later)

Findings resolved by a code change are counted as `fixed`.

---

## 6. Incident review addition

Add to the existing incident/postmortem template:

```markdown
## Which layer should have caught this?
- [ ] Verification before PR (tests/checks the author could have run)
- [ ] Required CI
- [ ] Automated code review
- [ ] Automated security review
- [ ] Human approval
- [ ] Merge queue
- [ ] Flags / staged release / monitoring

What was missing from that layer? (context in AGENTS.md/REVIEW.md, a lint rule, a test, a check)

Action: <the specific fix to the layer>
Eval case added: [ ] yes  (link)
```

---

## 7. Process decision log (`CHANGELOG-review-process.md`)

```markdown
# Review process changelog

## YYYY-MM-DD — <decision title>
- **Change:** what we changed
- **Why:** evidence (metric, incident, retro feedback)
- **Scope:** repos / paths affected
- **Revisit by:** date
```

---

## 8. Metrics tool

[`tools/review-metrics/`](tools/review-metrics/) is a Go CLI built on the libraries `gh` uses: go-gh for auth and API, cobra, and lipgloss. It produces the Phase 0 baseline and the ongoing numbers:
- Review timing: time to first human review, open to merge, reviews and approvals per PR
- Distribution: review and authorship distribution, an optional `--group`, and first-review wait per author
- Size and knowledge spread: PR size, and knowledge spread per area
- Friction: branch-update churn (merges from base and force-pushes)
- Breakage: follow-up fixes (`Follow-up to: #N`) and reverts
- Review-state timing: time from the `review/human` label to the first human review
- With `--ci`: CI re-run rate per workflow, which flags flaky tests

```bash
cd tools/review-metrics && go mod tidy && go build -o review-metrics .
./review-metrics k8ssandra/cass-operator --since 2024-10-01 --group alice,bob
./review-metrics org-a/repo1 org-b/repo2 --format markdown > report.md   # several repos, across orgs
./review-metrics k8ssandra/cass-operator --since 2026-07-01 --ci          # flaky-test signal
```

Output formats: styled tables on a terminal, plain text when piped, `--format markdown`, or `--format json`. Responses are cached on disk; `--no-cache` refreshes them. The README in the tool's folder has the full metric definitions.

---

## 9. CODEOWNERS with area pairs

```text
# Default: whole team (auto-assignment spreads reviews across everyone)
*                         @our-org/devs

# Area pairs: primary + secondary owner
/src/api/                 @owner-a @second-a
/src/ui/                  @owner-b @second-b
/src/billing/             @owner-c @second-c

# Always an owner: sensitive paths and CI
/src/auth/                @owner-a @owner-c
/.github/workflows/       @owner-a @owner-b
/CODEOWNERS               @owner-a @owner-b
```

Later lines win in CODEOWNERS, so put specific paths after the default.

---

## 10. Waiting-PR nudge (GitHub Action)

Keeps one issue titled "PRs waiting for review" up to date, twice a day. Swap the issue update for a chat webhook once T15 is decided.

```yaml
# .github/workflows/waiting-prs.yml
name: Waiting PRs
on:
  schedule:
    - cron: "0 6 * * 1-5"   # UTC; once a day, before most people's review slot
  workflow_dispatch:
permissions:
  issues: write
  pull-requests: read
jobs:
  nudge:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/github-script@v7
        with:
          script: |
            const TARGET_HOURS = 4;
            const { owner, repo } = context.repo;
            const prs = await github.paginate(github.rest.pulls.list, { owner, repo, state: "open" });
            const now = Date.now();
            const waiting = [];
            for (const pr of prs) {
              if (pr.draft) continue;
              const reviews = await github.rest.pulls.listReviews({ owner, repo, pull_number: pr.number });
              const human = reviews.data.filter(r => r.user.type !== "Bot" && r.user.login !== pr.user.login);
              if (human.length) continue;
              const h = (now - new Date(pr.created_at)) / 36e5;
              if (h >= TARGET_HOURS) {
                const who = pr.requested_reviewers.map(u => "@" + u.login).join(" ") || "(no reviewer)";
                waiting.push(`- [ ] #${pr.number} ${pr.title} — ${Math.round(h)}h — ${who}`);
              }
            }
            const title = "PRs waiting for review";
            const body = waiting.length ? waiting.join("\n") : "Nothing waiting. 🎉";
            const { data: issues } = await github.rest.issues.listForRepo({ owner, repo, state: "open", per_page: 100 });
            const existing = issues.find(i => i.title === title);
            if (existing) await github.rest.issues.update({ owner, repo, issue_number: existing.number, body });
            else await github.rest.issues.create({ owner, repo, title, body });
```

---

## 11. "Reviews before your own PRs" gate (GitHub Action, sketch)

This is a soft gate for the bootstrapping period ([`00-plan.md` §4.1](00-plan.md#41-cutting-review-wait)). It never fails a check. It only holds back requesting reviewers until the author's own review queue is clear. It is a sketch to prototype (T18), not production-ready.

```yaml
# .github/workflows/review-gate.yml
name: Review gate
on:
  pull_request:
    types: [opened, ready_for_review, reopened]
  schedule:
    - cron: "0 * * * 1-5"        # hourly re-check of held PRs
  workflow_dispatch:
permissions:
  pull-requests: write
  issues: write
jobs:
  gate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/github-script@v7
        with:
          script: |
            const { owner, repo } = context.repo;
            const LABEL = "waiting-on-author-reviews";
            const SKIP = "skip-review-gate";
            const TEAM = "devs";                       // GitHub team to request
            const MAX_AGE_H = 24;                      // "one working day", simplified
            const TEAM_ROLES = ["OWNER", "MEMBER", "COLLABORATOR"];

            // Approximation: a pending request = the user is still a requested
            // reviewer on an open PR last updated more than MAX_AGE_H ago.
            async function pendingReviews(user) {
              const q = `is:pr is:open review-requested:${user} org:${owner}`;
              const res = await github.rest.search.issuesAndPullRequests({ q, per_page: 50 });
              const cutoff = Date.now() - MAX_AGE_H * 36e5;
              return res.data.items.filter(i => new Date(i.updated_at) < cutoff);
            }

            async function evaluate(pr) {
              if (pr.draft || !TEAM_ROLES.includes(pr.author_association)) return;
              const labels = pr.labels.map(l => l.name);
              if (labels.includes(SKIP)) return;
              const pending = await pendingReviews(pr.user.login);
              if (pending.length) {
                if (!labels.includes(LABEL)) {
                  await github.rest.issues.addLabels({ owner, repo, issue_number: pr.number, labels: [LABEL] });
                  const list = pending.map(i => `- ${i.html_url}`).join("\n");
                  await github.rest.issues.createComment({ owner, repo, issue_number: pr.number,
                    body: `Reviewers will be requested once your pending reviews are done:\n${list}` });
                }
              } else if (labels.includes(LABEL) || context.eventName === "pull_request") {
                if (labels.includes(LABEL))
                  await github.rest.issues.removeLabel({ owner, repo, issue_number: pr.number, name: LABEL });
                await github.rest.pulls.requestReviewers({ owner, repo, pull_number: pr.number, team_reviewers: [TEAM] });
              }
            }

            if (context.eventName === "pull_request") {
              await evaluate(context.payload.pull_request);
            } else {
              const held = await github.paginate(github.rest.pulls.list, { owner, repo, state: "open" });
              for (const pr of held.filter(p => p.labels.some(l => l.name === LABEL))) await evaluate(pr);
            }
```

To work, this needs:
- Team auto-assignment turned **off** for the repo, because reviewers are requested by automation. Once the review-state workflow ([`00-plan.md` §4.2](00-plan.md#42-knowing-when-bots-are-done-review-state-labels)) exists, this gate becomes one more condition for `review/human` rather than a separate workflow. CODEOWNERS for sensitive paths still request their owners directly.
- A token that can request team reviewers. The default `GITHUB_TOKEN` may not be able to; this is part of T18.
- On public repos: `pull_request` from forks gets a read-only token. Since outside contributors are exempt anyway, the gate only matters for branches pushed to the repo itself.

---

## 12. Review-state labels and workflow (sketch)

Implements [`00-plan.md` §4.2](00-plan.md#42-knowing-when-bots-are-done-review-state-labels). Treat it as a starting point to prototype (T22), not production code.

### Label definitions

Create these in each repo. The sync tool can do it, since labels are per repo.

| Label | Colour | Description |
|---|---|---|
| `review/bots` | `#c5def5` | Bots are reviewing or fixing; humans wait |
| `review/author` | `#fbca04` | Bots done; CI red or Important findings open; author acts |
| `review/human` | `#0e8a16` | Bots done and clean; human review requested |
| `needs-human` | `#b60205` | Human approval required regardless of risk tier |
| `round-cap` | `#d93f0b` | 3 automated rounds reached; human reconsiders the approach |
| `deps` | `#0366d6` | Dependency update; [§4.4](00-plan.md#44-dependency-and-cve-updates) rules apply |
| `flaky-test` | `#e99695` | (issues) A test that failed then passed on the same commit |
| `skip-changelog` | `#ededed` | Docs-only or otherwise no changelog fragment needed |

### Workflow

Each automated reviewer, whatever its vendor, reports a check run named `ai-review/<name>` on the head SHA. A bot fix agent reports `ai-fix/<name>` while it runs. The state workflow doesn't care how many of each there are; it only lists the expected reviewers.

```yaml
# .github/workflows/review-state.yml
name: Review state
on:
  pull_request_target:            # write access on fork PRs; never checks out PR code
    types: [opened, reopened, ready_for_review, synchronize, converted_to_draft]
  check_run:
    types: [completed]
  pull_request_review:
    types: [submitted]
  schedule:
    - cron: "*/15 * * * 1-5"      # fallback for events Actions can't trigger on (e.g. thread resolution)
  workflow_dispatch:
permissions:
  pull-requests: write
  issues: write
  statuses: write
  checks: read
concurrency:
  group: review-state-${{ github.event.pull_request.number || github.event.check_run.head_sha || 'sweep' }}
  cancel-in-progress: false
jobs:
  evaluate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/github-script@v7
        with:
          script: |
            const { owner, repo } = context.repo;
            const EXPECTED_REVIEWERS = ["ai-review/primary"];   // add one entry per automated reviewer
            const REQUIRED_CI = ["Run unit tests"];             // names of required CI check runs
            const QUIET_MINUTES = 10;                           // no bot pushes for this long
            const ROUND_CAP = 3;
            const TEAM = "devs";
            const STATE = ["review/bots", "review/author", "review/human"];
            const BOT = u => u?.type === "Bot" || u?.login?.endsWith("[bot]");

            async function prsToEvaluate() {
              if (context.payload.pull_request) return [context.payload.pull_request];
              if (context.payload.check_run) return context.payload.check_run.pull_requests
                .map(p => p.number).map(async n => (await github.rest.pulls.get({ owner, repo, pull_number: n })).data);
              return github.paginate(github.rest.pulls.list, { owner, repo, state: "open" });
            }

            async function setState(pr, label, statusState, description) {
              const have = pr.labels.map(l => l.name);
              for (const l of STATE.filter(l => l !== label && have.includes(l)))
                await github.rest.issues.removeLabel({ owner, repo, issue_number: pr.number, name: l });
              if (!have.includes(label))
                await github.rest.issues.addLabels({ owner, repo, issue_number: pr.number, labels: [label] });
              await github.rest.repos.createCommitStatus({ owner, repo, sha: pr.head.sha,
                context: "review/bots-done", state: statusState, description });
            }

            async function unresolvedBotImportant(pr) {
              // Thread resolution is only exposed through GraphQL.
              const q = `query($o:String!,$r:String!,$n:Int!){ repository(owner:$o,name:$r){ pullRequest(number:$n){
                reviewThreads(first:100){ nodes{ isResolved comments(first:1){ nodes{ body author{ __typename } } } } } } } }`;
              const res = await github.graphql(q, { o: owner, r: repo, n: pr.number });
              return res.repository.pullRequest.reviewThreads.nodes.filter(t => !t.isResolved &&
                t.comments.nodes[0]?.author?.__typename === "Bot" && /\bImportant\b/.test(t.comments.nodes[0].body)).length;
            }

            async function evaluate(pr) {
              if (pr.state !== "open") return;
              if (pr.draft) return setState(pr, "review/bots", "pending", "Draft");
              const sha = pr.head.sha;
              const checks = await github.paginate(github.rest.checks.listForRef, { owner, repo, ref: sha });
              const byName = Object.fromEntries(checks.map(c => [c.name, c]));
              const done = n => byName[n]?.status === "completed";
              const ok = n => done(n) && ["success", "neutral", "skipped"].includes(byName[n].conclusion);

              const fixRunning = checks.some(c => c.name.startsWith("ai-fix/") && c.status !== "completed");
              const commits = await github.paginate(github.rest.pulls.listCommits, { owner, repo, pull_number: pr.number });
              const last = commits[commits.length - 1];
              const lastBotPush = BOT(last?.author) || BOT(last?.committer);
              const quiet = (Date.now() - new Date(last?.commit?.committer?.date ?? 0)) / 6e4 > QUIET_MINUTES;
              const rounds = commits.filter(c => BOT(c.author) || BOT(c.committer)).length;  // bot pushes ≈ rounds

              const allReviewed = EXPECTED_REVIEWERS.every(done);
              const ciDone = REQUIRED_CI.every(done);
              if (!allReviewed || !ciDone || fixRunning || (lastBotPush && !quiet))
                return setState(pr, "review/bots", "pending", "Automated review/fix in progress");

              if (rounds >= ROUND_CAP) {
                await github.rest.issues.addLabels({ owner, repo, issue_number: pr.number, labels: ["round-cap", "needs-human"] });
              } else {
                const ciGreen = REQUIRED_CI.every(ok);
                const open = await unresolvedBotImportant(pr);
                if (!ciGreen || open > 0)
                  return setState(pr, "review/author", "success", `Bots done: ${ciGreen ? "" : "CI red; "}${open} Important open`);
              }
              if (pr.labels.some(l => l.name === "waiting-on-author-reviews"))
                return setState(pr, "review/author", "success", "Author has pending reviews");

              await setState(pr, "review/human", "success", "Bots done; ready for human review");
              if (!pr.requested_teams?.length && !pr.requested_reviewers?.length)
                await github.rest.pulls.requestReviewers({ owner, repo, pull_number: pr.number, team_reviewers: [TEAM] });
            }

            for (const p of await prsToEvaluate()) await evaluate(await p);
```

Notes and known gaps to settle during the prototype (T22):
- **Round counting** here approximates rounds as "commits by a bot". A sturdier option is to keep a counter in a hidden comment (`<!-- review-rounds: N -->`) that is incremented on each `review/bots` → `review/author` transition.
- **Requesting a team** may need a token with more rights than the default `GITHUB_TOKEN`.
- **Fork PRs:** `check_run` events from forks carry no `pull_requests` list. The scheduled sweep covers those.
- The `review/bots-done` status can be made a required check once the workflow is trusted.
