package metrics

import (
	"math"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/example/review-metrics/internal/model"
)

// Percentile returns the nearest-rank percentile p (0-100) of values using index
// round(p/100*(n-1)). ok is false for an empty input.
func Percentile(values []float64, p float64) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}
	s := append([]float64(nil), values...)
	sort.Float64s(s)
	k := int(math.Round(p / 100 * float64(len(s)-1)))
	k = max(0, min(len(s)-1, k))
	return s[k], true
}

// NewStat summarises values as count, median and p90.
func NewStat(values []float64) Stat {
	st := Stat{Count: len(values)}
	if m, ok := Percentile(values, 50); ok {
		st.Median = m
	}
	if p, ok := Percentile(values, 90); ok {
		st.P90 = p
	}
	if len(values) > 0 {
		var sum float64
		for _, v := range values {
			sum += v
		}
		st.Mean = sum / float64(len(values))
	}
	return st
}

// IsBot reports whether a user is a bot account.
func IsBot(u *model.User) bool {
	if u == nil {
		return false
	}
	return u.Type == "Bot" || strings.HasSuffix(u.Login, "[bot]")
}

// HumanReviews returns submitted reviews by non-bot users other than the author,
// ordered by submission time.
func HumanReviews(pr model.PRData) []model.Review {
	author := pr.AuthorLogin()
	var out []model.Review
	for _, r := range pr.Reviews {
		if r.User == nil || r.User.Login == "" || IsBot(r.User) {
			continue
		}
		if strings.EqualFold(r.User.Login, author) {
			continue
		}
		if r.State == "PENDING" || r.SubmittedAt.IsZero() {
			continue
		}
		out = append(out, r)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].SubmittedAt.Before(out[j].SubmittedAt) })
	return out
}

// AreaOf maps a file path to its area: the first depth directory components.
// Files with fewer directory components use what they have; root files map to "(root)".
func AreaOf(file string, depth int) string {
	dir := path.Dir(file)
	if dir == "." || dir == "/" || dir == "" {
		return "(root)"
	}
	parts := strings.Split(dir, "/")
	if depth < 1 {
		depth = 1
	}
	if len(parts) > depth {
		parts = parts[:depth]
	}
	return strings.Join(parts, "/")
}

var (
	followUpLine = regexp.MustCompile(`(?im)^[\s>*_\-]*(?:follow-?up to|regression from|fixes problem (?:introduced in|from))[ \t*_]*:[ \t*_]*((?:#\d+[ \t,;]*(?:and[ \t]+)?)+)`)
	issueRef     = regexp.MustCompile(`#(\d+)`)
)

// ParseFollowUps returns the PR numbers referenced by "Follow-up to: #N" style
// lines in a PR body (also "Followup to:", "Regression from:",
// "Fixes problem introduced in:", "Fixes problem from:"), deduplicated in order.
func ParseFollowUps(body string) []int {
	var out []int
	seen := map[int]bool{}
	for _, m := range followUpLine.FindAllStringSubmatch(body, -1) {
		for _, ref := range issueRef.FindAllStringSubmatch(m[1], -1) {
			n, err := strconv.Atoi(ref[1])
			if err != nil || seen[n] {
				continue
			}
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}

// BranchUpdates counts "update from base" events in a PR: merge commits (two or more
// parents) and head_ref_force_pushed timeline events (rebases).
func BranchUpdates(pr model.PRData) (int, int) {
	merges := 0
	for _, c := range pr.Commits {
		if len(c.Parents) >= 2 {
			merges++
		}
	}
	pushes := 0
	for _, e := range pr.Timeline {
		if e.Event == "head_ref_force_pushed" {
			pushes++
		}
	}
	return merges, pushes
}

// ReadyLabelTime returns when label was first applied to the PR.
func ReadyLabelTime(pr model.PRData, label string) (time.Time, bool) {
	var first time.Time
	found := false
	for _, e := range pr.Timeline {
		if e.Event != "labeled" || e.Label == nil || !strings.EqualFold(e.Label.Name, label) {
			continue
		}
		if !found || e.CreatedAt.Before(first) {
			first, found = e.CreatedAt, true
		}
	}
	return first, found
}

// LabelToReview returns hours from the first application of label to the first human
// review at or after it. labeled is false if the label was never applied; reviewed is
// false if no human review followed it.
func LabelToReview(pr model.PRData, label string) (float64, bool, bool) {
	at, ok := ReadyLabelTime(pr, label)
	if !ok {
		return 0, false, false
	}
	for _, r := range HumanReviews(pr) {
		if !r.SubmittedAt.Before(at) {
			return hours(r.SubmittedAt.Sub(at)), true, true
		}
	}
	return 0, true, false
}

func hours(d time.Duration) float64 {
	return d.Hours()
}

// counter keeps insertion-independent counts with deterministic ordering helpers.
type counter map[string]int

func (c counter) top(n int, total int) []Count {
	out := make([]Count, 0, len(c))
	for k, v := range c {
		out = append(out, Count{Name: k, Count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Name < out[j].Name
	})
	if n > 0 && len(out) > n {
		out = out[:n]
	}
	if total > 0 {
		for i := range out {
			out[i].Share = float64(out[i].Count) / float64(total)
		}
	}
	return out
}
