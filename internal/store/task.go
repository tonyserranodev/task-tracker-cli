package store

import (
	"errors"
	"os"
	"sync/atomic"
	"time"
)

// Status represents the possible states of a task.
type Status int

const (
	// Todo indicates a task that has not been started.
	Todo Status = iota
	// InProgress indicates a task that is currently being worked on.
	InProgress
	// Done indicates a completed task.
	Done
)

// NewStatus converts a status string into a Status value.
func NewStatus(status string) Status {
	switch status {
	case "todo":
		return Todo
	case "in-progress":
		return InProgress
	case "done":
		return Done
	default:
		return Todo
	}
}

// statusName maps each Status to its string representation.
var statusName = map[Status]string{
	Todo:       "todo",
	InProgress: "in progress",
	Done:       "done",
}

// String returns the string representation of a task status.
func (d Status) String() string {
	return statusName[d]
}

// Task represents a single task in the task tracker.
type Task struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// taskID holds the counter used to generate unique task IDs.
var taskID atomic.Int64

// ResetTaskID resets the task ID counter to zero. It is intended for tests.
func ResetTaskID() {
	taskID.Store(0)
}

// NewTask creates a new Task with the given description and a unique ID.
func NewTask(desc string) Task {

	newTask := Task{
		ID:          taskID.Add(1),
		Description: desc,
		Status:      statusName[Todo],
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return newTask
}

// LoadTasks reads tasks from disk into the store if the file exists.
func (s *Store) LoadTasks() error {
	if !fileExists("tasks.json") {
		return nil
	}
	loadedTasks, err := s.GetAll()
	if err != nil {
		return err
	}
	s.Tasks = loadedTasks

	return nil
}

// fileExists reports whether the named file exists.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}
