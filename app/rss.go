package app

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

func (c Channel) generateRSSChannel() ([]byte, error) {
	xmlData, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}

	rss := []byte(xml.Header + string(xmlData))
	return rss, nil
}

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}
