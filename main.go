package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/scottyloveless/gator/internal/config"
	"github.com/scottyloveless/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
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
		db:     dbQueries,
		config: &globalConfig,
	}

	cmdMap := make(map[string]func(*state, command) error)
	cmds := commands{
		commandMap: cmdMap,
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	var newCommand command
	newCommand.name = args[1]
	newCommand.args = args[2:]

	if err := cmds.run(&appState, newCommand); err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s.config == nil {
		return fmt.Errorf("no state found")
	}

	f, exists := c.commandMap[cmd.name]
	if !exists {
		return fmt.Errorf("command does not exist")
	} else {
		if err := f(s, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) <= 0 {
		return fmt.Errorf("no username found. Please add username after login command")
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("too many arguemnts")
	}

	checkUser, _ := s.db.CheckUser(context.Background(), cmd.args[0])
	if checkUser.Name != cmd.args[0] {
		fmt.Println("user doesn't exist")
		os.Exit(1)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Println("user has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) <= 0 {
		return fmt.Errorf("no name found, please add your username after the register command")
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	checkUser, _ := s.db.CheckUser(context.Background(), cmd.args[0])
	if checkUser.Name == cmd.args[0] {
		fmt.Println("user already exists")
		os.Exit(1)
	}

	newUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("User created successfully!\nname: %v\nID: %v\ncreated_at: %v\n", newUser.Name, newUser.ID, newUser.CreatedAt)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("user table reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUserLoggedIn := s.config.CurrentUsername

	if len(users) == 0 {
		return fmt.Errorf("no users found")
	}

	for _, user := range users {
		if user.Name == currentUserLoggedIn {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}
