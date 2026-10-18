// Package cmd defines the review-metrics command line.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/spf13/cobra"

	"github.com/example/review-metrics/internal/app"
	"github.com/example/review-metrics/internal/cache"
	"github.com/example/review-metrics/internal/report"
	"github.com/example/review-metrics/internal/report/styled"
)

type flags struct {
	since      string
	limit      int
	group      []string
	areaDepth  int
	readyLabel string
	ci         bool
	format     string
	cacheDir   string
	noCache    bool
	top        int
	workers    int
}

// NewRootCmd returns the review-metrics command.
func NewRootCmd() *cobra.Command {
	f := flags{}
	c := &cobra.Command{
		Use:   "review-metrics OWNER/REPO...",
		Short: "Code review metrics for merged pull requests",
		Long: `Compute code review metrics for merged pull requests in one or more GitHub
repositories: review latency, review depth and distribution, PR size, knowledge
spread per area, branch-update churn, follow-up fixes, ready-label timing and
(optionally) CI re-run rate.

Authentication and host selection follow the gh CLI (gh auth login, GH_TOKEN,
GH_HOST). Only the REST API is used.`,
		Example: `  review-metrics k8ssandra/cass-operator
  review-metrics k8ssandra/cass-operator --since 2024-10-01 --group alice,bob
  review-metrics org/op-a org/op-b --format markdown > review-metrics.md
  review-metrics org/op-a --ci --since 2026-07-01 --format json | jq .reports[0].ci`,
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args, f)
		},
	}
	defaultSince := time.Now().AddDate(0, -6, 0).Format("2006-01-02")
	fl := c.Flags()
	fl.StringVar(&f.since, "since", defaultSince, "only PRs created on or after this date (YYYY-MM-DD)")
	fl.IntVar(&f.limit, "limit", 500, "maximum merged PRs per repository")
	fl.StringSliceVar(&f.group, "group", nil, "comma-separated logins to compare against everyone else")
	fl.IntVar(&f.areaDepth, "area-depth", 2, "directory depth that defines an area for knowledge spread")
	fl.StringVar(&f.readyLabel, "ready-label", "review/ready-for-human", "label that marks a PR ready for human review")
	fl.BoolVar(&f.ci, "ci", false, "also fetch GitHub Actions runs to compute the re-run rate")
	fl.StringVar(&f.format, "format", app.FormatText, "output format: text, json or markdown")
	fl.StringVar(&f.cacheDir, "cache-dir", "", "response cache directory (default $XDG_CACHE_HOME/review-metrics)")
	fl.BoolVar(&f.noCache, "no-cache", false, "do not read or write the response cache")
	fl.IntVar(&f.top, "top", 6, "number of entries in ranked lists")
	fl.IntVar(&f.workers, "workers", 6, "concurrent API requests")
	return c
}

func run(cmd *cobra.Command, repos []string, f flags) error {
	since, err := app.ParseSince(f.since)
	if err != nil {
		return err
	}
	if _, err := app.ParseFormat(f.format); err != nil {
		return err
	}

	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("creating GitHub client (try `gh auth login`): %w", err)
	}
	host, _ := auth.DefaultHost()

	var c *cache.Cache
	if !f.noCache {
		dir := f.cacheDir
		if dir == "" {
			if dir, err = cache.DefaultDir(); err != nil {
				return err
			}
		}
		if c, err = cache.New(dir); err != nil {
			return err
		}
	}

	t := term.FromEnv()
	out := cmd.OutOrStdout()
	env := app.Env{
		Client:       client,
		Cache:        c,
		Host:         host,
		Out:          out,
		Err:          cmd.ErrOrStderr(),
		ShowProgress: term.IsTerminal(os.Stderr),
	}
	if t.IsTerminalOutput() && out == io.Writer(os.Stdout) {
		width, _, _ := t.Size()
		env.Styled = func(w io.Writer, docs []report.Document) error {
			return styled.Render(w, docs, width)
		}
	}

	cfg := app.Config{
		Repos:      repos,
		Since:      since,
		Limit:      f.limit,
		Group:      trimAll(f.group),
		AreaDepth:  f.areaDepth,
		ReadyLabel: f.readyLabel,
		CI:         f.ci,
		Format:     f.format,
		Top:        f.top,
		Workers:    f.workers,
	}
	return app.Run(cmd.Context(), cfg, env)
}

func trimAll(in []string) []string {
	var out []string
	for _, s := range in {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
