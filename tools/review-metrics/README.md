# review-metrics

Code review metrics for merged pull requests in one or more GitHub repositories.
Uses the same libraries as the `gh` CLI
(`go-gh` for auth and REST, `cobra`, `lipgloss`), and only the REST API.

## Install

```sh
go install github.com/example/review-metrics@latest   # TODO: real module path
# or, from this directory:
go mod tidy && go build -o review-metrics .
```

Authentication follows `gh`: run `gh auth login`, or set `GH_TOKEN`. `GH_HOST`
selects a GitHub Enterprise host.

## Usage

```sh
review-metrics k8ssandra/cass-operator                       # last 6 months
review-metrics k8ssandra/cass-operator --since 2024-10-01
review-metrics org/op-a org/op-b --group alice,bob           # per repo + combined
review-metrics org/op-a --format markdown > metrics.md       # commit / paste in an issue
review-metrics org/op-a --ci --since 2026-07-01 --format json | jq '.reports[0].ci'
```

| Flag | Default | Meaning |
| --- | --- | --- |
| `--since` | 6 months ago | Only PRs created on or after this date (`YYYY-MM-DD`) |
| `--limit` | 500 | Max merged PRs per repository |
| `--group` | | Logins to compare against everyone else (e.g. current experts) |
| `--area-depth` | 2 | Directory depth that defines an area |
| `--ready-label` | `review/human` | Label that marks a PR ready for human review |
| `--ci` | off | Also fetch Actions runs for each PR head SHA (re-run rate) |
| `--format` | `text` | `text` (styled tables on a terminal, plain text otherwise), `json`, `markdown` |
| `--top` | 6 | Entries in ranked lists |
| `--cache-dir` | `$XDG_CACHE_HOME/review-metrics` | Response cache |
| `--no-cache` | off | Don't read or write the cache |
| `--workers` | 6 | Concurrent API requests |

Per-PR responses (reviews, files, commits, timeline, runs) are cached on disk by
URL, since merged PRs don't change; the PR list itself is always fetched fresh.
Use `--no-cache` if you need to pick up CI runs re-run after merge.

## Metric definitions

Only merged PRs created since `--since` are counted. A **human review** is a
submitted review by an account that is not a bot (`type == "Bot"` or login ending
in `[bot]`) and is not the PR author. Percentiles use nearest rank: the value at
index `round(p/100 * (n-1))` of the sorted list.

1. **Merged PRs**: count.
2. **Time to first human review**: PR creation to first human review (median, p90),
   plus the number of PRs merged without any human review.
3. **Open to merge**: PR creation to merge (median, p90).
4. **Review depth**: human reviews per PR (median, mean) and distinct human
   approvers per PR (0 / 1 / 2+).
5. **Distribution**: reviews per reviewer (and number of PRs they reviewed) and PRs
   per author, with shares. `--group` reports the group's share of reviews and of
   authored PRs.
6. **First-review wait by author**: median wait for each top author.
7. **PR size**: additions + deletions over the PR's files (median, p90).
8. **Knowledge spread**: an area is the first `--area-depth` directories of a file
   path (files at the root are `(root)`; shallower paths keep what they have). For
   each area: PRs touching it, distinct people (non-bot author plus human
   reviewers) and the top 3 people. Sorted fewest people first.
9. **Branch-update churn**: per PR, merge commits in the PR (commits with 2+
   parents, usually merging the base branch in) plus `head_ref_force_pushed`
   timeline events (rebases). Reports the share of PRs with at least one update,
   mean updates per PR, and merge-updates vs force-pushes.
10. **Follow-up fixes**: PRs whose body has a line like `Follow-up to: #123`
    (also `Followup to:`, `Regression from:`, `Fixes problem introduced in:`,
    `Fixes problem from:`; case-insensitive; one or more `#N`). Reports count, share
    and "#X was followed up by #Y" links. **This depends on the PR template having a
    `Follow-up to: #N` field**; until it does, expect 0. Revert PRs (title starts
    with "Revert") are still counted.
11. **Ready-label timing**: from the first time `--ready-label` was applied to the
    first human review at or after it (median, p90). Shows
    `n/a (label not in use yet)` when no PR has the label.
12. **CI re-runs** (`--ci`): for each PR's final head SHA, `pull_request` workflow
    runs; the share of SHAs with any run where `run_attempt > 1`, and workflows
    ranked by re-run count (a flaky-test signal).

## Layout

```
main.go                  entry point
cmd/                     cobra command, go-gh client, terminal detection
internal/app             orchestration (fetch -> compute -> render)
internal/github          REST fetch layer: pagination, bounded concurrency, cache
internal/cache           on-disk response cache
internal/model           GitHub REST payload types
internal/metrics         pure metric computations (unit tested)
internal/report          document model, plain text / Markdown / JSON renderers
internal/report/styled   lipgloss renderer for terminals
internal/progress        stderr spinner
```

## Development

```sh
go vet ./...
go test ./...
```


## Known follow-ups

- **Build not yet verified with real dependencies.** The code was written in an environment that couldn't download Go modules. Run `go mod tidy && go vet ./... && go test ./...` first, and fix anything that comes up.
- **GraphQL is allowed.** Fetching is REST-only today (about 400 calls cold for 98 PRs). One GraphQL query per page of PRs could replace most of those calls, and is needed anyway for review-thread resolution (finding outcomes).
- **Module path** in `go.mod` is a placeholder (`github.com/example/review-metrics`).
