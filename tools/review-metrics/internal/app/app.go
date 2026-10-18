// Package app wires fetching, metric computation and rendering together.
// It has no CLI or terminal dependencies so it can be driven by any front end.
package app

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/example/review-metrics/internal/cache"
	"github.com/example/review-metrics/internal/github"
	"github.com/example/review-metrics/internal/metrics"
	"github.com/example/review-metrics/internal/model"
	"github.com/example/review-metrics/internal/progress"
	"github.com/example/review-metrics/internal/report"
)

// Output formats.
const (
	FormatText     = "text"
	FormatJSON     = "json"
	FormatMarkdown = "markdown"
)

// Config is the user-facing configuration of a run.
type Config struct {
	Repos      []string
	Since      time.Time
	Limit      int
	Group      []string
	AreaDepth  int
	ReadyLabel string
	CI         bool
	Format     string
	Top        int
	Workers    int
}

// Env is everything a run needs from the outside world.
type Env struct {
	Client github.Doer
	Cache  *cache.Cache
	// Host namespaces cache entries (e.g. github.com or a GHES host).
	Host string
	Out  io.Writer
	Err  io.Writer
	// Styled renders text output for a terminal; nil means plain text.
	Styled func(io.Writer, []report.Document) error
	// ShowProgress draws a spinner on Err while fetching.
	ShowProgress bool
}

// ParseFormat normalises a --format value.
func ParseFormat(s string) (string, error) {
	switch strings.ToLower(s) {
	case "", "text", "table":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	}
	return "", fmt.Errorf("unknown format %q (want text, json or markdown)", s)
}

// ParseSince accepts YYYY-MM-DD or RFC 3339.
func ParseSince(s string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --since %q (want YYYY-MM-DD): %w", s, err)
	}
	return t, nil
}

// ValidateRepo checks that s looks like OWNER/REPO.
func ValidateRepo(s string) error {
	parts := strings.Split(s, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid repository %q: want OWNER/REPO", s)
	}
	return nil
}

// Run fetches data for every repo, computes metrics and writes the report.
func Run(ctx context.Context, cfg Config, env Env) error {
	format, err := ParseFormat(cfg.Format)
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	var repos []string
	for _, r := range cfg.Repos {
		if err := ValidateRepo(r); err != nil {
			return err
		}
		if key := strings.ToLower(r); !seen[key] {
			seen[key] = true
			repos = append(repos, r)
		}
	}
	cfg.Repos = repos

	var spin *progress.Spinner
	if env.ShowProgress {
		spin = progress.Start(env.Err, "listing pull requests…")
	}
	f := github.NewFetcher(env.Client, env.Cache, env.Host, cfg.Workers)
	f.Progress = func(repo string, done, total int) {
		spin.Set("%s: fetched details for %d/%d merged PRs", repo, done, total)
	}

	opts := metrics.Options{
		Group: cfg.Group, AreaDepth: cfg.AreaDepth, ReadyLabel: cfg.ReadyLabel, Top: cfg.Top, CI: cfg.CI,
	}
	var (
		reports []metrics.Report
		all     []model.PRData
	)
	for _, repo := range cfg.Repos {
		spin.Set("%s: listing merged pull requests…", repo)
		prs, err := f.FetchRepo(ctx, repo, github.Options{Since: cfg.Since, Limit: cfg.Limit, CI: cfg.CI})
		if err != nil {
			spin.Stop()
			return fmt.Errorf("fetching %s: %w", repo, err)
		}
		reports = append(reports, metrics.Compute(prs, []string{repo}, cfg.Since, opts))
		all = append(all, prs...)
	}
	spin.Stop()

	var combined *metrics.Report
	if len(cfg.Repos) > 1 {
		c := metrics.Compute(all, cfg.Repos, cfg.Since, opts)
		combined = &c
	}

	if format == FormatJSON {
		return report.JSON(env.Out, reports, combined)
	}
	docs := make([]report.Document, 0, len(reports)+1)
	for _, r := range reports {
		docs = append(docs, report.Build(r))
	}
	if combined != nil {
		docs = append(docs, report.Build(*combined))
	}
	if format == FormatMarkdown {
		return report.Markdown(env.Out, docs)
	}
	if env.Styled != nil {
		return env.Styled(env.Out, docs)
	}
	return report.Plain(env.Out, docs)
}
