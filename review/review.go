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

// Annotation is a message about a particular line in a hunk.
// It is made by particular user who may or may not have
// published it. If an annotation is published, other reviewers
// will be able to view it.
type Annotation struct {
	File      int    `json:"file"`
	Hunk      int    `json:"hunk"`
	Line      int    `json:"line"`
	Comment   string `json:"comment"`
	User      string `json:"user"`
	Published bool   `json:"published"`
}

// GetAnnotations gets the published annotations at the index provided.
// If none exists, it returns a nil slice.
func (r Revision) GetAnnotations(file, hunk, line int, currentUser string) []Annotation {
	var result []Annotation
	for _, a := range r.Annotations {
		if a.File == file && a.Hunk == hunk && a.Line == line && (a.Published || a.User == currentUser) {
			result = append(result, a)
		}
	}
	return result
}

// UnpublishedAnnotationsFor counts the number of un-published annotations
// made by user.
func (r Revision) UnpublishedAnnotationsFor(user string) (count int) {
	for _, a := range r.Annotations {
		if a.User == user && !a.Published {
			count += 1
		}
	}
	return
}
