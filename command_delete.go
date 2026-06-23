package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

// commandDelete removes the task with the given ID from the store.
func commandDelete(st *store.Store, args ...string) error {

	if len(args) < 1 {
		return errors.New("must provide the id of a task to delete")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if err := st.Delete(int64(id)); err != nil {
		return err
	}
	fmt.Printf("Task %v deleted!\n", id)
	return nil
}
