package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

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
	var readers []io.Reader

	// If no file is specified, read from stdin
	if len(args) == 0 {
		readers = append(readers, os.Stdin)
	} else {
		// Open specified files
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("cannot open file %s: %v", filename, err)
			}
			defer file.Close()
			readers = append(readers, file)
		}
	}

	// Compile filter pattern if provided
	var filter *regexp.Regexp
	if filterPattern != "" {
		var err error
		filter, err = regexp.Compile(filterPattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	// Process each source
	for i, reader := range readers {
		if len(args) > 1 {
			fmt.Printf("==> %s <==\n", args[i])
		}

		if err := processLogs(reader, filter); err != nil {
			return err
		}

		if i < len(readers)-1 {
			fmt.Println()
		}
	}

	return nil
}

func processLogs(reader io.Reader, filter *regexp.Regexp) error {
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
