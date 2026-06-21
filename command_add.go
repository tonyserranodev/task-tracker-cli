package main

import (
	"fmt"
	"os"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

func commandAdd(st *store.Store, args ...string) error {
	if len(args) < 1 {
		fmt.Println("Must provide a description of a task to add.")
		os.Exit(1)
	}

	taskDesc := args[0]
	task := store.NewTask(taskDesc)
	task.Status = store.Todo.String()
	err := st.Add(task)
	if err != nil {
		return err
	}
	fmt.Println("Task added succedfully!")

	return nil
}
