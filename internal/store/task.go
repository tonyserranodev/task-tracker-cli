package store

import "time"

type Description int

const (
	Todo Description = iota
	InProgress
	Done
)

var descriptionName = map[Description]string{
	Todo:       "todo",
	InProgress: "in progress",
	Done:       "done",
}

func (d Description) String() string {
	return descriptionName[d]
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
