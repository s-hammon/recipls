package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
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
	cfg := apiConfig{DB: dbQueries, App: app}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("GET /index.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.RSSPath)
	})

	mux.HandleFunc("GET /recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getRequestID(r)
		if err != nil {
			respondError(w, http.StatusNotFound, "recipe not found 😔")
			return
		}

		recipeDB, err := cfg.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
		if err != nil {
			respondError(w, http.StatusInternalServerError, "error getting recipe")
			return
		}
		recipe := dbToRecipe(recipeDB)

		userDB, err := cfg.DB.GetUserByID(r.Context(), uuidToPgType(recipe.UserID))
		if err != nil {
			respondError(w, http.StatusInternalServerError, "error getting user")
			return
		}
		user := dbToUser(userDB)

		tmpl := getTemplate("recipe.html", template.FuncMap{"splitLines": splitLines})
		data := struct {
			Recipe Recipe
			User   User
		}{
			Recipe: recipe,
			User:   user,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("error executing template: %v", err)
			respondError(w, http.StatusInternalServerError, err.Error())
		}
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
