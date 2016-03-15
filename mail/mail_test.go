package mail

import (
	"testing"
)

func TestMail(t *testing.T) {
	m := Message{
		Sender:     "bob@example.com",
		Recipient:  "alice@example.com",
		Repository: "github.com/echlebek/erickson",
		ReviewURL:  "/reviews/1",

		Annotations: []Annotation{
			{},
		},
	}
	// If the template execution fails then there will be a panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	if err := Nil.NotifyReviewAnnotated(m); err != nil {
		t.Error(err)
	}
}
