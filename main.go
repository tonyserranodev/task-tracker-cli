// task-tracker-cli is a CLI tool for managing a personal task list.
package main

import (
	"fmt"
	"os"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

// main loads saved tasks, parses the CLI command, and executes it.
func main() {
	var st = store.NewStore()

	err := st.LoadTasks()
	if err != nil {
		fmt.Println(err)
	}

	if len(os.Args) <= 1 {
		fmt.Println(`Usage: task-tracker-cli add "Buy Eggs"`)
		os.Exit(0)
	}
	userArgs := os.Args[1:]
	commandName := userArgs[0]
	cmd, ok := getCommands()[commandName]
	if !ok {
		fmt.Println("Unknown command")
		os.Exit(1)
	}
	cmdArgs := userArgs[1:]

	err = cmd.callback(st, cmdArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
