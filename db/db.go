package db

import (
	"errors"
	"fmt"

	"github.com/echlebek/erickson/review"
)

var (
	ErrNoDB       = errors.New("the database does not exist or is corrupt")
	ErrUserExists = errors.New("this username is already taken")
	ErrNoUser     = errors.New("no such user")
)

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

	// AddRevision adds a Revision to a Review.
	AddRevision(id int, r review.Revision) error

	// UpdateRevision replaces an existing revision with the one provided.
	UpdateRevision(id, revId int, r review.Revision) error

	// DeleteReview deletes a review.
	DeleteReview(id int) error

	// CreateUser creates a user.
	CreateUser(review.User) error

	// UpdateUser replaces an existing user with the one provided.
	UpdateUser(review.User) error

	// GetUser gets a user by name
	GetUser(string) (review.User, error)

	// DeleteUser deletes a user by name.
	DeleteUser(string) error
}
