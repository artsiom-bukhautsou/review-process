package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/example/review-metrics/internal/model"
)

var t0 = time.Date(2025, 1, 6, 9, 0, 0, 0, time.UTC)

func at(h float64) time.Time { return t0.Add(time.Duration(h * float64(time.Hour))) }

func user(login string) *model.User { return &model.User{Login: login, Type: "User"} }

func review(login, state string, h float64) model.Review {
	return model.Review{User: user(login), State: state, SubmittedAt: at(h)}
}

func pr(n int, author string, mergedH float64) model.PRData {
	m := at(mergedH)
	return model.PRData{
		Repo: "o/r",
		PR:   model.PullRequest{Number: n, User: user(author), CreatedAt: t0, MergedAt: &m},
	}
}

func TestPercentile(t *testing.T) {
	_, ok := Percentile(nil, 50)
	require.False(t, ok)

	v, ok := Percentile([]float64{7}, 90)
	require.True(t, ok)
	require.Equal(t, 7.0, v)

	vals := []float64{10, 1, 9, 2, 8, 3, 7, 4, 6, 5} // 1..10 unsorted
	v, _ = Percentile(vals, 50)
	require.Equal(t, 6.0, v, "index round(0.5*9)=5 -> 6th smallest")
	v, _ = Percentile(vals, 90)
	require.Equal(t, 9.0, v, "index round(0.9*9)=8 -> 9")
	v, _ = Percentile(vals, 0)
	require.Equal(t, 1.0, v)
	v, _ = Percentile(vals, 100)
	require.Equal(t, 10.0, v)
	require.Equal(t, []float64{10, 1, 9, 2, 8, 3, 7, 4, 6, 5}, vals, "input must not be mutated")

	st := NewStat([]float64{1, 2, 3})
	require.Equal(t, Stat{Count: 3, Median: 2, P90: 3, Mean: 2}, st)
	require.Equal(t, 0, NewStat(nil).Count)
}

func TestAreaOf(t *testing.T) {
	cases := map[string]string{
		"main.go":                         "(root)",
		"Makefile":                        "(root)",
		"pkg/foo.go":                      "pkg",
		"pkg/reconciliation/reconcile.go": "pkg/reconciliation",
		"pkg/reconciliation/sub/deep/file_test.go": "pkg/reconciliation",
		"config/crd/bases/x.yaml":                  "config/crd",
		".github/workflows/ci.yaml":                ".github/workflows",
	}
	for in, want := range cases {
		require.Equal(t, want, AreaOf(in, 2), in)
	}
	require.Equal(t, "pkg", AreaOf("pkg/reconciliation/reconcile.go", 1))
	require.Equal(t, "pkg/a/b", AreaOf("pkg/a/b/c/d.go", 3))
	require.Equal(t, "pkg/a", AreaOf("pkg/a/d.go", 3), "shallower paths keep what they have")
}

func TestParseFollowUps(t *testing.T) {
	cases := []struct {
		name string
		body string
		want []int
	}{
		{"plain", "Fixes a bug.\n\nFollow-up to: #123\n", []int{123}},
		{"case and no hyphen", "followup TO: #7", []int{7}},
		{"multiple refs", "Follow-up to: #1, #2 and #3", []int{1, 2, 3}},
		{"regression", "Regression from: #45", []int{45}},
		{"introduced in", "Fixes problem introduced in: #99 #100", []int{99, 100}},
		{"problem from", "Fixes problem from: #5", []int{5}},
		{"markdown bold and bullet", "- **Follow-up to:** #12\n", []int{12}},
		{"dedupe across lines", "Follow-up to: #8\nRegression from: #8 #9", []int{8, 9}},
		{"template placeholder", "Follow-up to: #\n", nil},
		{"not at line start", "This is a follow-up to: #1 maybe", nil},
		{"n/a", "Follow-up to: n/a", nil},
		{"windows newlines", "Summary\r\nFollow-up to: #77\r\n", []int{77}},
		{"empty", "", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, ParseFollowUps(c.body))
		})
	}
}

