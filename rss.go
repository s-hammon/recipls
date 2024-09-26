package main

import (
	"encoding/xml"
	"time"
)

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     time.Time `xml:"pubDate"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     time.Time `xml:"pubDate"`
	Items       []Item    `xml:"item"`
}

func (c Channel) generateRSSFeed() ([]byte, error) {
	xmlData, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}

	b := []byte(xml.Header + string(xmlData))
	return b, nil
}
