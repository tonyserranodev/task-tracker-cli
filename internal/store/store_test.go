package store

import (
	"os"
	"testing"
)

func TestStoreAdd(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	tt := map[string]struct {
		desc    string
		wantErr bool
	}{
		"valid task":   {"Buy milk", false},
		"another task": {"Walk dog", false},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := NewStore()
			task := NewTask(tc.desc)

			err := s.Add(task)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Add() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(s.Tasks) != 1 {
				t.Errorf("len(Tasks) = %d, want 1", len(s.Tasks))
			}
			if s.Tasks[0].Description != tc.desc {
				t.Errorf("Description = %q, want %q", s.Tasks[0].Description, tc.desc)
			}
		})
	}
}

func TestStoreUpdateDescription(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	s := NewStore()
	_ = s.Add(NewTask("Old desc"))

	tt := map[string]struct {
		id          int64
		description string
		wantErr     bool
	}{
		"existing": {1, "New desc", false},
		"missing":  {99, "Anything", true},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := s.UpdateDescription(tc.id, tc.description)
			if (err != nil) != tc.wantErr {
				t.Fatalf("UpdateDescription() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				got, _ := s.GetByID(tc.id)
				if got.Description != tc.description {
					t.Errorf("Description = %q, want %q", got.Description, tc.description)
				}
			}
		})
	}
}

func TestStoreUpdateStatus(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	s := NewStore()
	_ = s.Add(NewTask("Task"))

	tt := map[string]struct {
		id      int64
		status  Status
		wantErr bool
	}{
		"existing todo": {1, Todo, false},
		"existing done": {1, Done, false},
		"missing":       {99, InProgress, true},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := s.UpdateStatus(tc.id, tc.status)
			if (err != nil) != tc.wantErr {
				t.Fatalf("UpdateStatus() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				got, _ := s.GetByID(tc.id)
				want := statusName[tc.status]
				if got.Status != want {
					t.Errorf("Status = %q, want %q", got.Status, want)
				}
			}
		})
	}
}

func TestStoreGetByID(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	s := NewStore()
	_ = s.Add(NewTask("One"))
	_ = s.Add(NewTask("Two"))

	tt := map[string]struct {
		id      int64
		want    string
		wantErr bool
	}{
		"first":   {1, "One", false},
		"second":  {2, "Two", false},
		"missing": {99, "", true},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := s.GetByID(tc.id)
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetByID() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got.Description != tc.want {
				t.Errorf("Description = %q, want %q", got.Description, tc.want)
			}
		})
	}
}

func TestStoreDelete(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	s := NewStore()
	_ = s.Add(NewTask("One"))
	_ = s.Add(NewTask("Two"))

	tt := []struct {
		name       string
		id         int64
		wantErr    bool
		wantLen    int
		wantRemain string
	}{
		{"existing first", 1, false, 1, "Two"},
		{"missing", 99, true, 1, "Two"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Delete(tc.id)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tc.wantErr)
			}
			if len(s.Tasks) != tc.wantLen {
				t.Errorf("len(Tasks) = %d, want %d", len(s.Tasks), tc.wantLen)
			}
			if len(s.Tasks) > 0 && s.Tasks[0].Description != tc.wantRemain {
				t.Errorf("remaining Description = %q, want %q", s.Tasks[0].Description, tc.wantRemain)
			}
		})
	}
}

func TestStoreSaveAndGetAll(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	s := NewStore()
	_ = s.Add(NewTask("Saved"))

	loaded, err := s.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("len(loaded) = %d, want 1", len(loaded))
	}
	if loaded[0].Description != "Saved" {
		t.Errorf("Description = %q, want %q", loaded[0].Description, "Saved")
	}
}

func TestStoreLoadTasks(t *testing.T) {
	t.Chdir(t.TempDir())
	taskID.Store(0)

	// Create a store, add a task, and save it.
	s := NewStore()
	_ = s.Add(NewTask("Persisted"))

	// Start fresh and load from disk.
	s2 := NewStore()
	if err := s2.LoadTasks(); err != nil {
		t.Fatalf("LoadTasks() error = %v", err)
	}

	if len(s2.Tasks) != 1 {
		t.Fatalf("len(Tasks) = %d, want 1", len(s2.Tasks))
	}
	if s2.Tasks[0].Description != "Persisted" {
		t.Errorf("Description = %q, want %q", s2.Tasks[0].Description, "Persisted")
	}
}

func TestStoreLoadTasksNoFile(t *testing.T) {
	t.Chdir(t.TempDir())

	s := NewStore()
	if err := s.LoadTasks(); err != nil {
		t.Fatalf("LoadTasks() error = %v", err)
	}
	if len(s.Tasks) != 0 {
		t.Errorf("len(Tasks) = %d, want 0", len(s.Tasks))
	}
}

func TestUpdateCounter(t *testing.T) {
	tt := map[string]struct {
		tasks []Task
		want  int64
	}{
		"empty":    {[]Task{}, 0},
		"single":   {[]Task{{ID: 5}}, 5},
		"multiple": {[]Task{{ID: 1}, {ID: 10}, {ID: 3}}, 10},
		"zero id":  {[]Task{{ID: 0}}, 0},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			taskID.Store(0)
			updateCounter(tc.tasks)
			if got := taskID.Load(); got != tc.want {
				t.Errorf("taskID = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tt := map[string]struct {
		setup func(t *testing.T) string
		want  bool
	}{
		"exists": {
			setup: func(t *testing.T) string {
				f, err := os.CreateTemp("", "tasks-*.json")
				if err != nil {
					t.Fatal(err)
				}
				f.Close()
				return f.Name()
			},
			want: true,
		},
		"missing": {
			setup: func(t *testing.T) string {
				return "definitely-does-not-exist-12345.json"
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)
			got := fileExists(path)
			if got != tc.want {
				t.Errorf("fileExists(%q) = %v, want %v", path, got, tc.want)
			}
		})
	}
}
