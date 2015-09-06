package review

import (
	"errors"
	"time"

	"github.com/echlebek/erickson/diff"
)

const (
	Open      = "Open"
	Submitted = "Submitted"
	Discarded = "Discarded"
)

var (
	ErrNotFound = errors.New("review does not exist")
)

// R is a code review. It is composed of a Summary and zero or more Revisions.
type R struct {
	Summary
	Revisions []Revision `json:"revisions"`
}

// A Summary contains cursory facts about a Review.
type Summary struct {
	ID          int       `json:"id"`
	CommitMsg   string    `json:"commitMsg"`
	Submitter   string    `json:"submitter"`
	SubmittedAt time.Time `json:"submittedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Repository  string    `json:"repository"`
	Status      string    `json:"status"`
}

// A revision is a set of diff.Files coupled with annotations.
type Revision struct {
	Files       []diff.File  `json:"files"`
	Annotations []Annotation `json:"annotations"`
}

// An Annotation is a message that corresponds to a file and line number in a
// patch set.
type Annotation struct {
	FileNumber int    `json:"hunkNumber"`
	LineNumber int    `json:"lineNumber"`
	Message    string `json:"message"`
}
