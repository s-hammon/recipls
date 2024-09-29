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
	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/app"
	"github.com/s-hammon/recipls/internal/database"

	pgxUUID "github.com/jackc/pgx-gofrs-uuid"
)

const port = ":8080"

const xmlDomain = "http://localhost" + port

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("CONN_STRING")
	if dbURL == "" {
		log.Fatal("CONN_STRING must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set (make it good!)")
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

	apiSvc := api.NewService(dbQueries, jwtSecret)

	mux := http.NewServeMux()
	mux.Handle("/", apiSvc)

	app, err := app.New(dbQueries, xmlDomain)
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /index.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.RSSPath)
	})

	loggedMux := api.MiddlewareLogger(mux)

	srv := &http.Server{
		Addr:    port,
		Handler: loggedMux,
	}

	const requestInterval = time.Minute * 10
	go app.RSSUpdateWorker(requestInterval)

	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
