package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/echlebek/erickson/db"
	"github.com/gorilla/mux"
)

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
}

// handler performs partial function application on f to produce an http.HandlerFunc.
// The returned function will set the review and revision fields of c before
// calling f, providing f with context.
func (c context) handler(f func(context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		c.review, _ = strconv.Atoi(vars["id"])
		c.revision, _ = strconv.Atoi(vars["revision"])
		f(c, w, req)
	}
}

func (c context) reviewURL() (str string) {
	url, err := c.router.Get("review").URL("id", strconv.Itoa(c.review))
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
// /reviews/{id} GET, DELETE
// /reviews/{id}/rev POST
// /reviews/{id}/rev/{revision} GET, PATCH
//
//
func NewRootHandler(d db.Database, fsRoot string) *RootHandler {
	r := mux.NewRouter()
	handler := &RootHandler{Database: d, Router: r, URL: new(string), FSRoot: fsRoot}
	ctx := context{router: r, db: d, url: handler.URL, fsRoot: fsRoot}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", ctx.handler(home)).
		Methods("GET")

	r.HandleFunc("/reviews", ctx.handler(home)).
		Methods("GET")

	r.HandleFunc("/reviews/{id}", ctx.handler(getReview)).
		Name("review").Methods("GET")

	r.HandleFunc("/reviews/{id}", ctx.handler(deleteReview)).
		Methods("DELETE")

	r.HandleFunc("/reviews/{id}/rev/{revision}", ctx.handler(getReview)).
		Name("revision").
		Methods("GET")

	r.HandleFunc("/reviews", ctx.handler(postReview)).
		Methods("POST")

	r.HandleFunc("/reviews/{id}/rev", ctx.handler(postRevision)).
		Methods("POST").
		Headers("Accept", "application/json")

	r.HandleFunc("/reviews/{id}/rev/{revision}", ctx.handler(patchRevision)).
		Methods("PATCH")

	return handler
}
