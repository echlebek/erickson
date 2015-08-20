package resource

import (
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
	return s[i].SubmittedAt.After(s[j].SubmittedAt)
}

func (s SummaryBySubmitTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SummaryBySubmitTime) Len() int {
	return len(s)
}
