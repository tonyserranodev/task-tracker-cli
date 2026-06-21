package main

import (
	"fmt"
	"os"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

func main() {
	var st = store.NewStore()

	err := st.LoadTasks()
	if err != nil {
		fmt.Println(err)
	}

	commands := getCommands()
	if len(os.Args) <= 1 {
		fmt.Println(`Usage: task-tracker-cli add "Buy Eggs"`)
		os.Exit(0)
	}
	userArgs := os.Args[1:]
	commandName := userArgs[0]
	cmd, ok := commands[commandName]
	if !ok {
		fmt.Println("Unknown command")
	}
	cmdArgs := userArgs[1:]

	err = cmd.callback(st, cmdArgs...)
	if err != nil {
		fmt.Println(err)
	}
}
