package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("user table reset")
	return nil
}
