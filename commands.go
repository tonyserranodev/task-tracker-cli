package main

import "github.com/tonyserranodev/task-tracker-cli/internal/store"

// cliCommand describes a single CLI command and its handler.
type cliCommand struct {
	name        string
	usage       string
	description string
	callback    func(*store.Store, ...string) error
}

// getCommands returns the map of available CLI commands.
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"add": {
			name:        "add",
			usage:       "add <description>",
			description: "Add a task.",
			callback:    commandAdd,
		},
		"list": {
			name:        "list",
			usage:       "list",
			description: "List all tasks.",
			callback:    commandList,
		},
		"delete": {
			name:        "delete",
			usage:       "delete <id>",
			description: "Delete a task by its ID.",
			callback:    commandDelete,
		},
		"update": {
			name:        "update",
			usage:       "update <id> <description>",
			description: "Update a task description by its ID.",
			callback:    commandUpdate,
		},
		"mark": {
			name:        "mark",
			usage:       "mark <id> <status>",
			description: "Update a task status by its ID.",
			callback:    commandMark,
		},
		"help": {
			name:        "help",
			usage:       "help [command]",
			description: "Print descriptions for all commands.",
			callback:    commandHelp,
		},
	}
}
