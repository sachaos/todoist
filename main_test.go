package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHelpCommand(t *testing.T) {
	tests := []struct {
		name     string
		cliArgs  []string
		osArgs   []string
		expected bool
	}{
		{
			name:     "no args shows help",
			cliArgs:  []string{},
			osArgs:   []string{"todoist"},
			expected: true,
		},
		{
			name:     "help subcommand",
			cliArgs:  []string{"help"},
			osArgs:   []string{"todoist", "help"},
			expected: true,
		},
		{
			name:     "h alias",
			cliArgs:  []string{"h"},
			osArgs:   []string{"todoist", "h"},
			expected: true,
		},
		{
			name:     "--help flag",
			cliArgs:  []string{"list"},
			osArgs:   []string{"todoist", "list", "--help"},
			expected: true,
		},
		{
			name:     "-h flag",
			cliArgs:  []string{"list"},
			osArgs:   []string{"todoist", "list", "-h"},
			expected: true,
		},
		{
			name:     "top-level --help",
			cliArgs:  []string{},
			osArgs:   []string{"todoist", "--help"},
			expected: true,
		},
		{
			name:     "list command is not help",
			cliArgs:  []string{"list"},
			osArgs:   []string{"todoist", "list"},
			expected: false,
		},
		{
			name:     "add command is not help",
			cliArgs:  []string{"add", "buy milk"},
			osArgs:   []string{"todoist", "add", "buy milk"},
			expected: false,
		},
		{
			name:     "sync command is not help",
			cliArgs:  []string{"sync"},
			osArgs:   []string{"todoist", "sync"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHelpCommand(tt.cliArgs, tt.osArgs)
			assert.Equal(t, tt.expected, result)
		})
	}
}
