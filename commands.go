package main

import "github.com/tonyserranodev/task-tracker-cli/internal/store"

type cliCommand struct {
	name        string
	description string
	callback    func(*store.Store, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"add": {
			name:        "add",
			description: "Add a task.",
			callback:    commandAdd,
		},
		"list": {
			name:        "list",
			description: "List all tasks.",
			callback:    commandList,
		},
	}
}
