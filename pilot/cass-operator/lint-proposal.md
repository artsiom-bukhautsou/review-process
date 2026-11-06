# Lint proposal — cass-operator

**Goal:** every mechanical rule lives in `.golangci.yml`, and the code complies with it. Nothing the linter can enforce is written in [`REVIEW.md`](REVIEW.md) or the shared guidelines. Reviewers, human or automated, only spend time on judgement.

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

Tighten `nolintlint`, which is already enabled: set `require-specific: true` (it must name the linter) and `require-explanation: true` (it must give a reason). New suppressions also need human approval ([shared guidelines §8](../../shared/go-operator-guidelines.md#8-change-hygiene)).

Already enabled and worth keeping: `nonamedreturns`, `ginkgolinter` (focus and spec pollution), `nolintlint`, `unparam`, `whitespace`.

Stale generated files are already caught in CI: `operatorBuildAndDeploy.yml` fails when `make manifests generate` leaves a diff.

## Additions from the reconcile audit

From [review-process#2](https://github.com/burmanm/review-process/issues/2): a `go func` literal in reconcile code that refers to `rc.Ctx` or the `ctx` parameter of `Reconcile` is an error. The reconcile context is cancelled when `Reconcile` returns, so the goroutine's work stops silently. No existing linter covers this, so it runs through `gocritic`'s `ruleguard` checker, or a small `go/analysis` pass if ruleguard can't express it. It ships with the fix for k8ssandra/cass-operator#1006.

## Rollout

1. **One PR updates `.golangci.yml`** and fixes the codebase so `make lint` passes. If the diff gets large, split it per linter: one PR enables `revive` and fixes, the next enables `govet` and fixes, and so on.
2. **Mechanical doesn't mean behaviour-free.** Each fix PR states that its rewrites are equivalent, and the reviewer checks every rewritten conditional on its own. Fix PRs never ride along with dependency bumps ([plan §4.4](../../00-plan.md#44-dependency-and-cve-updates)).
3. **These fix PRs are good first work for someone learning the codebase,** with a reviewer who knows the code they touch.
4. **After that, CI fails on any violation.** No new-issues-only mode is needed, because the backlog is fixed up front.
5. **When a new mechanical rule comes up in review**, add the linter rule and fix the code in the same PR, instead of writing the rule down.

## Stays judgement (cannot be linted well)

These live in [`REVIEW.md`](REVIEW.md) and the shared guidelines: reconcile blocking, transient vs terminal errors, per-reconcile cost, comments that restate code, unnecessary pointers or packages, test meaningfulness, changed test assertions, whether a goroutine outlives `Reconcile`, which code writes a condition `False`, and how the loop leaves a new blocking state.
