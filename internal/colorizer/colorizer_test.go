package colorizer

import (
	"strings"
	"testing"

	"logtail/internal/parser"

	"github.com/fatih/color"
)

func TestColorizeLogLine(t *testing.T) {
	// Disable color for consistent testing
	originalNoColor := color.NoColor
	color.NoColor = true
	defer func() {
		color.NoColor = originalNoColor
	}()

	tests := []struct {
		name         string
		entry        parser.LogEntry
		originalLine string
		expectColor  bool
	}{
		{
			name: "ERROR level line",
			entry: parser.LogEntry{
				Level:   parser.LevelError,
				Message: "Database connection failed",
				Raw:     "2024-09-30 ERROR Database connection failed",
			},
			originalLine: "2024-09-30 ERROR Database connection failed",
			expectColor:  true,
		},
		{
			name: "WARN level line",
			entry: parser.LogEntry{
				Level:   parser.LevelWarn,
				Message: "Configuration missing",
				Raw:     "2024-09-30 WARN Configuration missing",
			},
			originalLine: "2024-09-30 WARN Configuration missing",
			expectColor:  true,
		},
		{
			name: "INFO level line",
			entry: parser.LogEntry{
				Level:   parser.LevelInfo,
				Message: "Application started",
				Raw:     "2024-09-30 INFO Application started",
			},
			originalLine: "2024-09-30 INFO Application started",
			expectColor:  true,
		},
		{
			name: "DEBUG level line",
			entry: parser.LogEntry{
				Level:   parser.LevelDebug,
				Message: "Loading configuration",
				Raw:     "2024-09-30 DEBUG Loading configuration",
			},
			originalLine: "2024-09-30 DEBUG Loading configuration",
			expectColor:  true,
		},
		{
			name: "FATAL level line",
			entry: parser.LogEntry{
				Level:   parser.LevelFatal,
				Message: "System crash",
				Raw:     "2024-09-30 FATAL System crash",
			},
			originalLine: "2024-09-30 FATAL System crash",
			expectColor:  true,
		},
		{
			name: "TRACE level line",
			entry: parser.LogEntry{
				Level:   parser.LevelTrace,
				Message: "Function entry",
				Raw:     "2024-09-30 TRACE Function entry",
			},
			originalLine: "2024-09-30 TRACE Function entry",
			expectColor:  true,
		},
		{
			name: "Unknown level with error keyword",
			entry: parser.LogEntry{
				Level:   parser.LevelUnknown,
				Message: "Connection error occurred",
				Raw:     "Connection error occurred",
			},
			originalLine: "Connection error occurred",
			expectColor:  false, // Special patterns should be applied
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizeLogLine(tt.entry, tt.originalLine)

			// With color disabled, result should equal original line for known levels
			// or potentially be modified for unknown levels with special patterns
			if tt.entry.Level != parser.LevelUnknown && result != tt.originalLine {
				t.Errorf("ColorizeLogLine() with color disabled should return original line, got %q", result)
			}

			// Test that function doesn't panic and returns a string
			if result == "" {
				t.Error("ColorizeLogLine() returned empty string")
			}
		})
	}
}

func TestColorizeSpecialPatterns(t *testing.T) {
	// Disable color for consistent testing
	originalNoColor := color.NoColor
	color.NoColor = true
	defer func() {
		color.NoColor = originalNoColor
	}()

	tests := []struct {
		name     string
		input    string
		hasURL   bool
		hasIP    bool
		hasError bool
		hasWarn  bool
	}{
		{
			name:   "Line with URL",
			input:  "Connecting to https://api.example.com/v1/users",
			hasURL: true,
		},
		{
			name:  "Line with IP address",
			input: "Request from 192.168.1.100 denied",
			hasIP: true,
		},
		{
			name:     "Line with error keyword",
			input:    "Connection error: timeout exceeded",
			hasError: true,
		},
		{
			name:     "Line with exception keyword",
			input:    "Unhandled exception in module",
			hasError: true,
		},
		{
			name:     "Line with failed keyword",
			input:    "Authentication failed for user",
			hasError: true,
		},
		{
			name:     "Line with failure keyword",
			input:    "Deployment failure detected",
			hasError: true,
		},
		{
			name:     "Line with panic keyword",
			input:    "System panic: out of memory",
			hasError: true,
		},
		{
			name:     "Line with fatal keyword",
			input:    "Fatal system error occurred",
			hasError: true,
		},
		{
			name:    "Line with warning keyword",
			input:   "Performance warning: high CPU usage",
			hasWarn: true,
		},
		{
			name:    "Line with warn keyword",
			input:   "Config warn: missing parameter",
			hasWarn: true,
		},
		{
			name:    "Line with deprecated keyword",
			input:   "Method is deprecated and will be removed",
			hasWarn: true,
		},
		{
			name:    "Line with obsolete keyword",
			input:   "Feature is obsolete",
			hasWarn: true,
		},
		{
			name:     "Complex line with multiple patterns",
			input:    "ERROR: Failed to connect to https://api.example.com from 10.0.0.1",
			hasURL:   true,
			hasIP:    true,
			hasError: true,
		},
		{
			name:  "Line with no special patterns",
			input: "Normal log message without special content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorizeSpecialPatterns(tt.input)

			// With color disabled, the function should still process the text
			// but not add color codes (though it might modify the text)
			if result == "" {
				t.Error("colorizeSpecialPatterns() returned empty string")
			}

			// Test that function handles the input without panicking
			// Actual color testing would require enabling colors and checking ANSI codes
		})
	}
}

