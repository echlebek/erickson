package resource

import (
	"time"

	"github.com/echlebek/erickson/review"
)

// File is a resource that represents a file.
type File struct {
	Name string
	Type string
	URL  string
}

type ReviewSummary struct {
	URL string
	review.Summary
}

// Review
type Review struct {
	review.R
	SelectedRevision int
	Diff             []DiffLine
	URL              string
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
	LHS string
	RHS string
}

const (
	yellow = "#FDFCDC"
	red    = "#FA8072"
	green  = "#7FFFD4"
	bg     = "#FFFFFF"
)

func (d DiffLine) LHSColour() string {
	if d.LHS != "" && d.LHS[0] == '-' && d.RHS == "" {
		return red
	}
	if d.LHS != "" && d.LHS[0] == '-' && d.RHS != "" && d.RHS[0] == '+' {
		return yellow
	}
	return bg
}

func (d DiffLine) RHSColour() string {
	if d.LHS != "" && d.LHS[0] == '-' && d.RHS != "" && d.RHS[0] == '+' {
		return yellow
	}
	if d.LHS == "" && d.RHS != "" && d.RHS[0] == '+' {
		return green
	}
	return bg
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
