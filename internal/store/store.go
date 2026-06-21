package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

type Store struct {
	Tasks []Task
}

func NewStore() *Store {

	return &Store{
		Tasks: []Task{},
	}
}

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

var idCounter atomic.Int64

func NewTask(desc string) Task {
	return Task{
		ID:          int(idCounter.Add(1)),
		Description: desc,
		CreatedAt:   time.Now(),
	}
}

func (s *Store) Add(task Task) error {
	s.Tasks = append(s.Tasks, task)
	file, err := os.Create("tasks.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.Tasks)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAll() ([]Task, error) {
	file, err := os.Open("tasks.json")
	if err != nil {
		fmt.Println("Error opening file")
		return []Task{}, err
	}
	defer file.Close()

	var tasks []Task

	dec := json.NewDecoder(file)
	err = dec.Decode(&tasks)
	if err != nil {
		return []Task{}, err
	}

	return tasks, nil
}

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
