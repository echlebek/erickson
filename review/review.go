package review

import (
	"encoding/json"
	"errors"
	"time"
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

// A Revision is a patch set, stored in Patches, along with Annotations that
// map to the files and line numbers in the patch set.
type Revision struct {
	Patches     string       `json:"patches"`
	Annotations []Annotation `json:"annotations"`
}

func (r *Revision) UnmarshalJSON(data []byte) error {
	// This function is implemented so that the patchSet field
	// will be populated upon unmarshaling the data.
	type revision Revision // avoid recursion on Unmarshal
	var rev revision
	var err error
	err = json.Unmarshal(data, &rev)
	if err != nil {
		return err
	}
	*r = Revision(rev)
	return nil
}

// An Annotation is a message that corresponds to a file and line number in a
// patch set.
type Annotation struct {
	FileNumber int    `json:"hunkNumber"`
	LineNumber int    `json:"lineNumber"`
	Message    string `json:"message"`
}
