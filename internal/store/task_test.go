package store

import (
	"testing"
)

func TestNewStatus(t *testing.T) {
	tt := map[string]struct {
		input string
		want  Status
	}{
		"todo":        {"todo", Todo},
		"in-progress": {"in-progress", InProgress},
		"done":        {"done", Done},
		"unknown":     {"unknown", Todo},
		"empty":       {"", Todo},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := NewStatus(tc.input)
			if got != tc.want {
				t.Errorf("NewStatus(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestStatusString(t *testing.T) {
	tt := map[string]struct {
		status Status
		want   string
	}{
		"todo":        {Todo, "todo"},
		"in-progress": {InProgress, "in-progress"},
		"done":        {Done, "done"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := tc.status.String()
			if got != tc.want {
				t.Errorf("%v.String() = %q, want %q", tc.status, got, tc.want)
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	original := taskID.Load()
	taskID.Store(0)
	defer taskID.Store(original)

	tt := []struct {
		name       string
		desc       string
		wantID     int64
		wantStatus string
	}{
		{"first task", "Buy milk", 1, "todo"},
		{"second task", "Walk dog", 2, "todo"},
		{"empty desc", "", 3, "todo"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := NewTask(tc.desc)

			if got.ID != tc.wantID {
				t.Errorf("ID = %d, want %d", got.ID, tc.wantID)
			}
			if got.Description != tc.desc {
				t.Errorf("Description = %q, want %q", got.Description, tc.desc)
			}
			if got.Status != tc.wantStatus {
				t.Errorf("Status = %q, want %q", got.Status, tc.wantStatus)
			}
			if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
				t.Errorf("CreatedAt/UpdatedAt should be set")
			}
		})
	}
}
