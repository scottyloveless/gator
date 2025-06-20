package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scottyloveless/gator/internal/database"
)

func browse(s *state, limit ...int) {
	limitActual := 2
	if len(limit) > 0 {
		limitActual = limit[0]
	}
	p := tea.NewProgram(initialModel(s, int32(limitActual)), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	// posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{Name: s.cfg.CurrentUserName, Limit: int32(limitActual)})
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	//
	// if len(posts) == 0 {
	// 	log.Printf("no posts found")
	// }
	//
	// for _, post := range posts {
	// 	printPost(post)
	// }
}

func handlerBrowse(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 1 {
		return fmt.Errorf("format is browse 10")
	}

	intify, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("enter valid integer after browse")
	}

	browse(s, intify)

	return nil
}

func printPost(post database.Post) {
	fmt.Printf(" * Title:          %v\n", post.Title)
	fmt.Printf(" * URL:            %v\n", post.Url)
	fmt.Printf(" * Description:    %v\n", post.Description)
}
