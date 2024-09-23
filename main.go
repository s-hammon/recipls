package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	fmt.Println("Listening on port :8080...")
	log.Fatal(srv.ListenAndServe())
}
