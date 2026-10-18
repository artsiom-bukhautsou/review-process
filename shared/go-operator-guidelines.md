# Go operator guidelines (shared)

**Applies to:** all our Go controller-runtime operators. Each repo's `AGENTS.md` links here; each repo's `REVIEW.md` adds what is specific to that repo. \
**Source:** harvested from past cass-operator review comments (2021–2026), corrected by the owner.

**Scope:** judgement rules only. Anything a linter can enforce lives in `.golangci.yml`, and the codebase is fixed to comply with it (see [`pilot/cass-operator/lint-proposal.md`](../pilot/cass-operator/lint-proposal.md)). If a rule here becomes lintable, it moves to the linter config and is deleted from this file.

---

## 1. Reconcile loop

1. **Never block reconcile.** No long-running calls, waits or polling loops inside reconcile. Do the minimum, return, and requeue to check again.
2. **Reconcile runs very often.** In large clusters with many custom resources it runs constantly. Every extra API call, Secret read or config unmarshal is multiplied. Exit early when a feature isn't used, before doing any work for it.
3. **Know whether an error is transient or terminal.** If retrying with the same input cannot succeed, don't requeue in a loop. Use a terminal error (`reconcile.TerminalError`) or surface a status condition. Retrying with the same data that already failed only spins.
4. **Logging is not handling.** Logging an error and then returning an error result for something that will never succeed creates an endless reconcile.
5. **Use the sequence's own invariants before adding state.** If reaching the end of reconcile already means an earlier step succeeded, don't add a boolean to track it.
6. **Keep errors where they belong.** Return an error at the point it happens and let the retry handle it, rather than threading it into an unrelated process such as decommission.
7. **Respect controller lag.** The StatefulSet controller (or any other) may not have caught up yet. Don't act on objects whose status doesn't reflect the latest change. Show where in the sequence that situation is handled.
8. **State in goroutines doesn't persist.** Anything a goroutine computes is lost unless it's written to an object, for example a status or an annotation. Don't mutate outer variables from inside one and expect it to matter.
9. **Fatal startup errors exit.** In `main`, an error that should stop the process calls `os.Exit(1)`. It must not carry on and start reconcilers.

## 2. Kubernetes semantics

1. **Label values and annotation values have different rules.** Keys are similar, values are not. Use the Kubernetes validation helpers rather than hand-rolled regexes or blacklists.
2. **Object names derive from `metadata.name`.** Overrides used inside the application (for example a datacenter name override) must not leak into Kubernetes object names, or two resources can collide.
3. **Handle optional fields the API allows.** If Kubernetes allows ports without names, an empty template, and so on, the code must cope with them.
4. **Check that the platform actually supports what you rely on.** Examples: a StorageClass that supports modification, the StatefulSet controller accepting an in-place change, a CSI driver approving a change. Then reconcile until the change is confirmed, not just requested.
5. **Prefer OpenAPI validation (CRD markers) for static field rules.** Use webhooks for rules that need the old object, such as invalid transitions.

## 3. API design (CRDs)

1. **Backward compatibility first.** Existing resources must keep working after an operator upgrade. API changes are called out as such in the CHANGELOG.
2. **One way to configure a thing.** Don't add a second mechanism for something already configurable elsewhere.
3. **No redundant switches.** An optional pointer struct doesn't need an extra `enabled` flag.
4. **Reuse existing API groups and packages** unless there is a stated reason for a new one.
5. **Don't hardcode versions or image references in Go code.** Put them in configuration, so changing them doesn't need a new operator release.
6. **Regenerate generated code and manifests** with the repo's pinned tool versions (`make manifests generate`), and commit the result.

## 4. Go style

Effective Go is the baseline; mechanical style is the linter's job. Reviewers look for what the linter can't see:

1. **Don't use pointers without a reason**, such as optional values, mutation or large structs.
2. **Use the right tool for the job,** for example `strings.Builder` for building strings rather than repeated `[]byte` casts.
3. **Name constants.** Annotation keys, label prefixes and magic strings become `const`s.
4. **Reuse existing helpers** rather than reimplementing naming or label logic. Look for existing functions before writing new ones.
5. **No new packages or files for a handful of helpers used in one place.**
6. **Names should say what something does.** If a function doesn't do what its name says, rename it or change it.

