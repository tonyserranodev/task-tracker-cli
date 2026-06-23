# Task Tracker CLI

A small command-line tool for managing a personal task list. Tasks are stored locally in a `tasks.json` file.

## Build

Requires Go 1.25 or later.

```bash
go build -o task-tracker-cli
```

This produces an executable named `task-tracker-cli` in the current directory.

## Usage

```bash
./task-tracker-cli <command> [args...]
```

Run without arguments to see a short usage message.

## Commands

### Add a task

```bash
./task-tracker-cli add "Buy groceries"
```

Creates a new task with status `todo` and prints the assigned task ID.

### List all tasks

```bash
./task-tracker-cli list
```

Displays all tasks in a bordered table with ID, description, status, creation time, and last update time.

### Update a task description

```bash
task-tracker-cli update 1 "Buy groceries and milk"
```

Changes the description of the task with ID `1`.

### Mark a task status

```bash
./task-tracker-cli mark 1 done
```

Valid statuses are:

- `todo`
- `in-progress`
- `done`

### Delete a task

```bash
./task-tracker-cli delete 1
```

Removes the task with ID `1`.

### Show help

```bash
./task-tracker-cli help
```

Prints a list of all commands and their descriptions.

To see help for a specific command:

```bash
./task-tracker-cli help add
```

## Data storage

Tasks are persisted to a `tasks.json` file in the current working directory. The file is created automatically when a task is added, updated, marked, or deleted. Keep this file in the directory where you run the tool if you want to keep your task history.

## Development

Run the tests:

```bash
go test ./...
```

Run the linter:

```bash
go vet ./...
```

Build the binary:

```bash
go build -o task-tracker-cli
```
