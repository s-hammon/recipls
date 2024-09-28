package main

import (
	"fmt"
	"net/http"

	"github.com/s-hammon/recipls/internal/auth"
)

func (a *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetToken("Bearer", r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = a.DB.DeleteRefreshTokenByValue(r.Context(), authToken); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't delete refresh token: %v", err))
	}

	w.WriteHeader(http.StatusNoContent)
}
