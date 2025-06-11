package main

import (
	"fmt"
	"os"

	"github.com/scottyloveless/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	globalConfig, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
	}
	appState := state{
		config: &globalConfig,
	}

	cmdMap := make(map[string]func(*state, command) error)
	cmds := commands{
		commandMap: cmdMap,
	}
	cmds.register("login", handlerLogin)

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

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Println("user has been set")

	return nil
}
