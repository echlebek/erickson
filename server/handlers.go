package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/echlebek/erickson/diff"
	"github.com/echlebek/erickson/resource"
	"github.com/echlebek/erickson/review"
)

func home(ctx context, w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/reviews.html")
	if err != nil {
		log.Println(err)
		jsonError(w, errors.New("couldn't load resources"), 500)
		return
	}
	sums, err := ctx.db.GetSummaries()
	if err != nil {
		jsonError(w, err, 500)
		return
	}
	res := make([]resource.ReviewSummary, 0, len(sums))
	for _, s := range sums {
		ctx := ctx
		ctx.review = s.ID
		res = append(res, resource.ReviewSummary{Summary: s, URL: ctx.reviewURL()})
	}
	sort.Sort(resource.SummaryBySubmitTime(res))
	wrap := struct {
		Reviews []resource.ReviewSummary
	}{res}
	if err := tmpl.Execute(w, wrap); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

func getReview(ctx context, w http.ResponseWriter, req *http.Request) {
	review, err := ctx.db.GetReview(ctx.review)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	if len(review.Revisions) <= ctx.revision {
		http.NotFound(w, req)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(ctx.fsRoot, "static/html/review.html"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	res := resource.Review{
		R:                review,
		SelectedRevision: ctx.revision,
		URL:              ctx.reviewURL(),
	}
	if err := tmpl.Execute(w, res); err != nil {
		log.Println(err)
		return
	}
}

func postReview(ctx context, w http.ResponseWriter, req *http.Request) {
	switch req.Header.Get("Content-Type") {
	case "application/json":
		postJSONReview(ctx, w, req)
	case "application/x-www-form-urlencoded":
		postFormReview(ctx, w, req)
	}
}

func patchReview(ctx context, w http.ResponseWriter, req *http.Request) {
	p := struct {
		Status string `json:"status"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	review, err := ctx.db.GetReview(ctx.review)
	if err != nil {
		http.Error(w, fmt.Sprintf("review %d not found", ctx.review), http.StatusNotFound)
		return
	}
	review.Summary.Status = p.Status
	ctx.db.SetSummary(ctx.review, review.Summary)
	http.Redirect(w, req, "/", 303)
}

func postFormReview(ctx context, w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "couldn't parse form", 400)
		return
	}
	raw := req.FormValue("diff")
	files, err := diff.ParseFiles(raw)
	if err != nil {
		log.Println(err)
		http.Error(w, "couldn't parse diff", http.StatusBadRequest)
		return
	}
	r := review.R{
		Summary: review.Summary{
			CommitMsg:   req.FormValue("commitmsg"),
			Submitter:   req.FormValue("username"),
			Repository:  req.FormValue("repository"),
			SubmittedAt: time.Now(),
			Status:      review.Open,
		},
		Revisions: []review.Revision{
			{Files: files},
		},
	}
	id, err := ctx.db.CreateReview(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "couldn't create review", 500)
		return
	}
	ctx.review = id
	url := ctx.reviewURL()
	http.Redirect(w, req, url, 303)
}

func postJSONReview(ctx context, w http.ResponseWriter, req *http.Request) {
	var r review.R
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(r.Revisions) == 0 {
		http.Error(w, "no revisions provided", 400)
		return
	}
	r.Summary.SubmittedAt = time.Now()
	r.Summary.Status = review.Open
	id, err := ctx.db.CreateReview(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "couldn't create review", 500)
		return
	}
	ctx.review = id
	url := ctx.reviewURL()
	http.Redirect(w, req, url, 303)
}

func postRevision(ctx context, w http.ResponseWriter, req *http.Request) {
	var r review.Revision
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := ctx.db.AddRevision(ctx.review, r); err != nil {
		log.Println(err)
		http.Error(w, "couldn't add revision", 500)
		return
	}
	url := ctx.reviewURL()
	response := struct {
		URL string `json:"url"`
	}{url}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "couldn't write response", 500)
		return
	}
}

func patchRevision(ctx context, w http.ResponseWriter, req *http.Request) {
	var anno review.Annotation
	if err := json.NewDecoder(req.Body).Decode(&anno); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := ctx.db.AddAnnotation(ctx.review, ctx.revision, anno); err != nil {
		http.Error(w, "can't annotate review", 500)
		return
	}
}

func deleteReview(ctx context, w http.ResponseWriter, req *http.Request) {
	if err := ctx.db.DeleteReview(ctx.review); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
