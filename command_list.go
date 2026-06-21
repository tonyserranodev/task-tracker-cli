package main

import (
	"fmt"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

func commandList(st *store.Store, _ ...string) error {
	if len(st.Tasks) == 0 {
		fmt.Println("no tasks yet")
	}

	for _, task := range st.Tasks {
		fmt.Printf("%v. %s %s %v", task.ID, task.Description, task.Status, task.CreatedAt)
	}

	return nil
}
