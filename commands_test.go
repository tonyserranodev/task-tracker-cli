package main

import (
	"fmt"
	"testing"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
)

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
		if cmd.callback == nil {
			t.Errorf("command %q has no callback", name)
		}
	}
}

func TestCommandAdd(t *testing.T) {
	t.Chdir(t.TempDir())

	tt := map[string]struct {
		args    []string
		wantErr bool
		wantLen int
	}{
		"with description": {[]string{"Buy milk"}, false, 1},
		"missing arg":      {[]string{}, true, 0},
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
		})
	}
}

func TestCommandDelete(t *testing.T) {
	t.Chdir(t.TempDir())

	for name, tc := range map[string]struct {
		args    func(id int64) []string
		wantErr bool
		wantLen int
	}{
		"existing":    {func(id int64) []string { return []string{fmt.Sprint(id)} }, false, 0},
		"missing arg": {func(id int64) []string { return []string{} }, true, 1},
		"bad id":      {func(id int64) []string { return []string{"not-a-number"} }, true, 1},
		"not found":   {func(id int64) []string { return []string{"99"} }, true, 1},
	} {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			task := store.NewTask("Task")
			_ = s.Add(task)

			err := commandDelete(s, tc.args(task.ID)...)
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

	for name, tc := range map[string]struct {
		args    func(id int64) []string
		wantErr bool
		want    func(desc string) string
	}{
		"existing":     {func(id int64) []string { return []string{fmt.Sprint(id), "New desc"} }, false, func(string) string { return "New desc" }},
		"missing args": {func(id int64) []string { return []string{fmt.Sprint(id)} }, true, func(desc string) string { return desc }},
		"bad id":       {func(id int64) []string { return []string{"abc", "New desc"} }, true, func(desc string) string { return desc }},
		"not found":    {func(id int64) []string { return []string{"99", "New desc"} }, true, func(desc string) string { return desc }},
	} {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			task := store.NewTask("Old desc")
			_ = s.Add(task)

			err := commandUpdate(s, tc.args(task.ID)...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandUpdate() error = %v, wantErr %v", err, tc.wantErr)
			}

			got, _ := s.GetByID(task.ID)
			want := tc.want("Old desc")
			if got.Description != want {
				t.Errorf("Description = %q, want %q", got.Description, want)
			}
		})
	}
}

func TestCommandMark(t *testing.T) {
	t.Chdir(t.TempDir())

	for name, tc := range map[string]struct {
		args    func(id int64) []string
		wantErr bool
		want    string
	}{
		"existing":     {func(id int64) []string { return []string{fmt.Sprint(id), "done"} }, false, "done"},
		"missing args": {func(id int64) []string { return []string{fmt.Sprint(id)} }, true, "todo"},
		"bad id":       {func(id int64) []string { return []string{"abc", "done"} }, true, "todo"},
		"not found":    {func(id int64) []string { return []string{"99", "done"} }, true, "todo"},
	} {
		t.Run(name, func(t *testing.T) {
			s := store.NewStore()
			task := store.NewTask("Task")
			_ = s.Add(task)

			err := commandMark(s, tc.args(task.ID)...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("commandMark() error = %v, wantErr %v", err, tc.wantErr)
			}

			got, _ := s.GetByID(task.ID)
			if got.Status != tc.want {
				t.Errorf("Status = %q, want %q", got.Status, tc.want)
			}
		})
	}
}

func TestCommandList(t *testing.T) {
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
				_ = s.Add(store.NewTask("One"))
				_ = s.Add(store.NewTask("Two"))
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