func TestBranchUpdates(t *testing.T) {
	p := pr(1, "alice", 10)
	p.Commits = []model.Commit{
		{SHA: "a", Parents: []model.CommitParent{{SHA: "x"}}},
		{SHA: "b", Parents: []model.CommitParent{{SHA: "a"}, {SHA: "main"}}},
		{SHA: "c", Parents: []model.CommitParent{{SHA: "b"}}},
		{SHA: "d", Parents: []model.CommitParent{{SHA: "c"}, {SHA: "main2"}}},
	}
	p.Timeline = []model.TimelineEvent{
		{Event: "head_ref_force_pushed"},
		{Event: "labeled"},
		{Event: "head_ref_force_pushed"},
		{Event: "committed"},
	}
	merges, pushes := BranchUpdates(p)
	require.Equal(t, 2, merges)
	require.Equal(t, 2, pushes)

	merges, pushes = BranchUpdates(pr(2, "bob", 1))
	require.Zero(t, merges)
	require.Zero(t, pushes)
}

func TestHumanReviewsExcludesBotsAndSelf(t *testing.T) {
	p := pr(1, "alice", 10)
	p.Reviews = []model.Review{
		review("carol", "COMMENTED", 5),
		review("alice", "COMMENTED", 1), // self
		{User: &model.User{Login: "coderabbitai[bot]", Type: "Bot"}, State: "COMMENTED", SubmittedAt: at(0.5)},
		{User: &model.User{Login: "some-app", Type: "Bot"}, State: "APPROVED", SubmittedAt: at(0.6)},
		{User: &model.User{Login: "copilot-pull-request-reviewer[bot]", Type: "User"}, State: "COMMENTED", SubmittedAt: at(0.7)},
		review("bob", "APPROVED", 3),
		{User: user("dave"), State: "PENDING"}, // not submitted
		{User: nil, State: "COMMENTED", SubmittedAt: at(2)},
	}
	got := HumanReviews(p)
	require.Len(t, got, 2)
	require.Equal(t, "bob", got[0].User.Login, "sorted by submission time")
	require.Equal(t, "carol", got[1].User.Login)

	require.True(t, IsBot(&model.User{Login: "dependabot[bot]"}))
	require.True(t, IsBot(&model.User{Login: "x", Type: "Bot"}))
	require.False(t, IsBot(&model.User{Login: "burmanm", Type: "User"}))
	require.False(t, IsBot(nil))
}

func TestLabelToReview(t *testing.T) {
	const label = "review/ready-for-human"

	p := pr(1, "alice", 50)
	p.Reviews = []model.Review{
		review("bob", "COMMENTED", 2),  // before label: ignored
		review("carol", "APPROVED", 9), // first human review after label
	}
	p.Timeline = []model.TimelineEvent{
		{Event: "labeled", CreatedAt: at(1), Label: &model.Label{Name: "kind/bug"}},
		{Event: "labeled", CreatedAt: at(4), Label: &model.Label{Name: "Review/Ready-For-Human"}},
		{Event: "unlabeled", CreatedAt: at(5), Label: &model.Label{Name: label}},
		{Event: "labeled", CreatedAt: at(6), Label: &model.Label{Name: label}},
	}
	h, labeled, reviewed := LabelToReview(p, label)
	require.True(t, labeled)
	require.True(t, reviewed)
	require.InDelta(t, 5.0, h, 1e-9, "from first label application (h4) to carol (h9)")

	noReview := pr(2, "alice", 50)
	noReview.Timeline = []model.TimelineEvent{{Event: "labeled", CreatedAt: at(1), Label: &model.Label{Name: label}}}
	_, labeled, reviewed = LabelToReview(noReview, label)
	require.True(t, labeled)
	require.False(t, reviewed)

	_, labeled, _ = LabelToReview(pr(3, "alice", 1), label)
	require.False(t, labeled)
}

