package main

import (
	"fmt"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

// commandHelp prints usage information for all commands or a specific command.
func commandHelp(st *store.Store, args ...string) error {
	commands := getCommands()

	if len(args) == 0 {
		fmt.Println("Usage:")
		for _, command := range commands {
			fmt.Printf("%s: %s\n", command.name, command.description)
		}
		return nil
	}

	command := commands[args[0]]
	fmt.Printf("%s: %s\n", command.name, command.description)

	return nil
}
