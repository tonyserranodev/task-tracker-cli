package main

import (
	"errors"
	"fmt"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

// commandAdd creates a new task from args and stores it.
func commandAdd(st *store.Store, args ...string) error {
	if len(args) < 1 {
		return errors.New("must provide a description of a task to add")
	}

	description := args[0]

	task := store.NewTask(description)
	task.Status = store.Todo.String()

	err := st.Add(task)
	if err != nil {
		return err
	}

	fmt.Println("Task added successfully!")

	return nil
}
