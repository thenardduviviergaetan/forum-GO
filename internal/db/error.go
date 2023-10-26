package forum

import (
	"net/http"
	"text/template"
	//"fmt"
)

func UserError(w http.ResponseWriter, r *http.Request, message string) {
	tmpl, err := template.ParseFiles(
		"web/templates/error.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ServerError(w, r, "Error 500: internal server error, original error: " + message)
		return
	}
	if err := tmpl.Execute(w, map[string]string{"ErrorMessage": message}); err != nil {
		ServerError(w, r, "Error 500: internal server error, original error: " + message)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, message string) {
	tmpl, err := template.ParseFiles(
		"web/templates/error.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, map[string]string{"ErrorMessage": message}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
