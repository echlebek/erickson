package db

import (
	"errors"
	"fmt"

	"github.com/echlebek/erickson/review"
)

var ErrNoDB = errors.New("the database does not exist or is corrupt")

type ErrNoReview int
type ErrNoRevision int
type ErrNoAnnotation int

func (e ErrNoReview) Error() string {
	return fmt.Sprintf("review %d does not exist", e)
}

func (e ErrNoRevision) Error() string {
	return fmt.Sprintf("revision %d does not exist", e)
}

func (e ErrNoAnnotation) Error() string {
	return fmt.Sprintf("annotation %d does not exist", e)
}

// Database defines erickson's storage interface. erickson will only use
// the methods defined here, and will not introspect types implementing
// Database in any way.
type Database interface {
	// CreateReview creates a new review. It returns the ID of the review and
	// an error if the review could not be created.
	CreateReview(review.R) (id int, err error)

	// GetReview gets a review by ID.
	GetReview(id int) (review.R, error)

	// GetSummaries gets all the ReviewSummaries.
	GetSummaries() ([]review.Summary, error)

	// SetSummary sets the Summary of a Review by ID.
	SetSummary(id int, summary review.Summary) error

	// DeleteReview deletes a review.
	DeleteReview(id int) error

	// AddRevision adds a Revision to a Review.
	AddRevision(id int, r review.Revision) error

	// AddAnnotation adds an annotation to a Revision.
	AddAnnotation(id, revId int, a review.Annotation) error
}
