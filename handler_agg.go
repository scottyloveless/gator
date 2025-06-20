package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/scottyloveless/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	timeLayout := "Mon, 02 Jan 2006 15:04:05 +0000"
	for _, item := range feedData.Channel.Item {
		parsedTime, err := time.Parse(timeLayout, item.PubDate)
		if err != nil {
			log.Printf("error parsing time: %v", err)
			return
		}
		post, err := db.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsedTime,
			FeedID:      feed.ID,
		})
		if strings.Contains(fmt.Sprintf("%v", err), "duplicate key value violates unique constraint") {
			continue
		}
		if err != nil {
			log.Printf("error adding post %v to database: %v\n", post.Title, err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
