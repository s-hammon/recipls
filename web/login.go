package web

import "net/http"

const tmplLogin = "login.html"

func (c *config) renderLoginTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	tmpl := getTemplate(tmplLogin, nil)
	if err := tmpl.Execute(w, nil); err != nil {
		respondError(w, http.StatusNotFound, "couldn't find login page")
	}
}
