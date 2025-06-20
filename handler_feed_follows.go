package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/scottyloveless/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		fmt.Println("not enough arguments. syntax is gator addfeed 'Hacker News' 'https://hackernews.com/rss'")
		os.Exit(1)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	params := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	newCmd := command{
		Name: "follow",
		Args: []string{url},
	}
	fmt.Printf("Feed added successfully: %v\n", feed.Name)

	if err := handlerFollow(s, newCmd, user); err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		return fmt.Errorf("no feeds found")
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %v\nURL: %v\nAdded by: %v\n", feed.Name, feed.Url, feed.Username.String)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Println("Please add URL after follow command")
		os.Exit(1)
	}
	url := cmd.Args[0]

	userId := user.ID

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Printf("No feed found at URL: %v\n", cmd.Args[0])
		os.Exit(1)
	}

	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed successfully followed by %v: %v\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Printf("no feeds for user: %v", user.Name)
	}

	for _, feed := range feeds {
		fmt.Printf("%v\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("please add url after unfollow command")
	}

	if err := s.db.Unfollow(context.Background(), database.UnfollowParams{Name: user.Name, Url: cmd.Args[0]}); err != nil {
		return err
	}

	return nil
}
