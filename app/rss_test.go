package app

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestGenerateRSSFeed(t *testing.T) {
	pubDate1 := time.Now()
	pubDate2 := time.Now().Add(-1 * time.Hour)

	items := []Item{
		{
			Title:       "Title 1",
			Link:        "http://example.com/1",
			Description: "Description 1",
			PubDate:     pubDate1,
		},
		{
			Title:       "Title 2",
			Link:        "http://example.com/2",
			Description: "Description 2",
			PubDate:     pubDate2,
		},
	}

	feed := Channel{
		Title:       "Test RSS Feed",
		Link:        "http://example.com",
		Description: "A test RSS feed",
		PubDate:     pubDate1,
		Items:       items,
	}

	xmlData, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []byte(xml.Header + string(xmlData))
	got, err := feed.generateRSSChannel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(want) != string(got) {
		t.Errorf("want %s, got %s", string(want), string(got))
	}
}
