package web

import (
	"embed"
	"net/http"
	"time"

	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

const (
	baseURL      = "http://localhost:8080/v1"
	templatePath = "web/templates"
)

//go:embed static/*
var staticFiles embed.FS

func NewService(db *database.Queries, jwtSecret string) *http.HandlerFunc {
	cfg := config{
		DB:        db,
		client:    &http.Client{Timeout: time.Second * 5},
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
	client    *http.Client
	jwtSecret string
}
