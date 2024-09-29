package web

import (
	"embed"
	"net/http"

	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

const templatePath = "web/templates"

//go:embed static/*
var staticFiles embed.FS

func NewService(db *database.Queries, jwtSecret string) *http.HandlerFunc {
	cfg := config{
		DB:        db,
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServer(http.FS(staticFiles)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("/login", cfg.renderLoginTemplate)
	mux.HandleFunc("/home", cfg.middlewareSession(cfg.renderHomeTemplate))

	mux.HandleFunc("GET /recipes/new", cfg.middlewareJWT(cfg.renderNewRecipeTemplate))
	mux.HandleFunc("GET /recipes/{id}", cfg.renderRecipeTemplate)

	loggedMux := api.MiddlewareLogger(mux)

	return &loggedMux
}

type config struct {
	DB        *database.Queries
	jwtSecret string
}
