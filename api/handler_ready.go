package api

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, fmt.Sprintf("method not allowed: %s", r.Method))
		return
	}

	respondJSON(w, http.StatusOK, &response{Status: "ok"})
}
