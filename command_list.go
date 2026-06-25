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
		msg, err := style.Render("no tasks yet", "yellow")
		if err != nil {
			return err
		}

		fmt.Println(msg)

		return nil
	}

	taskTable, err := formatTaskTable(st.Tasks)
	if err != nil {
		return err
	}

	fmt.Println(taskTable)

	return nil
}

// formatTaskTable returns a styled, boxed table of tasks.
func formatTaskTable(tasks []store.Task) (string, error) {
	idW := len("ID")
	descW := len("Description")
	statusW := len("Status")
	createdW := len("Created")
	updatedW := len("Updated")

	// Use the longest value in each column as the column width.
	for _, task := range tasks {
		idW = max(idW, len(fmt.Sprint(task.ID)))
		descW = max(descW, len(task.Description))
		statusW = max(statusW, len(task.Status))
		createdW = max(createdW, len(task.CreatedAt.Format(time.RFC822)))
		updatedW = max(updatedW, len(task.UpdatedAt.Format(time.RFC822)))
	}

	// Helper function to format columns
	row := func(cols ...string) string {
		return style.PadRight(cols[0], idW) + " " +
			style.PadRight(cols[1], descW) + " " +
			style.PadRight(cols[2], statusW) + " " +
			style.PadRight(cols[3], createdW) + " " +
			style.PadRight(cols[4], updatedW)
	}

	headers := []string{"ID", "Description", "Status", "Created", "Updated"}
	styledHeaders := make([]string, len(headers))

	for i, h := range headers {
		styled, err := style.Render(h, "blue", "bold")
		if err != nil {
			return "", err
		}
		styledHeaders[i] = styled
	}

	lines := []string{
		row(styledHeaders...),
	}

	for _, task := range tasks {
		styledID, err := style.Render(fmt.Sprint(task.ID), "blue")
		if err != nil {
			return "", err
		}

		styledStatus, err := style.Render(task.Status, "green")
		if err != nil {
			return "", err
		}

		lines = append(lines,
			row(
				styledID,
				task.Description,
				styledStatus,
				task.CreatedAt.Format(time.RFC822),
				task.UpdatedAt.Format(time.RFC822),
			))
	}

	return style.Box(0, lines, style.SingleBorders), nil
}
