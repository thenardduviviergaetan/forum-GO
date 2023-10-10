package forum

import (
	"html/template"
	"net/http"
)

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// data := FetchData()
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/login.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
