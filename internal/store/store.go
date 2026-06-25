// Package store handles persistent storage and retrieval of tasks.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Store holds the in-memory task list and persists it to disk.
type Store struct {
	Tasks []Task
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{
		Tasks: []Task{},
	}
}

// FILEPATH is the default file used to persist tasks.
const FILEPATH = `tasks.json`

// Add appends a task to the store and writes the task list to disk.
func (s *Store) Add(task Task) error {
	s.Tasks = append(s.Tasks, task)
	err := s.SaveTasks()
	if err != nil {
		return err
	}
	return nil
}

// UpdateDescription changes the description of the task with the given ID.
func (s *Store) UpdateDescription(id int64, newDescription string) error {
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks[i].Description = newDescription
			s.Tasks[i].UpdatedAt = time.Now()

			if err := s.SaveTasks(); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("no task found with id, %v", id)
}

// UpdateStatus changes the status of the task with the given ID.
func (s *Store) UpdateStatus(id int64, newStatus Status) error {
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks[i].Status = statusNames[newStatus]
			s.Tasks[i].UpdatedAt = time.Now()

			if err := s.SaveTasks(); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("no task found with id, %v", id)
}

// GetByID returns the task with the given ID, or an error if not found.
func (s *Store) GetByID(id int64) (Task, error) {
	for _, t := range s.Tasks {
		if t.ID == id {
			return t, nil
		}
	}

	return Task{}, fmt.Errorf("no task found with id, %v", id)
}

func (s *Store) GetTasksByStatus(statusString string) ([]Task, error) {
	filtered := []Task{}

	for _, task := range s.Tasks {
		if task.Status == statusString {

			filtered = append(filtered, task)
		}
	}

	if len(filtered) == 0 {

		return nil, fmt.Errorf("no task found with status: %s", statusString)
	}

	return filtered, nil
}

// GetAll reads all tasks from disk and updates the global ID counter.
// TODO: factor out the file loading logic into a Load function
// GetAll should just return the tasks from memory
func (s *Store) GetAll() ([]Task, error) {
	file, err := os.Open(FILEPATH)
	if err != nil {
		return []Task{}, err
	}
	defer file.Close()

	var tasks []Task

	dec := json.NewDecoder(file)
	err = dec.Decode(&tasks)
	if err != nil {
		return []Task{}, err
	}
	updateCounter(tasks)

	return tasks, nil
}

// updateCounter sets the task ID counter to the highest existing task ID.
func updateCounter(tasks []Task) {
	var max int64
	for _, task := range tasks {
		if task.ID > max {
			max = task.ID
		}
	}
	taskID.Store(max)
}

// Delete removes the task with the given ID from the store.
func (s *Store) Delete(id int64) error {
	for i, task := range s.Tasks {
		if task.ID == id {
			newTasks := append(s.Tasks[:i], s.Tasks[i+1:]...)
			s.Tasks = newTasks
			if err := s.SaveTasks(); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("no task with id %v", id)
}

// SaveTasks writes the current in-memory task list to disk.
func (s *Store) SaveTasks() error {

	file, err := os.Create(FILEPATH)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.Tasks)
	if err != nil {
		return err
	}

	return nil
}
