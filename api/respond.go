package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

type errResponse struct {
	Message   string    `json:"message"`
	RequestDT time.Time `json:"request_dt"`
}

func respondError(w http.ResponseWriter, code int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	payload := errResponse{
		Message:   errMessage,
		RequestDT: time.Now().UTC(),
	}
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
