# LogTail

An intelligent log analyzer for developers, written in Go.

## Features

- 🎨 **Syntax highlighting** : Automatic highlighting of log levels (ERROR, WARN, INFO, DEBUG)
- 🔍 **Real-time filtering** : Regular expression support for log filtering
- 📊 **Smart parser** : Automatic detection of timestamps, log levels and messages
- 📝 **Line numbering** : Option to display line numbers
- 🔄 **Follow mode** : Real-time file following like `tail -f`
- 📁 **Multi-file support** : Process multiple files simultaneously

## Installation

### From source

```bash
git clone https://github.com/educhastenier/logtail
cd logtail
go build -o logtail
```

### Using go install

```bash
go install github.com/educhastenier/logtail@latest
```

## Usage

### Basic examples

```bash
# Analyze a log file
./logtail app.log

# Read from stdin
cat app.log | ./logtail

# Filter with regex
./logtail -f "ERROR|FATAL" app.log

# Show line numbers
./logtail -n app.log

# Disable coloring
# Follow log files in real-time
./logtail --follow app.log

# Follow multiple log files with filtering
./logtail -F -f "ERROR|WARN" app.log error.log

# Use with pipes
tail -f app.log | ./logtail -f "ERROR"

# Disable colors for scripting
./logtail -c=false app.log > filtered.log

# Process multiple files
./logtail app.log error.log access.log
```

### Available options

- `-f, --filter` : Filter with regular expression
- `-c, --color` : Enable/disable coloring (default: true)
- `-n, --line-numbers` : Show line numbers
- `-F, --follow` : Follow file like tail -f for real-time monitoring

## Development

Project structure:
```
logtail/
├── main.go              # Entry point
├── cmd/                 # Command line interface
│   └── root.go
├── internal/            # Internal code
│   ├── parser/          # Log parsing
│   └── colorizer/       # Syntax highlighting
└── pkg/                 # Public packages (coming soon)
```

### Tests

```bash
go test ./...
```

### Build

```bash
go build -o logtail
```

## Supported log examples

LogTail automatically recognizes several log formats:

```
2024-09-30T10:30:45.123Z INFO This is an info message
2024/09/30 10:30:45 ERROR Database connection failed
Sep 30 10:30:45 WARN This is a warning
[ERROR] 2024-09-30 10:30:45 Something went wrong
```

## Performance

LogTail is designed to be fast and memory-efficient:
- Streaming processing: handles large files without loading them entirely into memory
- Optimized regex patterns for common log formats
- Minimal memory footprint with buffered I/O

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development setup

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code style

- Follow standard Go formatting (`go fmt`)
- Write meaningful commit messages
- Add comments for complex logic
- Include tests for new features

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [fatih/color](https://github.com/fatih/color) for terminal colors
- Inspired by various log analysis tools in the Go ecosystem

## Roadmap

- [x] Follow mode (`-F, --follow`)
- [ ] Export to different formats (JSON, CSV)
- [ ] Log statistics (counters per level)
- [ ] Common error pattern detection
- [ ] File-based configuration
- [ ] Integration with journald
- [ ] Plugins system for custom parsers
- [ ] Web interface for log analysis