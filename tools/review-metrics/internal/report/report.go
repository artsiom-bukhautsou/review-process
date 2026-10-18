// Package report turns metrics into a renderer-neutral document and renders it as
// plain text, Markdown or JSON. Styled terminal output lives in report/styled.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/example/review-metrics/internal/metrics"
)

// KV is one "label: value" line.
type KV struct {
	Key   string
	Value string
}

// Table is a simple grid with a header row.
type Table struct {
	Title   string
	Headers []string
	Rows    [][]string
}

// Section is one titled block of the report.
type Section struct {
	Title  string
	Lines  []KV
	Tables []Table
	List   []string
	Note   string
}

// Document is a full report for one repo or a combination of repos.
type Document struct {
	Title    string
	Subtitle string
	Sections []Section
}

// Build converts a metrics report into a Document.
func Build(r metrics.Report) Document {
	doc := Document{
		Title: strings.Join(r.Repos, ", "),
		Subtitle: fmt.Sprintf("merged PRs created since %s: %d",
			r.Since.Format("2006-01-02"), r.MergedPRs),
	}
	if len(r.Repos) > 1 {
		doc.Title = "Combined: " + doc.Title
	}
	n := r.MergedPRs

	doc.Sections = append(doc.Sections, Section{
		Title: "Review latency",
		Lines: []KV{
			{"Time to first human review", statHours(r.FirstReviewHours)},
			{"Merged without human review", fmt.Sprintf("%d of %d (%s)", r.NoHumanReview, n, pct(r.NoHumanReview, n))},
			{"Open to merge", statHours(r.OpenToMergeHours)},
		},
		Note: "Waits are measured from PR creation; reviews by bots and by the PR author are ignored.",
	})

	approvals := fmt.Sprintf("0: %d (%s)   1: %d (%s)   2+: %d (%s)",
		r.Approvals["0"], pct(r.Approvals["0"], n), r.Approvals["1"], pct(r.Approvals["1"], n),
		r.Approvals["2+"], pct(r.Approvals["2+"], n))
	doc.Sections = append(doc.Sections, Section{
		Title: "Review depth",
		Lines: []KV{
			{"Human reviews per PR", fmt.Sprintf("median %s   mean %.1f", num(r.HumanReviewsPerPR.Median), r.HumanReviewsPerPR.Mean)},
			{"Distinct approvers per PR", approvals},
		},
	})

	dist := Section{Title: "Review and authorship distribution"}
	rt := Table{Title: "Top reviewers", Headers: []string{"Reviewer", "PRs reviewed", "Reviews", "Share of reviews"}}
	for _, rv := range r.Reviewers {
		rt.Rows = append(rt.Rows, []string{rv.Login, fmt.Sprint(rv.PRsReviewed), fmt.Sprint(rv.Reviews), share(rv.Share)})
	}
	at := Table{Title: "Top authors", Headers: []string{"Author", "PRs", "Share of PRs"}}
	for _, a := range r.Authors {
		at.Rows = append(at.Rows, []string{a.Name, fmt.Sprint(a.Count), share(a.Share)})
	}
	dist.Tables = []Table{rt, at}
	dist.Lines = []KV{{"Human reviews (total)", fmt.Sprint(r.TotalReviews)}}
	if g := r.Group; g != nil {
		dist.Lines = append(dist.Lines,
			KV{"Group", strings.Join(g.Members, ", ")},
			KV{"Reviews by group", fmt.Sprintf("%d (%s)   by others %d (%s)", g.ReviewsByGroup, share(g.ReviewShare), g.ReviewsByOthers, share(1-g.ReviewShare))},
			KV{"PRs authored by group", fmt.Sprintf("%d (%s)   by others %d (%s)", g.PRsByGroup, share(g.AuthorshipShare), g.PRsByOthers, share(1-g.AuthorshipShare))},
		)
	}
	doc.Sections = append(doc.Sections, dist)

	wt := Table{Headers: []string{"Author", "PRs", "With human review", "Median wait"}}
	for _, a := range r.FirstReviewByAuthor {
		wait := "n/a"
		if a.Reviewed > 0 {
			wait = hoursStr(a.MedianHours)
		}
		wt.Rows = append(wt.Rows, []string{a.Author, fmt.Sprint(a.PRs), fmt.Sprint(a.Reviewed), wait})
	}
	doc.Sections = append(doc.Sections, Section{Title: "First human review wait by author", Tables: []Table{wt}})

	doc.Sections = append(doc.Sections, Section{
		Title: "PR size",
		Lines: []KV{{"Lines changed (additions + deletions)", fmt.Sprintf("median %s   p90 %s", num(r.SizeLines.Median), num(r.SizeLines.P90))}},
	})

	ks := Table{Headers: []string{"Area", "PRs", "People", "Top people"}}
	for _, a := range r.Areas {
		var tops []string
		for _, p := range a.TopPeople {
			tops = append(tops, fmt.Sprintf("%s %d", p.Name, p.Count))
		}
		ks.Rows = append(ks.Rows, []string{a.Path, fmt.Sprint(a.PRs), fmt.Sprint(a.People), strings.Join(tops, ", ")})
	}
	doc.Sections = append(doc.Sections, Section{
		Title:  "Knowledge spread by area",
		Tables: []Table{ks},
		Note:   "People = PR author (non-bot) plus human reviewers; fewest people first.",
	})

	bc := r.BranchChurn
	churn := Section{
		Title: "Branch-update churn",
		Lines: []KV{
			{"PRs updated from base at least once", fmt.Sprintf("%d of %d (%s)", bc.PRsWithUpdate, n, share(bc.ShareWithUpdate))},
			{"Updates per PR", fmt.Sprintf("mean %.2f   max %d", bc.MeanPerPR, bc.MaxPerPR)},
			{"Merge-commit updates", fmt.Sprint(bc.MergeUpdates)},
			{"Force pushes (rebases)", fmt.Sprint(bc.ForcePushes)},
		},
		Note: "Updates = commits with 2+ parents in the PR plus head_ref_force_pushed events.",
	}
	if bc.MaxPR > 0 {
		churn.Lines[1].Value += fmt.Sprintf(" (#%d)", bc.MaxPR)
	}
	doc.Sections = append(doc.Sections, churn)

	fu := r.FollowUps
	fs := Section{
		Title: "Follow-up fixes",
		Lines: []KV{
			{"Follow-up PRs", fmt.Sprintf("%d of %d (%s)", fu.PRs, n, share(fu.Share))},
			{"Revert PRs (title starts with \"Revert\")", fmt.Sprintf("%d (%s)", fu.Reverts, share(fu.RevertShare))},
		},
		Note: "Follow-ups are PRs whose body has a line like \"Follow-up to: #N\".",
	}
	for _, l := range fu.Links {
		prefix := ""
		if len(r.Repos) > 1 {
			prefix = l.Repo + ": "
		}
		fs.List = append(fs.List, prefix+l.String())
	}
	doc.Sections = append(doc.Sections, fs)

	lt := r.LabelTiming
	ls := Section{Title: fmt.Sprintf("Ready label (%s)", lt.Label)}
	if !lt.InUse {
		ls.Lines = []KV{{"Label -> first human review", "n/a (label not in use yet)"}}
	} else {
		ls.Lines = []KV{
			{"PRs with the label", fmt.Sprintf("%d of %d (%s)", lt.PRsLabeled, n, pct(lt.PRsLabeled, n))},
			{"Label -> first human review", statHours(lt.LabelToReview)},
			{"Labeled but no human review after", fmt.Sprint(lt.NoReviewAfter)},
		}
	}
	doc.Sections = append(doc.Sections, ls)

	if ci := r.CI; ci != nil {
		cs := Section{
			Title: "CI re-runs",
			Lines: []KV{
				{"Head SHAs with any re-run", fmt.Sprintf("%d of %d with runs (%s)", ci.SHAsWithRerun, ci.SHAsWithRuns, share(ci.RerunShare))},
				{"Runs re-run", fmt.Sprintf("%d of %d", ci.RerunRuns, ci.TotalRuns)},
			},
			Note: "A run counts as re-run when run_attempt > 1 (pull_request runs of the merged head SHA only).",
		}
		if len(ci.TopWorkflows) > 0 {
			t := Table{Title: "Workflows by re-run count", Headers: []string{"Workflow", "Re-run runs", "Share"}}
			for _, w := range ci.TopWorkflows {
				t.Rows = append(t.Rows, []string{w.Name, fmt.Sprint(w.Count), share(w.Share)})
			}
			cs.Tables = []Table{t}
		}
		doc.Sections = append(doc.Sections, cs)
	}
	return doc
}

// JSON writes per-repo reports plus a combined report (when there are several) as JSON.
func JSON(w io.Writer, reports []metrics.Report, combined *metrics.Report) error {
	out := struct {
		Reports  []metrics.Report `json:"reports"`
		Combined *metrics.Report  `json:"combined,omitempty"`
	}{reports, combined}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	return nil
}

func statHours(s metrics.Stat) string {
	if s.Count == 0 {
		return "n/a"
	}
	return fmt.Sprintf("median %s   p90 %s   (n=%d)", hoursStr(s.Median), hoursStr(s.P90), s.Count)
}

func hoursStr(h float64) string {
	if h >= 48 {
		return fmt.Sprintf("%.1fh (%.1fd)", h, h/24)
	}
	return fmt.Sprintf("%.1fh", h)
}

func num(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprint(int64(v))
	}
	return fmt.Sprintf("%.1f", v)
}

func pct(a, b int) string {
	if b == 0 {
		return "n/a"
	}
	return share(float64(a) / float64(b))
}

func share(f float64) string {
	return fmt.Sprintf("%.1f%%", f*100)
}
