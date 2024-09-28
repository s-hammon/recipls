package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/s-hammon/recipls/app"
	"github.com/s-hammon/recipls/internal/database"

	pgxUUID "github.com/jackc/pgx-gofrs-uuid"
)

const port = ":8080"

const xmlPath = "content/xml"
const xmlName = "Recipls"
const xmlDomain = "http://localhost" + port
const xmlDescription = "A recipe feed"

type apiConfig struct {
	DB        *database.Queries
	App       *app.App
	jwtSecret string
}

//go:embed static/*
var staticFiles embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("CONN_STRING")
	if dbURL == "" {
		log.Fatal("CONN_STRING must be set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("error parsing db url: %v", err)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	app, err := app.New(xmlPath, xmlName, xmlDomain, xmlDescription)
	if err != nil {
		log.Fatal(err)
	}
	cfg := apiConfig{
		DB:        dbQueries,
		App:       app,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(staticFiles)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("GET /index.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.RSSPath)
	})
	mux.HandleFunc("/login", cfg.renderLoginTemplate)
	mux.HandleFunc("/home", cfg.middlewareJWT(cfg.renderHomeTemplate))

	mux.HandleFunc("GET /recipes/{id}", cfg.renderRecipeTemplate)

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("POST /v1/login", cfg.handlerLogin)
	mux.HandleFunc("POST /v1/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /v1/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handleGetUserByAPIKey))

	mux.HandleFunc("GET /v1/recipes/{id}", cfg.handlerGetRecipeByID)
	mux.HandleFunc("PUT /v1/recipes/{id}", cfg.middlewareAuth(cfg.handlerUpdateRecipe))
	mux.HandleFunc("DELETE /v1/recipes/{id}", cfg.middlewareAuth(cfg.handlerDeleteRecipe))
	mux.HandleFunc("POST /v1/recipes", cfg.middlewareAuth(cfg.handlerCreateRecipe))

	loggedMux := cfg.middlewareLogger(mux)

	srv := &http.Server{
		Addr:    port,
		Handler: loggedMux,
	}

	const requestInterval = time.Minute * 10
	go cfg.rssUpdateWorker(requestInterval)

	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
