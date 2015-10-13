package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/echlebek/erickson/assets"
	"github.com/echlebek/erickson/review"
	"github.com/gorilla/csrf"
)

func getSignup(ctx context, w http.ResponseWriter, req *http.Request) {
	wrap := map[string]interface{}{
		"Stylesheets":    assets.StylesheetHandlers,
		"Scripts":        assets.ScriptHandlers,
		csrf.TemplateTag: csrf.TemplateField(req),
	}
	if err := assets.Templates["signup.html"].Execute(w, wrap); err != nil {
		log.Println(err)
		http.Error(w, err500, 500)
		return
	}
}

func postSignup(ctx context, w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "couldn't parse form", 400)
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := review.NewUser(username, password)
	if err != nil {
		log.Println(err)
		http.Error(w, "couldn't create user", 500)
		return
	}
	if err := ctx.db.CreateUser(user); err != nil {
		log.Println(err)
		http.Error(w, "couldn't create user", 500)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func getLogin(ctx context, w http.ResponseWriter, req *http.Request) {
	wrap := struct {
		Stylesheets map[string]http.Handler
		Scripts     map[string]http.Handler
		CSRFField   template.HTML
		AuthFailure bool
	}{assets.StylesheetHandlers, assets.ScriptHandlers, csrf.TemplateField(req), false}
	if err := assets.Templates["login.html"].Execute(w, wrap); err != nil {
		log.Println(err)
		http.Error(w, err500, 500)
		return
	}
}

func postLogin(ctx context, w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "couldn't parse form", 400)
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := ctx.db.GetUser(username)
	if err != nil {
		http.Error(w, "bad credentials", 401)
		return
	}
	if ok, err := user.Authenticate(password); err != nil {
		http.Error(w, "bad credentials", 401)
		return
	} else if !ok {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	session, err := ctx.store.Get(req, SessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, err500, 500)
		return
	}
	session.Values["username"] = username
	if err := session.Save(req, w); err != nil {
		log.Println(err)
		http.Error(w, err500, 500)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