func TestComputeEndToEnd(t *testing.T) {
	a := pr(10, "alice", 48)
	a.Reviews = []model.Review{review("bob", "COMMENTED", 4), review("bob", "APPROVED", 6), review("carol", "APPROVED", 8)}
	a.Files = []model.File{{Filename: "pkg/x/a.go", Additions: 10, Deletions: 5}, {Filename: "README.md", Additions: 1}}
	a.Commits = []model.Commit{{Parents: []model.CommitParent{{}, {}}}}

	b := pr(11, "alice", 24)
	b.PR.Body = "Follow-up to: #10"
	b.Reviews = []model.Review{review("alice", "COMMENTED", 1), {User: &model.User{Login: "bot[bot]"}, State: "APPROVED", SubmittedAt: at(1)}}
	b.Files = []model.File{{Filename: "pkg/x/a.go", Additions: 2, Deletions: 2}}
	b.Timeline = []model.TimelineEvent{{Event: "head_ref_force_pushed"}, {Event: "head_ref_force_pushed"}}

	c := pr(12, "bob", 12)
	c.PR.Title = "Revert \"something\""
	c.Reviews = []model.Review{review("alice", "APPROVED", 2)}
	c.Files = []model.File{{Filename: "apis/v1/types.go", Additions: 100, Deletions: 0}}
	c.Runs = []model.WorkflowRun{{Name: "Tests", RunAttempt: 2}, {Name: "Lint", RunAttempt: 1}}

	rep := Compute([]model.PRData{a, b, c}, []string{"o/r"}, t0, Options{
		Group: []string{"Alice"}, AreaDepth: 2, ReadyLabel: "review/ready-for-human", Top: 5, CI: true,
	})

	require.Equal(t, 3, rep.MergedPRs)
	require.Equal(t, 1, rep.NoHumanReview)
	require.Equal(t, []int{11}, rep.NoHumanReviewPRs)
	require.Equal(t, 2, rep.FirstReviewHours.Count)
	require.Equal(t, map[string]int{"0": 1, "1": 1, "2+": 1}, rep.Approvals)
	require.Equal(t, 4, rep.TotalReviews)
	require.Equal(t, "bob", rep.Reviewers[0].Login)
	require.Equal(t, 2, rep.Reviewers[0].Reviews)
	require.Equal(t, 1, rep.Reviewers[0].PRsReviewed)
	require.Equal(t, "alice", rep.Authors[0].Name)
	require.InDelta(t, 2.0/3, rep.Authors[0].Share, 1e-9)

	require.NotNil(t, rep.Group)
	require.Equal(t, 1, rep.Group.ReviewsByGroup)
	require.Equal(t, 3, rep.Group.ReviewsByOthers)
	require.InDelta(t, 2.0/3, rep.Group.AuthorshipShare, 1e-9)

	require.Equal(t, Stat{Count: 3, Median: 16, P90: 100, Mean: 40}, rep.SizeLines)

	// Areas: pkg/x (2 PRs: alice, bob, carol) , (root) (1 PR: alice, bob, carol), apis/v1 (1 PR: bob, alice)
	require.Len(t, rep.Areas, 3)
	require.Equal(t, "apis/v1", rep.Areas[0].Path)
	require.Equal(t, 2, rep.Areas[0].People)
	require.Equal(t, "pkg/x", rep.Areas[1].Path)
	require.Equal(t, 2, rep.Areas[1].PRs)
	require.Equal(t, Count{Name: "alice", Count: 2, Share: 1}, rep.Areas[1].TopPeople[0])

	require.Equal(t, 2, rep.BranchChurn.PRsWithUpdate)
	require.Equal(t, 1, rep.BranchChurn.MergeUpdates)
	require.Equal(t, 2, rep.BranchChurn.ForcePushes)
	require.InDelta(t, 1.0, rep.BranchChurn.MeanPerPR, 1e-9)
	require.Equal(t, 11, rep.BranchChurn.MaxPR)

	require.Equal(t, 1, rep.FollowUps.PRs)
	require.Equal(t, []FollowUp{{Repo: "o/r", Original: 10, FollowUp: 11}}, rep.FollowUps.Links)
	require.Equal(t, 1, rep.FollowUps.Reverts)

	require.False(t, rep.LabelTiming.InUse)

	require.NotNil(t, rep.CI)
	require.Equal(t, 1, rep.CI.SHAsWithRuns)
	require.Equal(t, 1, rep.CI.SHAsWithRerun)
	require.Equal(t, []Count{{Name: "Tests", Count: 1, Share: 1}}, rep.CI.TopWorkflows)
}
