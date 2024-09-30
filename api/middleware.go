package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

const (
	ErrUserNotFoundAPI = "couldn't find user with that API key"
	ErrUserNotFoundJWT = "couldn't find user with that JWT"
	ErrValidateJWT     = "couldn't validate JWT"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (c *config) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(auth.APIKeyTokenType, r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := c.DB.GetUserByAPIKey(r.Context(), token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, ErrUserNotFoundAPI)
			return
		}

		handler(w, r, user)
	}
}

func (c *config) middlewareJWT(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(auth.AccessTokenType, r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userID, err := auth.ValidateJWT(token, c.jwtSecret)
		if err != nil {
			respondError(w, http.StatusUnauthorized, ErrValidateJWT)
			return
		}
		id, err := uuid.Parse(userID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, ErrValidateJWT)
			return
		}

		user, err := c.DB.GetUserByID(r.Context(), uuidToPgType(id))
		if err != nil {
			respondError(w, http.StatusUnauthorized, ErrUserNotFoundJWT)
			return
		}

		handler(w, r, user)
	}
}

func MiddlewareLogger(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mw := &mwResponseWriter{w, http.StatusOK, nil}
		handler.ServeHTTP(mw, r)

		msg := fmt.Sprintf("%d %s %s", mw.StatusCode, r.Method, r.URL.Path)
		if mw.StatusCode > 499 {
			slog.Warn(msg, "error", string(mw.Message), "duration", time.Since(start))
		}
		slog.Info(msg, "duration", time.Since(start))
	}
}

type mwResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Message    []byte
}

func (w *mwResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *mwResponseWriter) Write(b []byte) (int, error) {
	w.Message = b
	return w.ResponseWriter.Write(b)
}
