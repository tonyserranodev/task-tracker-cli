package main

import (
	"strings"
	"testing"
	"time"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

// newTestTask returns a Task with a fixed ID and timestamps for use in tests.
func newTestTask(id int64, description string) store.Task {
	fixed := time.Date(2024, time.June, 1, 12, 0, 0, 0, time.UTC)
	return store.Task{
		ID:          id,
		Description: description,
		Status:      "todo",
		CreatedAt:   fixed,
		UpdatedAt:   fixed,
	}
}

func TestGetCommands(t *testing.T) {
	want := []string{"add", "list", "delete", "update", "mark", "help"}

	cmds := getCommands()
	if len(cmds) != len(want) {
		t.Fatalf("len(commands) = %d, want %d", len(cmds), len(want))
	}

	for _, name := range want {
		cmd, ok := cmds[name]
		if !ok {
			t.Errorf("missing command %q", name)
			continue
		}
		if cmd.name != name {
			t.Errorf("command %q has name %q", name, cmd.name)
		}
		if cmd.usage == "" {
			t.Errorf("command %q has no usage", name)
		}
		if cmd.callback == nil {
			t.Errorf("command %q has no callback", name)
		}
	}
}

func TestCommandAdd(t *testing.T) {
	t.Chdir(t.TempDir())
	store.ResetTaskID()
	defer store.ResetTaskID()

	tt := map[string]struct {
		args    []string
		wantErr bool
		wantLen int
		wantID  int64
	}{
		"with description": {[]string{"Buy milk"}, false, 1, 1},
		"missing arg":      {[]string{}, true, 0, 0},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			err := commandAdd(s, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandAdd() error = %v, wantErr %v", err, tc.wantErr)
			}
			if len(s.Tasks) != tc.wantLen {
				t.Errorf("len(Tasks) = %d, want %d", len(s.Tasks), tc.wantLen)
			}
			if tc.wantID != 0 && s.Tasks[0].ID != tc.wantID {
				t.Errorf("ID = %d, want %d", s.Tasks[0].ID, tc.wantID)
			}
		})
	}
}

func TestCommandDelete(t *testing.T) {
	t.Chdir(t.TempDir())

	tt := map[string]struct {
		args    []string
		wantErr bool
		wantLen int
	}{
		"existing":    {[]string{"1"}, false, 0},
		"missing arg": {[]string{}, true, 1},
		"bad id":      {[]string{"not-a-number"}, true, 1},
		"not found":   {[]string{"99"}, true, 1},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			s.Tasks = []store.Task{newTestTask(1, "Task")}

			err := commandDelete(s, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandDelete() error = %v, wantErr %v", err, tc.wantErr)
			}
			if len(s.Tasks) != tc.wantLen {
				t.Errorf("len(Tasks) = %d, want %d", len(s.Tasks), tc.wantLen)
			}
		})
	}
}

func TestCommandUpdate(t *testing.T) {
	t.Chdir(t.TempDir())

	tt := []struct {
		name    string
		args    []string
		wantErr bool
		want    string
	}{
		{"existing", []string{"1", "New desc"}, false, "New desc"},
		{"missing args", []string{"1"}, true, "Old desc"},
		{"bad id", []string{"abc", "New desc"}, true, "Old desc"},
		{"not found", []string{"99", "New desc"}, true, "Old desc"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := store.NewStore()
			s.Tasks = []store.Task{newTestTask(1, "Old desc")}

			err := commandUpdate(s, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandUpdate() error = %v, wantErr %v", err, tc.wantErr)
			}

			got, _ := s.GetByID(1)
			if got.Description != tc.want {
				t.Errorf("Description = %q, want %q", got.Description, tc.want)
			}
		})
	}
}

func TestCommandMark(t *testing.T) {
	t.Chdir(t.TempDir())

	tt := []struct {
		name    string
		args    []string
		wantErr bool
		want    string
	}{
		{"existing", []string{"1", "done"}, false, "done"},
		{"missing args", []string{"1"}, true, "todo"},
		{"bad id", []string{"abc", "done"}, true, "todo"},
		{"not found", []string{"99", "done"}, true, "todo"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := store.NewStore()
			s.Tasks = []store.Task{newTestTask(1, "Task")}

			err := commandMark(s, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandMark() error = %v, wantErr %v", err, tc.wantErr)
			}

			got, _ := s.GetByID(1)
			if got.Status != tc.want {
				t.Errorf("Status = %q, want %q", got.Status, tc.want)
			}
		})
	}
}

func TestCommandList(t *testing.T) {
	t.Chdir(t.TempDir())

	tt := map[string]struct {
		setup   func(*store.Store)
		wantErr bool
	}{
		"empty": {
			setup:   func(s *store.Store) {},
			wantErr: false,
		},
		"with tasks": {
			setup: func(s *store.Store) {
				s.Tasks = []store.Task{
					newTestTask(1, "One"),
					newTestTask(2, "Two"),
				}
			},
			wantErr: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			tc.setup(s)

			err := commandList(s)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandList() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestCommandHelp(t *testing.T) {
	tt := map[string]struct {
		args    []string
		wantErr bool
		want    []string
	}{
		"all commands": {
			args:    []string{},
			wantErr: false,
			want:    []string{"add:", "list:", "delete:", "help:"},
		},
		"specific command": {
			args:    []string{"add"},
			wantErr: false,
			want:    []string{"add:", "Add a task."},
		},
		"unknown command": {
			args:    []string{"nope"},
			wantErr: true,
			want:    []string{},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			err := commandHelp(s, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandHelp() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestFormatHelpLine(t *testing.T) {
	got, _ := formatHelpLine("add <description>", "Add a task.", 20)
	if !strings.Contains(got, "add <description>") {
		t.Errorf("formatHelpLine() missing usage")
	}
	if !strings.Contains(got, "Add a task.") {
		t.Errorf("formatHelpLine() missing description")
	}
}

func TestMaxUsageLength(t *testing.T) {
	got := maxUsageLength(getCommands())
	if got <= 0 {
		t.Errorf("maxUsageLength() = %d, want > 0", got)
	}
}

func TestFormatTaskTable(t *testing.T) {
	fixed := time.Date(2024, time.June, 1, 12, 0, 0, 0, time.UTC)

	tt := map[string]struct {
		tasks []store.Task
		want  []string
	}{
		"single task": {
			tasks: []store.Task{
				{ID: 1, Description: "Buy milk", Status: "todo", CreatedAt: fixed, UpdatedAt: fixed},
			},
			want: []string{"ID", "Description", "Buy milk", "todo"},
		},
		"multiple tasks": {
			tasks: []store.Task{
				{ID: 1, Description: "One", Status: "todo", CreatedAt: fixed, UpdatedAt: fixed},
				{ID: 2, Description: "Two", Status: "done", CreatedAt: fixed, UpdatedAt: fixed},
			},
			want: []string{"One", "Two", "done"},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, _ := formatTaskTable(tc.tasks)
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("formatTaskTable() output missing %q", w)
				}
			}
		})
	}
}
