package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

const ApiAuthKey = "ApiKey"

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (a *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(ApiAuthKey, r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := a.DB.GetUserByAPIKey(r.Context(), token)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		handler(w, r, user)
	}
}

func (a *apiConfig) middlewareLogger(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mw := &mwResponseWriter{w, http.StatusOK}
		handler.ServeHTTP(mw, r)

		msg := fmt.Sprintf("%d %s %s", mw.StatusCode, r.Method, r.URL.Path)
		slog.Info(msg, "duration", time.Since(start))
		if strings.Contains(r.URL.Path, "/static/") {
			fmt.Println()
		}
	}
}

type mwResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *mwResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}
