package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
	"github.com/tonyserranodev/task-tracker-cli/internal/style"
)

// commandMark updates the status of the task with the given ID.
func commandMark(st *store.Store, args ...string) error {
	if len(args) != 2 {
		return errors.New("must provide the id of a task to mark and a new status")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	status := store.NewStatus(args[1])

	err = st.UpdateStatus(int64(id), status)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Task with id %v has been marked %v!\n", id, status)
	fmt.Println(style.Style{Foreground: style.Green, Bold: true}.Apply(msg))

	return nil
}
