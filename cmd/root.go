package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"logtail/internal/colorizer"
	"logtail/internal/parser"

	"github.com/spf13/cobra"
)

var (
	filterPattern string
	colorOutput   bool
	followMode    bool
	showLineNum   bool
)

var rootCmd = &cobra.Command{
	Use:   "logtail [file...]",
	Short: "An intelligent log analyzer for developers",
	Long: `LogTail is a powerful command-line tool for parsing, filtering, and analyzing log files.
It provides real-time filtering, syntax highlighting, and pattern detection.`,
	RunE: runLogTail,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&filterPattern, "filter", "f", "", "Filter logs with regex pattern")
	rootCmd.Flags().BoolVarP(&colorOutput, "color", "c", true, "Enable colorized output")
	rootCmd.Flags().BoolVarP(&followMode, "follow", "F", false, "Follow log file like tail -f")
	rootCmd.Flags().BoolVarP(&showLineNum, "line-numbers", "n", false, "Show line numbers")
}

func runLogTail(cmd *cobra.Command, args []string) error {
	// Compile filter pattern if provided
	var filter *regexp.Regexp
	if filterPattern != "" {
		var err error
		filter, err = regexp.Compile(filterPattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	// Handle stdin case
	if len(args) == 0 {
		return processLogs(os.Stdin, filter, "", false)
	}

	// Follow mode only works with files
	if followMode {
		return followFiles(args, filter)
	}

	// Normal mode: process files sequentially
	for i, filename := range args {
		if len(args) > 1 {
			fmt.Printf("==> %s <==\n", filename)
		}

		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("cannot open file %s: %v", filename, err)
		}

		err = processLogs(file, filter, filename, false)
		file.Close()

		if err != nil {
			return err
		}

		if i < len(args)-1 {
			fmt.Println()
		}
	}

	return nil
}

type FollowFile struct {
	file     *os.File
	filename string
	scanner  *bufio.Scanner
	lineNum  int
}

func followFiles(filenames []string, filter *regexp.Regexp) error {
	// For follow mode, we need to track file positions and watch for changes
	files := make([]*FollowFile, len(filenames))

	// Initialize files
	for i, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("cannot open file %s: %v", filename, err)
		}

		// First, read existing content from the beginning
		scanner := bufio.NewScanner(file)
		lineNum := 1

		if len(filenames) > 1 {
			fmt.Printf("==> %s <==\n", filename)
		}

		// Process existing content
		for scanner.Scan() {
			line := scanner.Text()

			// Apply filter if defined
			if filter != nil && !filter.MatchString(line) {
				lineNum++
				continue
			}

			// Parse the log line
			logEntry := parser.ParseLogLine(line)

			// Display the line
			output := line
			if colorOutput {
				output = colorizer.ColorizeLogLine(logEntry, line)
			}

			// Add filename prefix for multiple files
			if len(filenames) > 1 {
				prefix := fmt.Sprintf("[%s] ", filename)
				if showLineNum {
					fmt.Printf("%s%6d: %s\n", prefix, lineNum, output)
				} else {
					fmt.Printf("%s%s\n", prefix, output)
				}
			} else {
				if showLineNum {
					fmt.Printf("%6d: %s\n", lineNum, output)
				} else {
					fmt.Println(output)
				}
			}

			lineNum++
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			file.Close()
			return fmt.Errorf("error reading file %s: %v", filename, err)
		}

		files[i] = &FollowFile{
			file:     file,
			filename: filename,
			scanner:  scanner,
			lineNum:  lineNum,
		}
	}

	// Cleanup
	defer func() {
		for _, f := range files {
			f.file.Close()
		}
	}()

	// Follow loop
	for {
		hasNewContent := false

		for _, f := range files {
			for f.scanner.Scan() {
				hasNewContent = true
				line := f.scanner.Text()

				// Apply filter if defined
				if filter != nil && !filter.MatchString(line) {
					f.lineNum++
					continue
				}

				// Parse the log line
				logEntry := parser.ParseLogLine(line)

				// Display the line
				output := line
				if colorOutput {
					output = colorizer.ColorizeLogLine(logEntry, line)
				}

				// Add filename prefix for multiple files
				if len(filenames) > 1 {
					prefix := fmt.Sprintf("[%s] ", f.filename)
					if showLineNum {
						fmt.Printf("%s%6d: %s\n", prefix, f.lineNum, output)
					} else {
						fmt.Printf("%s%s\n", prefix, output)
					}
				} else {
					if showLineNum {
						fmt.Printf("%6d: %s\n", f.lineNum, output)
					} else {
						fmt.Println(output)
					}
				}

				f.lineNum++
			}

			// Check for scanner errors
			if err := f.scanner.Err(); err != nil {
				return fmt.Errorf("error reading file %s: %v", f.filename, err)
			}
		}

		// If no new content, sleep briefly before checking again
		if !hasNewContent {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func processLogs(reader io.Reader, filter *regexp.Regexp, filename string, isFollow bool) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()

		// Apply filter if defined
		if filter != nil && !filter.MatchString(line) {
			lineNum++
			continue
		}

		// Parse the log line
		logEntry := parser.ParseLogLine(line)

		// Display the line
		output := line
		if colorOutput {
			output = colorizer.ColorizeLogLine(logEntry, line)
		}

		if showLineNum {
			fmt.Printf("%6d: %s\n", lineNum, output)
		} else {
			fmt.Println(output)
		}

		lineNum++
	}

	return scanner.Err()
}
