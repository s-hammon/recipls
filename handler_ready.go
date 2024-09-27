package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, &response{Status: "error"})
		return
	}

	respondJSON(w, http.StatusOK, &response{Status: "ok"})
}
