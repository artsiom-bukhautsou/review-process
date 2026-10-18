// Package metrics computes code review metrics from fetched PR data.
// Everything here is pure: no I/O, no clocks.
package metrics

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/example/review-metrics/internal/model"
)

// Options configures Compute.
type Options struct {
	// Group is a set of logins (e.g. current experts) compared against everyone else.
	Group []string
	// AreaDepth is the number of leading directories that make up an area.
	AreaDepth int
	// ReadyLabel is the label that marks a PR as ready for human review.
	ReadyLabel string
	// Top limits ranked lists.
	Top int
	// CI enables workflow re-run metrics (requires PRData.Runs).
	CI bool
}

// Stat summarises a distribution. Values are meaningless when Count is 0.
type Stat struct {
	Count  int     `json:"count"`
	Median float64 `json:"median"`
	P90    float64 `json:"p90"`
	Mean   float64 `json:"mean"`
}

// Count is a named count with its share of a total.
type Count struct {
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Share float64 `json:"share"`
}

// Reviewer is review activity for one person.
type Reviewer struct {
	Login       string  `json:"login"`
	Reviews     int     `json:"reviews"`
	PRsReviewed int     `json:"prs_reviewed"`
	Share       float64 `json:"share_of_reviews"`
}

// AuthorWait is the first-review wait for PRs by one author.
type AuthorWait struct {
	Author      string  `json:"author"`
	PRs         int     `json:"prs"`
	Reviewed    int     `json:"prs_with_human_review"`
	MedianHours float64 `json:"median_hours"`
}

// Area is knowledge spread for one path prefix.
type Area struct {
	Path      string  `json:"path"`
	PRs       int     `json:"prs"`
	People    int     `json:"people"`
	TopPeople []Count `json:"top_people"`
}

// GroupShare compares a group against everyone else.
type GroupShare struct {
	Members         []string `json:"members"`
	ReviewsByGroup  int      `json:"reviews_by_group"`
	ReviewsByOthers int      `json:"reviews_by_others"`
	ReviewShare     float64  `json:"review_share"`
	PRsByGroup      int      `json:"prs_by_group"`
	PRsByOthers     int      `json:"prs_by_others"`
	AuthorshipShare float64  `json:"authorship_share"`
}

// BranchChurn is metric 9: updates of a PR branch from its base.
type BranchChurn struct {
	PRsWithUpdate   int     `json:"prs_with_update"`
	ShareWithUpdate float64 `json:"share_with_update"`
	MeanPerPR       float64 `json:"mean_updates_per_pr"`
	MergeUpdates    int     `json:"merge_commit_updates"`
	ForcePushes     int     `json:"force_pushes"`
	MaxPerPR        int     `json:"max_updates_per_pr"`
	MaxPR           int     `json:"max_updates_pr,omitempty"`
}

// FollowUp links a PR to an earlier PR it follows up on.
type FollowUp struct {
	Repo     string `json:"repo"`
	Original int    `json:"original"`
	FollowUp int    `json:"follow_up"`
}

// FollowUps is metric 10.
type FollowUps struct {
	PRs         int        `json:"follow_up_prs"`
	Share       float64    `json:"share"`
	Links       []FollowUp `json:"links"`
	Reverts     int        `json:"revert_prs"`
	RevertShare float64    `json:"revert_share"`
}

// LabelTiming is metric 11.
type LabelTiming struct {
	Label         string `json:"label"`
	InUse         bool   `json:"in_use"`
	PRsLabeled    int    `json:"prs_labeled"`
	NoReviewAfter int    `json:"prs_without_review_after_label"`
	LabelToReview Stat   `json:"label_to_first_human_review_hours"`
}

// CIRuns is metric 12.
type CIRuns struct {
	SHAsWithRuns  int     `json:"head_shas_with_runs"`
	SHAsWithRerun int     `json:"head_shas_with_rerun"`
	RerunShare    float64 `json:"rerun_share"`
	TotalRuns     int     `json:"total_runs"`
	RerunRuns     int     `json:"rerun_runs"`
	TopWorkflows  []Count `json:"top_rerun_workflows"`
}

// Report holds all metrics for one repository or a combination of repositories.
type Report struct {
	Repos     []string  `json:"repos"`
	Since     time.Time `json:"since"`
	MergedPRs int       `json:"merged_prs"`

	FirstReviewHours Stat  `json:"first_human_review_hours"`
	NoHumanReview    int   `json:"prs_without_human_review"`
	NoHumanReviewPRs []int `json:"prs_without_human_review_numbers,omitempty"`
	OpenToMergeHours Stat  `json:"open_to_merge_hours"`

	HumanReviewsPerPR Stat           `json:"human_reviews_per_pr"`
	Approvals         map[string]int `json:"approvals_distribution"`

	Reviewers    []Reviewer  `json:"top_reviewers"`
	TotalReviews int         `json:"total_human_reviews"`
	Authors      []Count     `json:"top_authors"`
	Group        *GroupShare `json:"group,omitempty"`

	FirstReviewByAuthor []AuthorWait `json:"first_review_by_author"`

	SizeLines Stat `json:"pr_size_lines"`

	Areas []Area `json:"areas"`

	BranchChurn BranchChurn `json:"branch_update_churn"`
	FollowUps   FollowUps   `json:"follow_ups"`
	LabelTiming LabelTiming `json:"label_timing"`
	CI          *CIRuns     `json:"ci,omitempty"`
}

