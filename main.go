package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/s-hammon/recipls/internal/database"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("CONN_STRING")
	if dbURL == "" {
		log.Fatal("CONN_STRING must be set")
	}

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbConf, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	defer db.Close(context.Background())

	dbQueries := database.New(db)
	cfg := apiConfig{DB: dbQueries}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUser)

	srv := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	fmt.Println("Listening on port :8080...")
	log.Fatal(srv.ListenAndServe())
}
