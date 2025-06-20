package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/scottyloveless/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("no username found. Please add username after login command")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguemnts")
	}

	checkUser, _ := s.db.GetUser(context.Background(), cmd.Args[0])
	if checkUser.Name != cmd.Args[0] {
		fmt.Println("user doesn't exist")
		os.Exit(1)
	}

	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Println("user has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("no name found, please add your username after the register command")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	checkUser, _ := s.db.GetUser(context.Background(), cmd.Args[0])
	if checkUser.Name == cmd.Args[0] {
		fmt.Println("user already exists")
		os.Exit(1)
	}

	newUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("User created successfully!\nname: %v\nID: %v\ncreated_at: %v\n", newUser.Name, newUser.ID, newUser.CreatedAt)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUserLoggedIn := s.cfg.CurrentUserName

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
