package main

import (
	"fmt"
	"time"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
	"github.com/tonyserranodev/task-tracker-cli/internal/style"
)

// commandList prints all tasks in the store.
func commandList(st *store.Store, _ ...string) error {
	if len(st.Tasks) == 0 {
		fmt.Println("no tasks yet")
		return nil
	}

	blueBold := style.Style{Foreground: style.Blue, Bold: true}
	green := style.Style{Foreground: style.Green}

	// Compute the visible width needed for each column so everything lines up.
	idW := len("ID")
	descW := len("Description")
	statusW := len("Status")
	createdW := len("Created")
	updatedW := len("Updated")

	for _, task := range st.Tasks {
		idW = max(idW, len(fmt.Sprint(task.ID)))
		descW = max(descW, len(task.Description))
		statusW = max(statusW, len(task.Status))
		createdW = max(createdW, len(task.CreatedAt.Format(time.RFC822)))
		updatedW = max(updatedW, len(task.UpdatedAt.Format(time.RFC822)))
	}

	// Helper to join one row's columns padded to their visible widths.
	row := func(cols ...string) string {
		return style.PadRight(cols[0], idW) + " " +
			style.PadRight(cols[1], descW) + " " +
			style.PadRight(cols[2], statusW) + " " +
			style.PadRight(cols[3], createdW) + " " +
			style.PadRight(cols[4], updatedW)
	}

	lines := []string{
		row(blueBold.Apply("ID"), blueBold.Apply("Description"), blueBold.Apply("Status"), blueBold.Apply("Created"), blueBold.Apply("Updated")),
	}

	for _, task := range st.Tasks {
		lines = append(lines, row(
			blueBold.Apply(fmt.Sprint(task.ID)),
			task.Description,
			green.Apply(task.Status),
			task.CreatedAt.Format(time.RFC822),
			task.UpdatedAt.Format(time.RFC822),
		))
	}

	fmt.Println(style.Box(0, lines, style.SingleBorders))

	return nil
}
