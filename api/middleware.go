package api

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
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
			respondError(w, http.StatusNotFound, err.Error())
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
			respondError(w, http.StatusUnauthorized, "couldn't validate JWT")
			return
		}
		id, err := uuid.Parse(userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "error parsing UUID stirng")
			return
		}

		user, err := c.DB.GetUserByID(r.Context(), uuidToPgType(id))
		if err != nil {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}

		handler(w, r, user)
	}
}

func (c *config) middlewareSession(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("recipls_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			default:
				log.Println(err)
				respondError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		refreshToken, err := c.DB.GetRefreshTokenByValue(r.Context(), cookie.Value)
		if err != nil {
			slog.Info("redirecting user to login")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if refreshToken.ExpiresAt.Time.Before(time.Now().UTC()) {
			respondError(w, http.StatusUnauthorized, "cookie expired")
			return
		}

		user, err := c.DB.GetUserByID(r.Context(), refreshToken.UserID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "couldn't get user")
			return
		}

		handler(w, r, user)
	}
}

func MiddlewareLogger(handler http.Handler) http.HandlerFunc {
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
