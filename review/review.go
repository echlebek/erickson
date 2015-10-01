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

// A revision is a set of diff.Files with corresponding annotations.
// Annotations are a list of comments on a given line of the diff,
// identified by file, hunk and line-number indices.
type Revision struct {
	Files       []diff.File  `json:"files"`
	Annotations []Annotation `json:"annotations"`
}

// Annotate annotates a file at a particular hunk and line number.
func (r *Revision) Annotate(a Annotation) {
	r.Annotations = append(r.Annotations, a)
}

type Annotation struct {
	File    int    `json:"file"`
	Hunk    int    `json:"hunk"`
	Line    int    `json:"line"`
	Comment string `json:"comment"`
	User    string `json:"user"`
}

// GetAnnotation gets the annotation at the index provided.
// If none exists, it returns an empty slice.
func (r Revision) GetAnnotations(file, hunk, line int) []Annotation {
	var result []Annotation
	for _, a := range r.Annotations {
		if a.File == file && a.Hunk == hunk && a.Line == line {
			result = append(result, a)
		}
	}
	return result
}
