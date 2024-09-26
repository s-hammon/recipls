package app

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"time"
)

const index = "index.xml"

type App struct {
	RSSPath string
	RSSFeed Feed
}

func New(xmlPath, title, domain, description string) (*App, error) {
	if err := checkPath(xmlPath); err != nil {
		return nil, err
	}

	link := filepath.Join(domain, index)
	a := App{
		RSSPath: filepath.Join(xmlPath, index),
		RSSFeed: Feed{
			Version: "2.0",
			Channel: Channel{
				Title:       title,
				Link:        link,
				Description: description,
				PubDate:     time.Now().UTC(),
				Items:       []Item{},
			},
		},
	}

	if _, err := os.Stat(a.RSSPath); err != nil {
		if err := a.compile(); err != nil {
			return nil, err
		}
	}

	return &a, nil
}

func (a *App) AddItems(items []Item) error {
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

func checkPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
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
