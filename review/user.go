package review

import "github.com/echlebek/erickson/sec"

// User holds the information about an erickson user's account.
type User struct {
	sec.Credentials
}

// NewUser creates a new user.
func NewUser(username, password string) (u User, err error) {
	u.Credentials, err = sec.NewCredentials(username, password)
	return
}