func TestColorizeByLevel(t *testing.T) {
	tests := []struct {
		name  string
		level parser.LogLevel
	}{
		{
			name:  "ERROR level",
			level: parser.LevelError,
		},
		{
			name:  "FATAL level",
			level: parser.LevelFatal,
		},
		{
			name:  "WARN level",
			level: parser.LevelWarn,
		},
		{
			name:  "INFO level",
			level: parser.LevelInfo,
		},
		{
			name:  "DEBUG level",
			level: parser.LevelDebug,
		},
		{
			name:  "TRACE level",
			level: parser.LevelTrace,
		},
		{
			name:  "UNKNOWN level",
			level: parser.LevelUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorFunc := ColorizeByLevel(tt.level)

			// Test that we get a valid function
			if colorFunc == nil {
				t.Error("ColorizeByLevel() returned nil function")
			}

			// Test that the function works
			result := colorFunc("test message")
			if result == "" {
				t.Error("Color function returned empty string")
			}
		})
	}
}

func TestDisableColor(t *testing.T) {
	// Save original state
	originalNoColor := color.NoColor
	defer func() {
		color.NoColor = originalNoColor
	}()

	// Test enabling color first
	color.NoColor = false
	DisableColor()

	if !color.NoColor {
		t.Error("DisableColor() did not disable colors")
	}
}

// Test edge cases and error conditions
func TestColorizeLogLineEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		entry        parser.LogEntry
		originalLine string
	}{
		{
			name: "Empty line",
			entry: parser.LogEntry{
				Level: parser.LevelUnknown,
				Raw:   "",
			},
			originalLine: "",
		},
		{
			name: "Very long line",
			entry: parser.LogEntry{
				Level: parser.LevelInfo,
				Raw:   strings.Repeat("a", 10000),
			},
			originalLine: strings.Repeat("a", 10000),
		},
		{
			name: "Line with special characters",
			entry: parser.LogEntry{
				Level: parser.LevelError,
				Raw:   "Error: файл не найден (file not found) 错误",
			},
			originalLine: "Error: файл не найден (file not found) 错误",
		},
		{
			name: "Line with ANSI escape codes",
			entry: parser.LogEntry{
				Level: parser.LevelWarn,
				Raw:   "\x1b[31mRed text\x1b[0m normal text",
			},
			originalLine: "\x1b[31mRed text\x1b[0m normal text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := ColorizeLogLine(tt.entry, tt.originalLine)
			if len(result) == 0 && len(tt.originalLine) > 0 {
				t.Error("ColorizeLogLine() returned empty result for non-empty input")
			}
		})
	}
}

// Benchmark tests
func BenchmarkColorizeLogLine(b *testing.B) {
	entry := parser.LogEntry{
		Level:   parser.LevelError,
		Message: "Database connection failed",
		Raw:     "2024-09-30 ERROR Database connection failed",
	}
	line := "2024-09-30 ERROR Database connection failed"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ColorizeLogLine(entry, line)
	}
}

func BenchmarkColorizeSpecialPatterns(b *testing.B) {
	line := "ERROR: Failed to connect to https://api.example.com from 192.168.1.100"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		colorizeSpecialPatterns(line)
	}
}
