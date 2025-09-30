package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		input          string
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "Help command",
			args: []string{"--help"},
			expectedOutput: []string{
				"LogTail is a powerful command-line tool",
				"Usage:",
				"logtail [file...]",
			},
		},
		{
			name: "Version information in help",
			args: []string{"-h"},
			expectedOutput: []string{
				"filter",
				"color",
				"line-numbers",
			},
		},
		{
			name:        "Invalid flag",
			args:        []string{"--invalid-flag"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test to avoid state pollution
			cmd := &cobra.Command{
				Use:   "logtail [file...]",
				Short: "An intelligent log analyzer for developers",
				Long: `LogTail is a powerful command-line tool for parsing, filtering, and analyzing log files.
It provides real-time filtering, syntax highlighting, and pattern detection.`,
				RunE: runLogTail,
			}

			cmd.Flags().StringVarP(&filterPattern, "filter", "f", "", "Filter logs with regex pattern")
			cmd.Flags().BoolVarP(&colorOutput, "color", "c", true, "Enable colorized output")
			cmd.Flags().BoolVarP(&followMode, "follow", "F", false, "Follow log file like tail -f")
			cmd.Flags().BoolVarP(&showLineNum, "line-numbers", "n", false, "Show line numbers")

			cmd.SetArgs(tt.args)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := cmd.Execute()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			output := buf.String()
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but got: %s", expected, output)
				}
			}
		})
	}
}

// Simple integration test
func TestRunLogTailIntegration(t *testing.T) {
	// Test that the function executes without panic
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.log")

	testContent := `2024-09-30T10:30:45.123Z INFO Application started
2024-09-30T10:30:46.456Z ERROR Database connection failed`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test that it doesn't crash with a valid file
	// We just test that the function can execute, not the exact output
	filterPattern = ""
	colorOutput = false
	showLineNum = false

	// Simulate command execution
	err = runLogTail(nil, []string{testFile})
	if err != nil {
		t.Errorf("runLogTail should not return error for valid file: %v", err)
	}
}

func TestRunLogTailErrors(t *testing.T) {
	// Create a temp file for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.log")

	err := os.WriteFile(testFile, []byte("test log line"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test invalid regex with existing file
	filterPattern = "[invalid"
	colorOutput = true
	showLineNum = false

	err = runLogTail(nil, []string{testFile})
	if err == nil {
		t.Error("Expected error for invalid regex but got none")
	}

	if !strings.Contains(err.Error(), "invalid regex pattern") {
		t.Errorf("Expected error about invalid regex, got: %v", err)
	}
}
