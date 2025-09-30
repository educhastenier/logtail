package colorizer

import (
	"regexp"
	"strings"

	"logtail/internal/parser"

	"github.com/fatih/color"
)

var (
	// Color functions for different elements
	errorColor     = color.New(color.FgRed, color.Bold)
	warnColor      = color.New(color.FgYellow, color.Bold)
	infoColor      = color.New(color.FgCyan)
	debugColor     = color.New(color.FgMagenta)
	timestampColor = color.New(color.FgBlue)
	sourceColor    = color.New(color.FgGreen)

	// Patterns to identify special elements
	urlPattern    = regexp.MustCompile(`https?://[^\s]+`)
	ipPattern     = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	numberPattern = regexp.MustCompile(`\b\d+\b`)
)

// ColorizeLogLine colorizes a log line based on its parsed content
func ColorizeLogLine(entry parser.LogEntry, originalLine string) string {
	line := originalLine

	// Colorize according to log level
	switch entry.Level {
	case parser.LevelError, parser.LevelFatal:
		line = errorColor.Sprint(line)
	case parser.LevelWarn:
		line = warnColor.Sprint(line)
	case parser.LevelInfo:
		line = infoColor.Sprint(line)
	case parser.LevelDebug, parser.LevelTrace:
		line = debugColor.Sprint(line)
	default:
		// Apply specific colorizations if no level detected
		line = colorizeSpecialPatterns(line)
	}

	return line
}

// colorizeSpecialPatterns colors special patterns in the line
func colorizeSpecialPatterns(line string) string {
	// Colorize URLs
	line = urlPattern.ReplaceAllStringFunc(line, func(url string) string {
		return color.BlueString(url)
	})

	// Colorize IP addresses
	line = ipPattern.ReplaceAllStringFunc(line, func(ip string) string {
		return color.CyanString(ip)
	})

	// Identify and colorize error keywords even without formal level
	errorKeywords := []string{"error", "exception", "failed", "failure", "panic", "fatal"}
	for _, keyword := range errorKeywords {
		if strings.Contains(strings.ToLower(line), keyword) {
			pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
			line = pattern.ReplaceAllStringFunc(line, func(match string) string {
				return errorColor.Sprint(match)
			})
		}
	}

	// Identify and colorize warning keywords
	warnKeywords := []string{"warning", "warn", "deprecated", "obsolete"}
	for _, keyword := range warnKeywords {
		if strings.Contains(strings.ToLower(line), keyword) {
			pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
			line = pattern.ReplaceAllStringFunc(line, func(match string) string {
				return warnColor.Sprint(match)
			})
		}
	}

	return line
}

// ColorizeByLevel returns a coloring function based on the level
func ColorizeByLevel(level parser.LogLevel) func(...interface{}) string {
	switch level {
	case parser.LevelError, parser.LevelFatal:
		return errorColor.Sprint
	case parser.LevelWarn:
		return warnColor.Sprint
	case parser.LevelInfo:
		return infoColor.Sprint
	case parser.LevelDebug, parser.LevelTrace:
		return debugColor.Sprint
	default:
		return color.New().Sprint
	}
}

// DisableColor disables coloring (useful for pipes and redirections)
func DisableColor() {
	color.NoColor = true
}
