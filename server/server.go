package server

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/echlebek/erickson/assets"
	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/mail"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const SessionName = "erickson"

type RootHandler struct {
	*mux.Router
	db.Database
	URL    *string
	FSRoot string
}

type context struct {
	router *mux.Router
	db     db.Database

	review   int
	revision int

	url    *string
	fsRoot string

	store *sessions.CookieStore

	mailer mail.Mailer

	// The request lead to auth failure. This is here because
	// I want to call getLogin after postLogin fails, but I don't
	// have any way for getLogin to know about the failure.
	// badAuth is passed to the template executor, so that the
	// page can show that the login failed.
	badAuth bool
}

type contextHandlerFunc func(context, http.ResponseWriter, *http.Request)

// authHandler checks the user's session before proceeding with the request.
func (c context) authHandler(f contextHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if session, err := c.store.Get(req, SessionName); err != nil || session.IsNew {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		vars := mux.Vars(req)
		c.review, _ = strconv.Atoi(vars["id"])
		c.revision, _ = strconv.Atoi(vars["revision"])
		f(c, w, req)
	}
}

// handler adds context to erickson's handlers.
func (c context) handler(f contextHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		f(c, w, req)
	}
}

func (c context) reviewURL() (str string) {
	var (
		url *url.URL
		err error
	)
	if c.revision > 0 {
		url, err = c.router.Get("review").URL("id", strconv.Itoa(c.review), "rev", strconv.Itoa(c.revision))
	} else {
		url, err = c.router.Get("review").URL("id", strconv.Itoa(c.review))
	}
	if err != nil {
		log.Println(err)
	}
	str = *c.url + url.String()
	return
}

func (c context) revisionURL(id, revision int) (str string) {
	sid := strconv.Itoa(id)
	srev := strconv.Itoa(revision)
	url, err := c.router.Get("revision").URL("id", sid, "revision", srev)
	if err != nil {
		log.Println(err)
	}
	str = *c.url + url.String()
	return
}

// NewRootHandler creates a RootHandler with the following routes defined:
//
// / GET
// /reviews GET, POST
// /reviews/{id} GET, PATCH, DELETE
// /reviews/{id}/rev POST
// /reviews/{id}/rev/{revision} GET, PATCH
// /reviews/{id}/annotations POST
// /reviews/{id}/rev/{revision}/annotations POST
//
//
func NewRootHandler(d db.Database, fsRoot string, sessionKey []byte, mailer mail.Mailer, urlRoot string) *RootHandler {
	r := mux.NewRouter()
	handler := &RootHandler{Database: d, Router: r, URL: &urlRoot, FSRoot: fsRoot}
	store := sessions.NewCookieStore(sessionKey)

	ctx := context{
		router: r,
		db:     d,
		url:    handler.URL,
		fsRoot: fsRoot,
		store:  store,
		mailer: mailer,
	}

	for name, handler := range assets.StylesheetHandlers {
		r.Handle("/assets/"+name, handler).Methods("GET")
	}

	for name, handler := range assets.ScriptHandlers {
		r.Handle("/assets/"+name, handler).Methods("GET")
	}

	r.HandleFunc("/", ctx.authHandler(home)).
		Methods("GET")

	r.HandleFunc("/", getCsrfToken).Methods("HEAD")

	r.HandleFunc("/reviews", ctx.authHandler(home)).
		Methods("GET")

	r.HandleFunc("/reviews/{id}", ctx.authHandler(getReview)).
		Name("review").Methods("GET")

	r.HandleFunc("/reviews/{id}", headReview).
		Name("review").Methods("HEAD")

	r.HandleFunc("/reviews/{id}", ctx.authHandler(deleteReview)).
		Methods("DELETE")

	r.HandleFunc("/reviews/{id}/status", ctx.authHandler(postStatus)).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/reviews/{id}/annotations", ctx.authHandler(postAnnotation)).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/reviews/{id}", ctx.authHandler(patchReview)).
		Methods("PATCH")

	r.HandleFunc("/reviews/{id}/rev/{revision}", ctx.authHandler(getReview)).
		Name("revision").
		Methods("GET")

	r.HandleFunc("/reviews/{id}/rev/{revision}/annotations", ctx.authHandler(postAnnotation)).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/reviews/{id}/rev/{revision}/annotations/publish", ctx.authHandler(publishAnnotations)).
		Methods("POST")

	r.HandleFunc("/reviews", ctx.authHandler(postJSONReview)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	r.HandleFunc("/reviews", ctx.authHandler(postFormReview)).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/reviews/{id}/rev", ctx.authHandler(postRevision)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	r.HandleFunc("/reviews/{id}/rev/{revision}", ctx.authHandler(patchRevision)).
		Methods("PATCH")

	r.HandleFunc("/signup", ctx.handler(getSignup)).Methods("GET")

	r.HandleFunc("/signup", ctx.handler(postSignup)).Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/login", ctx.handler(getLogin)).
		Methods("GET")

	r.HandleFunc("/login", ctx.handler(postLogin)).Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/logout", ctx.authHandler(postLogout)).Methods("POST")

	return handler
}
