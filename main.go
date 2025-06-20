package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/scottyloveless/gator/internal/config"
	"github.com/scottyloveless/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	globalConfig, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
	}

	db, err := sql.Open("postgres", globalConfig.DBurl)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
	}

	dbQueries := database.New(db)

	appState := state{
		db:  dbQueries,
		cfg: &globalConfig,
	}

	cmdMap := make(map[string]func(*state, command) error)
	cmds := commands{
		commandMap: cmdMap,
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	var newCommand command
	newCommand.Name = args[1]
	newCommand.Args = args[2:]

	if err := cmds.run(&appState, newCommand); err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		if err := handler(s, cmd, user); err != nil {
			return err
		}

		return nil
	}
}
