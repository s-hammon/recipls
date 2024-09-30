package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/app"
	"github.com/s-hammon/recipls/internal/database"

	pgxUUID "github.com/jackc/pgx-gofrs-uuid"
)

var (
	host               = flag.String("host", "0.0.0.0", "server host")
	port               = flag.Int("port", 8080, "listening port for server")
	rssRefreshInterval = flag.Int("rss-interval", 10, "interval (in minutes) by which the RSS feed updates")
)

func main() {
	flag.Parse()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
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

	domain := net.JoinHostPort(*host, strconv.Itoa(*port))
	app, err := app.New(dbQueries, domain)
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /index.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.RSSPath)
	})

	loggedMux := api.MiddlewareLogger(mux)

	srv := &http.Server{
		Addr:    domain,
		Handler: loggedMux,
	}

	requestInterval := time.Minute * time.Duration(*rssRefreshInterval)
	go app.RSSUpdateWorker(requestInterval)

	fmt.Printf("Listening on port %d...\n", *port)
	log.Fatal(srv.ListenAndServe())
}
