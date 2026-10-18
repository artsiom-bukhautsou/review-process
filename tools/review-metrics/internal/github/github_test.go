package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/example/review-metrics/internal/cache"
	"github.com/example/review-metrics/internal/model"
)

// fakeDoer serves canned JSON by exact path and records requests.
type fakeDoer struct {
	mu        sync.Mutex
	responses map[string]interface{}
	calls     []string
}

func (f *fakeDoer) DoWithContext(_ context.Context, method, path string, _ io.Reader, resp interface{}) error {
	f.mu.Lock()
	f.calls = append(f.calls, path)
	body, ok := f.responses[path]
	f.mu.Unlock()
	if !ok {
		body = []interface{}{}
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, resp)
}

func closedPR(n int, created time.Time, merged bool) map[string]interface{} {
	m := map[string]interface{}{"number": n, "created_at": created, "user": map[string]string{"login": "a"}}
	if merged {
		m["merged_at"] = created.Add(time.Hour)
	}
	return m
}

func TestMergedPRsPagesAndStopsAtSince(t *testing.T) {
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	var page1 []interface{}
	for i := 0; i < perPage; i++ {
		page1 = append(page1, closedPR(1000-i, since.Add(time.Duration(1000-i)*time.Hour), i%2 == 0))
	}
	page2 := []interface{}{
		closedPR(900, since.Add(2*time.Hour), true),
		closedPR(899, since.Add(-time.Hour), true), // older than since: stop here
		closedPR(898, since.Add(-2*time.Hour), true),
	}
	d := &fakeDoer{responses: map[string]interface{}{
		"repos/o/r/pulls?state=closed&sort=created&direction=desc&per_page=100&page=1": page1,
		"repos/o/r/pulls?state=closed&sort=created&direction=desc&per_page=100&page=2": page2,
	}}
	f := NewFetcher(d, nil, "github.com", 2)

	prs, err := f.MergedPRs(context.Background(), "o/r", since, 0)
	require.NoError(t, err)
	require.Len(t, prs, 51)
	require.Equal(t, 900, prs[50].Number)
	require.Len(t, d.calls, 2)

	prs, err = f.MergedPRs(context.Background(), "o/r", since, 10)
	require.NoError(t, err)
	require.Len(t, prs, 10)
}

func TestDetailsUsesCacheAndBuildsPageURLs(t *testing.T) {
	var reviews []interface{}
	for i := 0; i < perPage; i++ {
		reviews = append(reviews, map[string]interface{}{"state": "COMMENTED", "user": map[string]string{"login": fmt.Sprint("u", i)}})
	}
	d := &fakeDoer{responses: map[string]interface{}{
		"repos/o/r/pulls/7/reviews?per_page=100&page=1": reviews,
		"repos/o/r/pulls/7/reviews?per_page=100&page=2": []interface{}{map[string]interface{}{"state": "APPROVED"}},
		"repos/o/r/pulls/7/commits?per_page=100&page=1": []interface{}{map[string]interface{}{"sha": "x", "parents": []map[string]string{{"sha": "a"}, {"sha": "b"}}}},
		"repos/o/r/actions/runs?head_sha=abc&event=pull_request&per_page=100&page=1": map[string]interface{}{
			"total_count": 1, "workflow_runs": []map[string]interface{}{{"name": "CI", "run_attempt": 2}},
		},
	}}
	c, err := cache.New(t.TempDir())
	require.NoError(t, err)
	f := NewFetcher(d, c, "github.com", 1)
	pr := model.PullRequest{Number: 7, Head: model.Ref{SHA: "abc"}}

	got, err := f.Details(context.Background(), "o/r", pr, true)
	require.NoError(t, err)
	require.Len(t, got.Reviews, perPage+1)
	require.Len(t, got.Commits, 1)
	require.Len(t, got.Commits[0].Parents, 2)
	require.Len(t, got.Runs, 1)
	require.Equal(t, 2, got.Runs[0].RunAttempt)
	first := len(d.calls)
	for _, p := range d.calls {
		require.False(t, strings.Contains(p, "repositories/"), p)
	}

	again, err := f.Details(context.Background(), "o/r", pr, true)
	require.NoError(t, err)
	require.Equal(t, got, again)
	require.Len(t, d.calls, first, "second fetch must be served from cache")
}
