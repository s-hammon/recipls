package api

import (
	"net/http"

	"github.com/s-hammon/recipls/internal/database"
)

func NewService(db *database.Queries, jwtSecret string) *http.HandlerFunc {
	cfg := config{
		DB:        db,
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/metrics", cfg.middlewareAuth(cfg.handlerGetMetrics))
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("POST /v1/login", cfg.handlerLogin)
	mux.HandleFunc("POST /v1/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /v1/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareJWT(cfg.handleGetUserByAPIKey))

	mux.HandleFunc("GET /v1/recipes/{id}", cfg.handlerGetRecipeByID)
	mux.HandleFunc("PUT /v1/recipes/{id}", cfg.middlewareJWT(cfg.handlerUpdateRecipe))
	mux.HandleFunc("DELETE /v1/recipes/{id}", cfg.middlewareJWT(cfg.handlerDeleteRecipe))
	mux.HandleFunc("POST /v1/recipes", cfg.middlewareJWT(cfg.handlerCreateRecipe))

	loggedMux := MiddlewareLogger(mux)

	return &loggedMux
}

type config struct {
	DB        *database.Queries
	jwtSecret string
}

func NewConfig(db *database.Queries, jwtSecret string) config {
	return config{
		DB:        db,
		jwtSecret: jwtSecret,
	}
}
