package app

import (
	"context"
	"encoding/xml"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

const index = "index.xml"
const xmlPath = "content/xml"
const xmlTitle = "Recipls"
const xmlDescription = "A recipe feed"

type App struct {
	DB      *database.Queries
	RSSPath string
	RSSFeed Feed
	domain  string
}

func New(db *database.Queries, domain string) (*App, error) {
	if err := checkPath(xmlPath); err != nil {
		return nil, err
	}

	link := filepath.Join(domain, index)
	a := App{
		DB:      db,
		RSSPath: filepath.Join(xmlPath, index),
		RSSFeed: Feed{
			Version: "2.0",
			Channel: Channel{
				Title:       xmlTitle,
				Link:        link,
				Description: xmlDescription,
				PubDate:     time.Now().UTC(),
				Items:       []Item{},
			},
		},
		domain: domain,
	}

	if _, err := os.Stat(a.RSSPath); err != nil {
		if err := a.compile(); err != nil {
			return nil, err
		}
	}

	return &a, nil
}

func (a *App) AddItems() error {
	recipesDB, err := a.DB.GetRecipesWithLimit(context.Background(), 100)
	if err != nil {
		return err
	}

	items := []Item{}
	for _, recipe := range recipesDB {
		r := api.DBToRecipe(recipe)
		items = append(items, Item{
			Title:       r.Title,
			Link:        a.domain + "/recipes/" + r.ID.String(),
			Description: r.Description,
			PubDate:     r.CreatedAt,
		})
	}

	a.RSSFeed.Channel.Items = items
	rss, err := a.RSSFeed.Channel.generateRSSChannel()
	if err != nil {
		return err
	}

	return a.RefreshRSS(rss)
}

func (a *App) RefreshRSS(rss []byte) error {
	if err := checkPath(filepath.Dir(a.RSSPath)); err != nil {
		return err
	}

	return os.WriteFile(a.RSSPath, rss, os.ModePerm)
}

func (a *App) RSSUpdateWorker(requestInterval time.Duration) {
	ticker := time.NewTicker(requestInterval)

	for ; ; <-ticker.C {
		if err := a.AddItems(); err != nil {
			slog.Warn("couldn't update RSS feed", "error", err)
			continue
		}

		slog.Info("RSS feed updated")
	}
}

func (a *App) compile() error {
	xmlData, err := xml.MarshalIndent(a.RSSFeed, "", "  ")
	if err != nil {
		return err
	}

	b := []byte(xml.Header + string(xmlData))
	if _, err := os.Create(a.RSSPath); err != nil {
		return err
	}

	return os.WriteFile(a.RSSPath, b, os.ModePerm)
}

func checkPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}