## 5. Comments

1. **Don't comment what the code already says.** No comments restating a type or a field. That is what a typed language and the editor are for.
2. **Reasoning usually goes in the commit message,** where `git blame` leads to it.
3. **Exception: a short comment for code that would otherwise confuse a reader,** such as a non-obvious rule, invariant or external constraint. Keep it to the minimum that prevents someone from "fixing" it. This is the exception, not the default.
4. **Delete commented-out code.**

## 6. Tests

Three levels. Each catches different problems; none replaces another.

| Level | Where | Framework | Use it for |
|---|---|---|---|
| **Unit** | Next to the code (`pkg/...`) | Go testing + testify (`require` for checks that must stop the test) | Pure logic: builders, parsing, calculations |
| **envtest** (integration) | In the controller and webhook packages, since the suite starts a real API server | Ginkgo/Gomega | Controller and webhook behaviour against a real API server: object creation, status updates, requeues, validation |
| **End-to-end** | `tests/<scenario>/` | Ginkgo/Gomega on a kind cluster | Real Cassandra behaviour: scaling, decommission, upgrades, version-specific behaviour |

**Prefer envtest for controller behaviour.** It gives far better feedback than unit tests with mocked clients and runs much faster than end-to-end tests. A controller change with only unit tests should explain why envtest wasn't possible.

All levels:
1. **Tests must actually exercise the feature.** Check that fakes, such as a fake API server, expose the capability being tested. A test that passes because the fake lacks the feature proves nothing.
2. **Cover three cases for new behaviour:** existing behaviour unchanged, new behaviour works, and the failure path.
3. **Test at the level that owns the behaviour** (table above).
4. **Keep test setup simple.** Inline data is often clearer than overriding methods. Use existing builders and helpers.
5. **Don't change test data unless the test needs it.** Test data doesn't have to mirror production.
6. **Check that assertions can fail.** For example, `-1 <= 1` passes when the value is an error sentinel.

Waiting and polling (envtest and end-to-end):
1. **Poll, don't sleep.** Use `Eventually` and the existing wait helpers.
2. **Assertions inside `Eventually` use the provided handle:** Gomega's `Eventually(func(g Gomega) { g.Expect(...) })`, or testify's `require.EventuallyWithT(t, func(c *assert.CollectT) { ... })`. A plain `Expect` or `require` inside the polled function fails the test on the first attempt, so the check stops being eventual.
3. **Timeouts should be realistic.** A step that needs ten minutes means the test has already failed.

## 7. Build and tooling

1. **Makefile values must be overridable,** using `?=` variables. Nothing hardcoded that breaks local builds.
2. **Call tools through their Makefile variables,** for example `$(HELM)`, and mark targets `.PHONY`.
3. **Don't duplicate what the build platform already does,** such as setting `GOARCH` when buildx `--platform` handles it.

## 8. Change hygiene

1. **No unrelated changes.** No mass re-indents, no touching files without a change, no reformatting outside the diff.
2. **Every change has an issue.** Follow-ups get a ticket in the repo, not a TODO.
3. **CHANGELOG entry** in the repo's format (see each repo's `REVIEW.md`).
4. **Classify changes correctly.** If nothing was broken, it isn't a bugfix.
5. **Changed assertions in existing tests need a human.** If a PR modifies or removes an assertion in an existing test (a different expected value, a relaxed matcher, a deleted check, a skipped test), the reviewer flags it as Important and the PR requires human approval regardless of risk tier. Sometimes the change is correct because behaviour changed on purpose. But it is also the easiest way to hide a newly introduced bug. The PR description must say why the expectation changed.
6. **New lint suppressions need a human.** A new `//nolint` directive must name the specific linter and give a reason (`nolintlint` enforces both). Even then, the PR needs human approval regardless of risk tier. Suppressions are sometimes necessary, but each one switches off a check for good.
