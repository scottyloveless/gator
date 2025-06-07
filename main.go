package main

import (
	"encoding/json"
	"fmt"

	"github.com/scottyloveless/gator/internal/config"
)

func main() {
	cfgRead, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
	}

	cfgRead.SetUser()

	configAgain, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
	}

	data, err := json.Marshal(configAgain)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	fmt.Printf("%v", string(data))
}
