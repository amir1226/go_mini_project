package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/amir1226/go_mini_project/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRquest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration\n", concurrency, timeBetweenRquest)
	t := time.NewTicker(timeBetweenRquest)
	// This is an infinite loop that will run every timeBetweenRequest duration
	for ; ; <-t.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pub, err := time.Parse(time.RFC3339, item.PubDate)
		if err != nil {
			pub, err = time.Parse(time.RFC1123Z, item.PubDate)
		}
		if err != nil {
			log.Printf("error parsing date %v with error %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pub,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("error creating post:", err)
		}
	}
	log.Printf("Feed %s collected, %vposts found", feed.Name, len(rssFeed.Channel.Items))
}
