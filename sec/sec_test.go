package sec

import "testing"

func TestAuth(t *testing.T) {
	username := "foobar@baz.com"
	password := "sup3rs3cr3t"
	badpassword := "supersecret"

	u, err := NewCredentials(username, password)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := u.Authenticate(password); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Errorf("auth error")
	}
	if ok, err := u.Authenticate(badpassword); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Errorf("auth error")
	}
}

func TestBadCredentials(t *testing.T) {
	if _, err := NewCredentials("", ""); err != ErrUsernameTooShort {
		t.Errorf("expected error")
	}
	if _, err := NewCredentials("foobar@baz.com", "123456"); err != ErrPasswordTooShort {
		t.Errorf("expected error")
	}
}
