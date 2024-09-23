package main

import (
	"net/http"

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
