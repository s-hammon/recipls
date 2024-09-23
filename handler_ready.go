package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}

	respondJSON(w, http.StatusOK, &response{Status: "ok"})
}
