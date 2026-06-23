package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

// commandUpdate changes the description of the task with the given ID.
func commandUpdate(st *store.Store, args ...string) error {

	if len(args) != 2 {
		return errors.New("must provide the id of a task to update and a new description")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	newDescription := args[1]

	err = st.UpdateDescription(int64(id), newDescription)
	if err != nil {
		return err
	}

	fmt.Printf("Task with id %v has been updated!\n", id)
	return nil
}
