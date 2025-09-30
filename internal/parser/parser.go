package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelTrace   LogLevel = "TRACE"
	LevelDebug   LogLevel = "DEBUG"
	LevelInfo    LogLevel = "INFO"
	LevelWarn    LogLevel = "WARN"
	LevelError   LogLevel = "ERROR"
	LevelFatal   LogLevel = "FATAL"
	LevelUnknown LogLevel = "UNKNOWN"
)

// LogEntry represents a parsed log line
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Source    string
	Raw       string
}

var (
	// Common log patterns
	timestampPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d{3})?(?:Z|[+-]\d{2}:\d{2})?)`),
		regexp.MustCompile(`(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})`),
		regexp.MustCompile(`(\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2})`),
		regexp.MustCompile(`(\w{3} \d{2} \d{2}:\d{2}:\d{2})`),
	}

	levelPattern = regexp.MustCompile(`(?i)\b(TRACE|DEBUG|INFO|WARN|WARNING|ERROR|ERR|FATAL|PANIC)\b`)
)

// ParseLogLine attempts to parse a log line and extract structured information
func ParseLogLine(line string) LogEntry {
	entry := LogEntry{
		Raw:     line,
		Level:   LevelUnknown,
		Message: line,
	}

	// Try to extract timestamp
	for _, pattern := range timestampPatterns {
		if matches := pattern.FindStringSubmatch(line); len(matches) > 1 {
			if timestamp, err := parseTimestamp(matches[1]); err == nil {
				entry.Timestamp = timestamp
				break
			}
		}
	}

	// Extract log level
	if matches := levelPattern.FindStringSubmatch(line); len(matches) > 1 {
		level := strings.ToUpper(matches[1])
		switch level {
		case "WARNING":
			entry.Level = LevelWarn
		case "ERR":
			entry.Level = LevelError
		case "PANIC":
			entry.Level = LevelFatal
		default:
			entry.Level = LogLevel(level)
		}
	}

	// Extract message (everything after level, or full line if no level found)
	if entry.Level != LevelUnknown {
		// Find the original level text in the line to extract message after it
		if matches := levelPattern.FindStringSubmatch(line); len(matches) > 1 {
			originalLevel := matches[1]
			if idx := strings.Index(line, originalLevel); idx != -1 {
				start := idx + len(originalLevel)
				if start < len(line) {
					entry.Message = strings.TrimSpace(line[start:])
				}
			}
		}
	}

	return entry
}

// parseTimestamp tries to parse timestamp from string using common formats
func parseTimestamp(timestampStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006/01/02 15:04:05",
		"01/02/2006 15:04:05",
		"Jan 02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestampStr)
}

// IsErrorLevel returns true if the log level indicates an error or fatal condition
func (entry LogEntry) IsErrorLevel() bool {
	return entry.Level == LevelError || entry.Level == LevelFatal
}

// IsWarningLevel returns true if the log level indicates a warning
func (entry LogEntry) IsWarningLevel() bool {
	return entry.Level == LevelWarn
}

// IsInfoLevel returns true if the log level indicates info, debug, or trace
func (entry LogEntry) IsInfoLevel() bool {
	return entry.Level == LevelInfo || entry.Level == LevelDebug || entry.Level == LevelTrace
}
