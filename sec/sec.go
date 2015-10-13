package sec

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/scrypt"
)

type Base64 string

func (b Base64) MarshalJSON() ([]byte, error) {
	enc := base64.StdEncoding
	result := enc.EncodeToString([]byte(b))
	return json.Marshal(result)
}

func (b *Base64) UnmarshalJSON(src []byte) error {
	enc := base64.StdEncoding
	var stringsrc string
	if err := json.Unmarshal(src, &stringsrc); err != nil {
		return err
	}
	result, err := enc.DecodeString(stringsrc)
	*b = Base64(result)
	return err
}

// Credentials holds a user's hashed password and salt.
type Credentials struct {
	Name           string `json:"name"`
	Salt           Base64 `json:"salt"`
	HashedPassword Base64 `json:"hashedPassword"`
}

var ErrPasswordTooShort = errors.New("passwords must be at least 8 bytes long")
var ErrUsernameTooShort = errors.New("usernames must be at least 1 byte long")

const (
	n       = 32768
	p       = 16
	r       = 2
	keySize = 32
)

// NewCredentials creates a new user with a salt and hashed password
func NewCredentials(username string, password string) (Credentials, error) {
	u := Credentials{Name: username}
	if len(username) < 1 {
		return u, ErrUsernameTooShort
	}
	if len(password) < 8 {
		return u, ErrPasswordTooShort
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return u, err
	}
	u.Salt = Base64(salt)
	hashed, err := scrypt.Key([]byte(password), salt, n, p, r, keySize)
	if err != nil {
		return u, err
	}
	u.HashedPassword = Base64(hashed)
	return u, nil
}

// Authenticate checks if the password is correct.
func (c Credentials) Authenticate(password string) (bool, error) {
	got, err := scrypt.Key([]byte(password), []byte(c.Salt), n, p, r, keySize)
	if err != nil {
		return false, err
	}
	want := []byte(c.HashedPassword)
	if len(got) != len(want) {
		return false, nil
	}
	for i := range got {
		if want[i] != got[i] {
			return false, nil
		}
	}
	return true, nil
}
