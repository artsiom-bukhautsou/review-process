// Package github fetches pull request data over the GitHub REST API.
//
// It only depends on a small Doer interface, which go-gh's *api.RESTClient
// satisfies, so the transport (auth, GH_HOST, enterprise hosts) is go-gh's.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/example/review-metrics/internal/cache"
	"github.com/example/review-metrics/internal/model"
)

const perPage = 100

// Doer performs one REST request and decodes the JSON response.
// *github.com/cli/go-gh/v2/pkg/api.RESTClient implements it.
type Doer interface {
	DoWithContext(ctx context.Context, method string, path string, body io.Reader, response interface{}) error
}

// Fetcher retrieves merged PRs and their details with bounded concurrency.
type Fetcher struct {
	client  Doer
	cache   *cache.Cache
	host    string
	workers int
	// Progress, if set, is called after each PR's details are fetched.
	Progress func(repo string, done, total int)
}

// NewFetcher returns a Fetcher. host only namespaces cache keys; c may be nil.
func NewFetcher(client Doer, c *cache.Cache, host string, workers int) *Fetcher {
	if workers < 1 {
		workers = 1
	}
	return &Fetcher{client: client, cache: c, host: host, workers: workers}
}

// Options controls what FetchRepo retrieves.
type Options struct {
	Since time.Time
	Limit int
	CI    bool
}

// FetchRepo returns merged PRs created at or after opts.Since (newest first, at most
// opts.Limit) along with their reviews, files, commits, timeline and optionally CI runs.
func (f *Fetcher) FetchRepo(ctx context.Context, repo string, opts Options) ([]model.PRData, error) {
	prs, err := f.MergedPRs(ctx, repo, opts.Since, opts.Limit)
	if err != nil {
		return nil, err
	}
	out := make([]model.PRData, len(prs))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
		done     int
	)
	sem := make(chan struct{}, f.workers)
	for i, pr := range prs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()
			d, err := f.Details(ctx, repo, pr, opts.CI)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
					cancel()
				}
				return
			}
			out[i] = d
			done++
			if f.Progress != nil {
				f.Progress(repo, done, len(prs))
			}
		}()
	}
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return out, nil
}

// MergedPRs pages through closed PRs (newest created first) and keeps merged ones
// created at or after since. Paging stops at the first PR older than since.
func (f *Fetcher) MergedPRs(ctx context.Context, repo string, since time.Time, limit int) ([]model.PullRequest, error) {
	var merged []model.PullRequest
	for page := 1; ; page++ {
		path := fmt.Sprintf("repos/%s/pulls?state=closed&sort=created&direction=desc&per_page=%d&page=%d", repo, perPage, page)
		var batch []model.PullRequest
		// The PR list changes over time, so it is never cached.
		if err := f.get(ctx, path, false, &batch); err != nil {
			return nil, err
		}
		for _, pr := range batch {
			if pr.CreatedAt.Before(since) {
				return merged, nil
			}
			if pr.MergedAt == nil {
				continue
			}
			merged = append(merged, pr)
			if limit > 0 && len(merged) >= limit {
				return merged, nil
			}
		}
		if len(batch) < perPage {
			return merged, nil
		}
	}
}

// Details fetches reviews, files, commits and timeline (and CI runs if ci) for one PR.
// Merged PRs do not change, so these responses are cached.
func (f *Fetcher) Details(ctx context.Context, repo string, pr model.PullRequest, ci bool) (model.PRData, error) {
	d := model.PRData{Repo: repo, PR: pr}
	base := fmt.Sprintf("repos/%s", repo)
	n := pr.Number
	if err := getAll(ctx, f, fmt.Sprintf("%s/pulls/%d/reviews", base, n), &d.Reviews); err != nil {
		return d, err
	}
	if err := getAll(ctx, f, fmt.Sprintf("%s/pulls/%d/files", base, n), &d.Files); err != nil {
		return d, err
	}
	if err := getAll(ctx, f, fmt.Sprintf("%s/pulls/%d/commits", base, n), &d.Commits); err != nil {
		return d, err
	}
	if err := getAll(ctx, f, fmt.Sprintf("%s/issues/%d/timeline", base, n), &d.Timeline); err != nil {
		return d, err
	}
	if ci && pr.Head.SHA != "" {
		runs, err := f.Runs(ctx, repo, pr.Head.SHA)
		if err != nil {
			return d, err
		}
		d.Runs = runs
	}
	return d, nil
}

// Runs returns pull_request-triggered workflow runs for a head SHA.
func (f *Fetcher) Runs(ctx context.Context, repo, sha string) ([]model.WorkflowRun, error) {
	runs := []model.WorkflowRun{}
	for page := 1; ; page++ {
		path := fmt.Sprintf("repos/%s/actions/runs?head_sha=%s&event=pull_request&per_page=%d&page=%d",
			repo, url.QueryEscape(sha), perPage, page)
		var resp struct {
			TotalCount   int                 `json:"total_count"`
			WorkflowRuns []model.WorkflowRun `json:"workflow_runs"`
		}
		if err := f.get(ctx, path, true, &resp); err != nil {
			return nil, err
		}
		runs = append(runs, resp.WorkflowRuns...)
		if len(resp.WorkflowRuns) < perPage || len(runs) >= resp.TotalCount {
			return runs, nil
		}
	}
}

// getAll pages through a list endpoint by building ?page=N URLs itself (rather than
// following Link headers, which may point at repositories/{id} paths).
func getAll[T any](ctx context.Context, f *Fetcher, path string, out *[]T) error {
	items := []T{}
	for page := 1; ; page++ {
		var batch []T
		p := fmt.Sprintf("%s?per_page=%d&page=%d", path, perPage, page)
		if err := f.get(ctx, p, true, &batch); err != nil {
			return err
		}
		items = append(items, batch...)
		if len(batch) < perPage {
			*out = items
			return nil
		}
	}
}

func (f *Fetcher) get(ctx context.Context, path string, cacheable bool, out interface{}) error {
	key := f.host + "/" + path
	if cacheable {
		if b, ok := f.cache.Get(key); ok {
			if err := json.Unmarshal(b, out); err == nil {
				return nil
			}
		}
	}
	var raw json.RawMessage
	if err := f.client.DoWithContext(ctx, "GET", path, nil, &raw); err != nil {
		return fmt.Errorf("GET %s: %w", strings.SplitN(path, "?", 2)[0], err)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decoding %s: %w", path, err)
	}
	if cacheable {
		if err := f.cache.Put(key, raw); err != nil {
			return err
		}
	}
	return nil
}
