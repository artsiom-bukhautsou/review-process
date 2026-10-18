# Lint proposal — cass-operator

**Goal:** every mechanical rule lives in `.golangci.yml`, and the code complies with it. Nothing the linter can enforce is written in `REVIEW.md` or the shared guidelines. Reviewers, human or automated, only spend time on judgement.

## Quick wins already in the config

The current `.golangci.yml` has a full **`revive` rule set and `govet` settings configured, but both linters are commented out** in `linters.enable`. Turning them on enforces rules that come up in review today:

| Linter | Enforces |
|---|---|
| `revive` (existing rule set) | `if err := …; err != nil` flow (`indent-error-flow`), no `else` after `return` (`superfluous-else`), `:=` over `var` (`var-declaration`), godoc form on exported identifiers (`exported`), `error-strings`, `receiver-naming`, `context-as-argument` |
| `govet` (existing settings, `fieldalignment` off) | Shadowing, printf misuse, copylocks |

## Additions that map to recurring review comments

| Linter | Rule | Review comment it replaces |
|---|---|---|
| `testifylint` | `require-error` and others | "use `require.NoError(err)`" in unit tests |
| `forbidigo` | Forbid `time.Sleep` and `fmt.Print*` in `tests/**` and envtest suites | "No sleeps, use Eventually"; "no fmt.Println, logging handles it" |
| `errorlint` | `errors.Is` / `errors.As`, `%w` wrapping | Error comparison mistakes |
| `gocritic` | Small curated set (`ifElseChain`, `elseif`, `unlambda`, `assignOp`) | Assorted readability nits |
| `usestdlibvars` | `http.MethodGet` and similar | Magic strings |

Tighten `nolintlint`, which is already enabled: set `require-specific: true` (it must name the linter) and `require-explanation: true` (it must give a reason). New suppressions also need human approval (shared guidelines §8).

Already enabled and worth keeping: `nonamedreturns`, `ginkgolinter` (focus and spec pollution), `nolintlint`, `unparam`, `whitespace`.

Candidate to look into: a check that fails when `make manifests generate` produces a diff, so stale generated files are caught by CI and not by a reviewer.

## Rollout

1. **One PR updates `.golangci.yml`** and fixes the codebase so `make lint` passes. If the diff gets large, split it per linter: one PR enables `revive` and fixes, the next enables `govet` and fixes, and so on. Each PR is mechanical and easy to review.
2. **These fix PRs are good first work for someone learning the codebase.** The fixes touch every package without changing behaviour.
3. **After that, CI fails on any violation.** No new-issues-only mode is needed, because the backlog is fixed up front.
4. **When a new mechanical rule comes up in review**, add the linter rule and fix the code in the same PR, instead of writing the rule down.

## Stays judgement (cannot be linted well)

These live in `REVIEW.md` and the shared guidelines: reconcile blocking, transient vs terminal errors, per-reconcile cost, comments that restate code, unnecessary pointers or packages, test meaningfulness, and changed test assertions.
