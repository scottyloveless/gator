package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s.cfg == nil {
		return fmt.Errorf("no state found")
	}

	f, exists := c.commandMap[cmd.Name]
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
