// Package model holds the subset of GitHub REST API payloads that review-metrics uses.
package model

import "time"

// User is a GitHub account as returned in REST payloads.
type User struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

// Ref is the head or base ref of a pull request.
type Ref struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

// PullRequest is an item from GET repos/{o}/{r}/pulls.
type PullRequest struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	User      *User      `json:"user"`
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at"`
	Head      Ref        `json:"head"`
	Draft     bool       `json:"draft"`
}

// Review is an item from GET repos/{o}/{r}/pulls/{n}/reviews.
type Review struct {
	User        *User     `json:"user"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// File is an item from GET repos/{o}/{r}/pulls/{n}/files.
type File struct {
	Filename  string `json:"filename"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// CommitParent is a parent reference of a commit.
type CommitParent struct {
	SHA string `json:"sha"`
}

// Commit is an item from GET repos/{o}/{r}/pulls/{n}/commits.
type Commit struct {
	SHA     string         `json:"sha"`
	Parents []CommitParent `json:"parents"`
}

// Label is a label reference in a timeline event.
type Label struct {
	Name string `json:"name"`
}

// TimelineEvent is an item from GET repos/{o}/{r}/issues/{n}/timeline.
// Only the fields shared by the events we care about are decoded.
type TimelineEvent struct {
	Event     string    `json:"event"`
	CreatedAt time.Time `json:"created_at"`
	Actor     *User     `json:"actor"`
	Label     *Label    `json:"label"`
}

// WorkflowRun is an item from GET repos/{o}/{r}/actions/runs.
type WorkflowRun struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	WorkflowID int64  `json:"workflow_id"`
	HeadSHA    string `json:"head_sha"`
	RunAttempt int    `json:"run_attempt"`
	Conclusion string `json:"conclusion"`
}

// PRData is everything fetched for one merged pull request.
type PRData struct {
	Repo     string          `json:"repo"`
	PR       PullRequest     `json:"pr"`
	Reviews  []Review        `json:"reviews"`
	Files    []File          `json:"files"`
	Commits  []Commit        `json:"commits"`
	Timeline []TimelineEvent `json:"timeline"`
	// Runs is nil unless CI data was requested.
	Runs []WorkflowRun `json:"runs,omitempty"`
}

// AuthorLogin returns the PR author's login, or "ghost" for deleted accounts.
func (p PRData) AuthorLogin() string {
	if p.PR.User == nil || p.PR.User.Login == "" {
		return "ghost"
	}
	return p.PR.User.Login
}