// Compute derives all metrics from prs. repos and since are echoed into the report.
func Compute(prs []model.PRData, repos []string, since time.Time, opts Options) Report {
	if opts.AreaDepth < 1 {
		opts.AreaDepth = 2
	}
	rep := Report{
		Repos:     repos,
		Since:     since,
		MergedPRs: len(prs),
		Approvals: map[string]int{"0": 0, "1": 0, "2+": 0},
	}
	group := map[string]bool{}
	for _, g := range opts.Group {
		if g = strings.TrimSpace(g); g != "" {
			group[strings.ToLower(g)] = true
		}
	}
	multiRepo := len(repos) > 1

	var (
		firstReview, openToMerge, reviewsPerPR, sizes []float64
		authorWaits                                   = map[string][]float64{}
		authorPRs                                     = counter{}
		reviewCounts                                  = counter{}
		reviewedPRs                                   = counter{}
		areaPRs                                       = map[string]int{}
		areaPeople                                    = map[string]counter{}
		gs                                            = GroupShare{}
		labelWaits                                    []float64
		updatesTotal                                  int
		workflowReruns                                = counter{}
	)

	for _, pr := range prs {
		author := pr.AuthorLogin()
		authorPRs[author]++
		if group[strings.ToLower(author)] {
			gs.PRsByGroup++
		} else {
			gs.PRsByOthers++
		}
		if pr.PR.MergedAt != nil {
			openToMerge = append(openToMerge, hours(pr.PR.MergedAt.Sub(pr.PR.CreatedAt)))
		}

		human := HumanReviews(pr)
		reviewsPerPR = append(reviewsPerPR, float64(len(human)))
		approvers := map[string]bool{}
		reviewersHere := map[string]bool{}
		for _, r := range human {
			login := r.User.Login
			reviewCounts[login]++
			reviewersHere[login] = true
			if r.State == "APPROVED" {
				approvers[strings.ToLower(login)] = true
			}
			if group[strings.ToLower(login)] {
				gs.ReviewsByGroup++
			} else {
				gs.ReviewsByOthers++
			}
		}
		for login := range reviewersHere {
			reviewedPRs[login]++
		}
		switch len(approvers) {
		case 0:
			rep.Approvals["0"]++
		case 1:
			rep.Approvals["1"]++
		default:
			rep.Approvals["2+"]++
		}

		if len(human) > 0 {
			w := hours(human[0].SubmittedAt.Sub(pr.PR.CreatedAt))
			firstReview = append(firstReview, w)
			authorWaits[author] = append(authorWaits[author], w)
		} else {
			rep.NoHumanReview++
			rep.NoHumanReviewPRs = append(rep.NoHumanReviewPRs, pr.PR.Number)
			if _, ok := authorWaits[author]; !ok {
				authorWaits[author] = nil
			}
		}

		size := 0
		for _, f := range pr.Files {
			size += f.Additions + f.Deletions
		}
		sizes = append(sizes, float64(size))

		// Knowledge spread: the author (unless a bot) plus human reviewers, once per PR.
		people := map[string]bool{}
		if !IsBot(pr.PR.User) && author != "ghost" {
			people[author] = true
		}
		for login := range reviewersHere {
			people[login] = true
		}
		areasHere := map[string]bool{}
		for _, f := range pr.Files {
			a := AreaOf(f.Filename, opts.AreaDepth)
			if multiRepo {
				a = pr.Repo + ":" + a
			}
			areasHere[a] = true
		}
		for a := range areasHere {
			areaPRs[a]++
			if areaPeople[a] == nil {
				areaPeople[a] = counter{}
			}
			for p := range people {
				areaPeople[a][p]++
			}
		}

		merges, pushes := BranchUpdates(pr)
		rep.BranchChurn.MergeUpdates += merges
		rep.BranchChurn.ForcePushes += pushes
		if u := merges + pushes; u > 0 {
			rep.BranchChurn.PRsWithUpdate++
			updatesTotal += u
			if u > rep.BranchChurn.MaxPerPR {
				rep.BranchChurn.MaxPerPR = u
				rep.BranchChurn.MaxPR = pr.PR.Number
			}
		}

		if refs := ParseFollowUps(pr.PR.Body); len(refs) > 0 {
			rep.FollowUps.PRs++
			for _, n := range refs {
				rep.FollowUps.Links = append(rep.FollowUps.Links, FollowUp{Repo: pr.Repo, Original: n, FollowUp: pr.PR.Number})
			}
		}
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(pr.PR.Title)), "revert") {
			rep.FollowUps.Reverts++
		}

		if opts.ReadyLabel != "" {
			w, labeled, reviewed := LabelToReview(pr, opts.ReadyLabel)
			if labeled {
				rep.LabelTiming.PRsLabeled++
				if reviewed {
					labelWaits = append(labelWaits, w)
				} else {
					rep.LabelTiming.NoReviewAfter++
				}
			}
		}

		if opts.CI && len(pr.Runs) > 0 {
			if rep.CI == nil {
				rep.CI = &CIRuns{}
			}
			rep.CI.SHAsWithRuns++
			rerun := false
			for _, r := range pr.Runs {
				rep.CI.TotalRuns++
				if r.RunAttempt > 1 {
					rerun = true
					rep.CI.RerunRuns++
					workflowReruns[r.Name]++
				}
			}
			if rerun {
				rep.CI.SHAsWithRerun++
			}
		}
	}

	n := len(prs)
	rep.FirstReviewHours = NewStat(firstReview)
	rep.OpenToMergeHours = NewStat(openToMerge)
	rep.HumanReviewsPerPR = NewStat(reviewsPerPR)
	rep.SizeLines = NewStat(sizes)

	for _, c := range reviewCounts {
		rep.TotalReviews += c
	}
	for _, c := range reviewCounts.top(opts.Top, rep.TotalReviews) {
		rep.Reviewers = append(rep.Reviewers, Reviewer{
			Login: c.Name, Reviews: c.Count, PRsReviewed: reviewedPRs[c.Name], Share: c.Share,
		})
	}
	rep.Authors = authorPRs.top(opts.Top, n)

	if len(group) > 0 {
		for g := range group {
			gs.Members = append(gs.Members, g)
		}
		sort.Strings(gs.Members)
		if t := gs.ReviewsByGroup + gs.ReviewsByOthers; t > 0 {
			gs.ReviewShare = float64(gs.ReviewsByGroup) / float64(t)
		}
		if n > 0 {
			gs.AuthorshipShare = float64(gs.PRsByGroup) / float64(n)
		}
		rep.Group = &gs
	}

	for _, a := range authorPRs.top(opts.Top, 0) {
		aw := AuthorWait{Author: a.Name, PRs: a.Count, Reviewed: len(authorWaits[a.Name])}
		if s := NewStat(authorWaits[a.Name]); s.Count > 0 {
			aw.MedianHours = s.Median
		}
		rep.FirstReviewByAuthor = append(rep.FirstReviewByAuthor, aw)
	}

	for a, c := range areaPRs {
		people := areaPeople[a]
		rep.Areas = append(rep.Areas, Area{Path: a, PRs: c, People: len(people), TopPeople: people.top(3, c)})
	}
	sort.Slice(rep.Areas, func(i, j int) bool {
		x, y := rep.Areas[i], rep.Areas[j]
		if x.People != y.People {
			return x.People < y.People
		}
		if x.PRs != y.PRs {
			return x.PRs > y.PRs
		}
		return x.Path < y.Path
	})

	if n > 0 {
		rep.BranchChurn.ShareWithUpdate = float64(rep.BranchChurn.PRsWithUpdate) / float64(n)
		rep.BranchChurn.MeanPerPR = float64(updatesTotal) / float64(n)
		rep.FollowUps.Share = float64(rep.FollowUps.PRs) / float64(n)
		rep.FollowUps.RevertShare = float64(rep.FollowUps.Reverts) / float64(n)
	}
	sort.Slice(rep.FollowUps.Links, func(i, j int) bool {
		a, b := rep.FollowUps.Links[i], rep.FollowUps.Links[j]
		if a.Repo != b.Repo {
			return a.Repo < b.Repo
		}
		if a.Original != b.Original {
			return a.Original < b.Original
		}
		return a.FollowUp < b.FollowUp
	})
	sort.Ints(rep.NoHumanReviewPRs)

	rep.LabelTiming.Label = opts.ReadyLabel
	rep.LabelTiming.InUse = rep.LabelTiming.PRsLabeled > 0
	rep.LabelTiming.LabelToReview = NewStat(labelWaits)

	if opts.CI {
		if rep.CI == nil {
			rep.CI = &CIRuns{}
		}
		if rep.CI.SHAsWithRuns > 0 {
			rep.CI.RerunShare = float64(rep.CI.SHAsWithRerun) / float64(rep.CI.SHAsWithRuns)
		}
		rep.CI.TopWorkflows = workflowReruns.top(opts.Top, rep.CI.RerunRuns)
	}
	return rep
}

// String renders a follow-up link for humans.
func (f FollowUp) String() string {
	return fmt.Sprintf("#%d was followed up by #%d", f.Original, f.FollowUp)
}
