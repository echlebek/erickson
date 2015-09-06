package resource

import (
	"html/template"
	"time"

	"github.com/echlebek/erickson/review"
)

type ReviewSummary struct {
	URL string
	review.Summary
}

// Review
type Review struct {
	review.R
	SelectedRevision int
	URL              string
	Header           template.HTML
}

// CSS label class will be rendered according to the status
func (r ReviewSummary) StatusLabel() string {
	switch r.Status {
	case review.Open:
		return "label-primary"
	case review.Submitted:
		return "label-success"
	case review.Discarded:
		return "label-danger"
	}
	return "label-default"
}

func (r Review) StatusOpen() bool {
	return r.Status == review.Open
}

func (r ReviewSummary) SubmittedAt() string {
	return r.Summary.SubmittedAt.Format(time.UnixDate)
}

type DiffLine struct {
	LHS     string
	RHS     string
	lhsline int
	rhsline int
}

const (
	deletion     = "deletion"
	modification = "modification"
	insertion    = "insertion"
	unchanged    = "unchanged"
)

func (d DiffLine) LHSType() string {
	if d.LHS != "" && d.LHS[0] == '-' && d.RHS == "" {
		return "deletion"
	}
	if d.LHS != "" && d.LHS[0] == '-' && d.RHS != "" && d.RHS[0] == '+' {
		return "modification"
	}
	return unchanged
}

func (d DiffLine) RHSColour() string {
	if d.LHS != "" && d.LHS[0] == '-' && d.RHS != "" && d.RHS[0] == '+' {
		return modification
	}
	if d.LHS == "" && d.RHS != "" && d.RHS[0] == '+' {
		return insertion
	}
	return unchanged
}

type SummaryBySubmitTime []ReviewSummary

func (s SummaryBySubmitTime) Less(i, j int) bool {
	return s[i].Summary.SubmittedAt.After(s[j].Summary.SubmittedAt)
}

func (s SummaryBySubmitTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SummaryBySubmitTime) Len() int {
	return len(s)
}
