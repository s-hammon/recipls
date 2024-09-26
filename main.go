package main

import (
	"context"
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
	DB  *database.Queries
	App *app.App
}

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
	fmt.Printf("using XML path: %s\n", app.RSSPath)
	cfg := apiConfig{DB: dbQueries, App: app}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /index.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.RSSPath)
	})

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handleGetUserByAPIKey))

	mux.HandleFunc("POST /v1/recipes", cfg.middlewareAuth(cfg.handlerCreateRecipe))
	mux.HandleFunc("GET /v1/recipes/{id}", cfg.handlerGetRecipeByID)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	const requestInterval = time.Minute * 10
	go cfg.rssUpdateWorker(requestInterval)

	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
