package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSS struct {
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSS, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSS{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSS{}, err
	}

	var rss RSS
	err = xml.Unmarshal(dat, &rss)
	if err != nil {
		return RSS{}, err
	}

	return rss, nil
}
