package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (a *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(auth.APIKeyTokenType, r.Header)
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

func (a *apiConfig) middlewareJWT(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(auth.AccessTokenType, r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userID, err := auth.ValidateJWT(token, a.jwtSecret)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "couldn't validate JWT")
			return
		}
		id, err := uuid.Parse(userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "error parsing UUID stirng")
			return
		}

		user, err := a.DB.GetUserByID(r.Context(), uuidToPgType(id))
		if err != nil {
			respondError(w, http.StatusNotFound, "user not found")
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
