package main

import (
	"fmt"
	"sort"

	"github.com/tonyserranodev/task-tracker-cli/internal/store"
	"github.com/tonyserranodev/task-tracker-cli/internal/style"
)

// commandHelp prints usage information for all commands or a specific command.
func commandHelp(st *store.Store, args ...string) error {
	commands := getCommands()

	if len(args) == 0 {
		names := sortedCommandNames(commands)
		maxUsageLen := maxUsageLength(commands)
		lines := make([]string, 0, len(commands))
		for _, name := range names {
			cmd := commands[name]
			lines = append(lines, formatHelpLine(cmd.usage, cmd.description, maxUsageLen))
		}
		fmt.Println(style.Box(0, lines, style.SingleBorders))
		return nil
	}

	command, ok := commands[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %s", args[0])
	}

	fmt.Println(style.Box(0, []string{formatHelpLine(command.usage, command.description, len(command.usage))}, style.SingleBorders))
	return nil
}

// sortedCommandNames returns the command names in alphabetical order.
func sortedCommandNames(commands map[string]cliCommand) []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// maxUsageLength returns the visible length of the longest command usage string.
func maxUsageLength(commands map[string]cliCommand) int {
	max := 0
	for _, cmd := range commands {
		if l := len(cmd.usage); l > max {
			max = l
		}
	}
	return max
}

// formatHelpLine returns a styled "usage: description" help entry.
// usageWidth is the length of the longest usage string; it is used to pad
// "usage:" so the descriptions align across lines.
func formatHelpLine(usage, description string, usageWidth int) string {
	usageStyle := style.Style{Foreground: style.Cyan, Bold: true}
	padded := style.PadRight(usageStyle.Apply(usage+":"), usageWidth+1)
	return padded + " " + description
}
