package client

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/mail"
	"github.com/echlebek/erickson/review"
	"github.com/echlebek/erickson/server"
)

const diff = `diff --git a/server/server.go b/server/server.go
index c749bae..45afe23 100644
--- a/server/server.go
+++ b/server/server.go
@@ -102,6 +102,8 @@ func NewRootHandler(d db.Database, fsRoot string, sessionKey []byte) *RootH
        r.HandleFunc("/", ctx.authHandler(home)).
                Methods("GET")

+       r.HandleFunc("/", getCsrfToken).Methods("HEAD")
+
        r.HandleFunc("/reviews", ctx.authHandler(home)).
                Methods("GET")

@@ -153,7 +155,8 @@ func NewRootHandler(d db.Database, fsRoot string, sessionKey []byte) *RootH
        r.HandleFunc("/signup", ctx.handler(postSignup)).Methods("POST").
                Headers("Content-Type", "application/x-www-form-urlencoded")

-       r.HandleFunc("/login", ctx.handler(getLogin)).Methods("GET")
+       r.HandleFunc("/login", ctx.handler(getLogin)).
+               Methods("GET")

        r.HandleFunc("/login", ctx.handler(postLogin)).Methods("POST").
                Headers("Content-Type", "application/x-www-form-urlencoded")
`

func TestClient(t *testing.T) {
	server, done := newTestServer(t)
	defer func() {
		done <- true
	}()
	client, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	username := "bob"
	password := "supersecret"
	response, err := client.Authenticate(username, password)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode >= 400 {
		t.Fatalf("HTTP status %d", response.StatusCode)
	}
	commitmsg := "Changed some things"
	repo := "ToolTools"
	location, err := client.PostReview(diff, username, commitmsg, repo)
	if err != nil {
		t.Fatal(err)
	}
	if want := client.URL.String() + "/reviews/1"; location != want {
		t.Errorf("bad location: got %s, want %s", location, want)
	}
}

func newTestServer(t *testing.T) (*httptest.Server, chan bool) {
	done := make(chan bool)
	tmpd, err := ioutil.TempDir("/tmp", "erickson")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		<-done
		if err := os.RemoveAll(tmpd); err != nil {
			t.Fatal(err)
		}
	}()
	db, err := db.NewBoltDB(tmpd + "/erickson.db")
	if err != nil {
		t.Fatal(err)
	}
	u, err := review.NewUser("bob", "supersecret")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.CreateUser(u); err != nil {
		t.Fatal(err)
	}
	key := []byte("12345678901234567890123456789012")
	return httptest.NewServer(server.NewRootHandler(db, tmpd, key, mail.Nil)), done
}
