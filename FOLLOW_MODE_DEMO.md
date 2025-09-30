# LogTail Follow Mode Demo

## Demo Script for Follow Mode Functionality

This script demonstrates the new `--follow` functionality of LogTail.

### Test 1: Basic Follow Mode

```bash
# Create initial log file
echo "2024-09-30T10:00:00.000Z INFO Application startup" > demo.log

# Start LogTail in follow mode (in one terminal)
./logtail --follow demo.log

# In another terminal, add new log entries
echo "2024-09-30T10:01:00.000Z ERROR Database connection failed" >> demo.log
echo "2024-09-30T10:02:00.000Z WARN High memory usage detected" >> demo.log
echo "2024-09-30T10:03:00.000Z INFO Database connection restored" >> demo.log
```

### Test 2: Follow Mode with Filtering

```bash
# Follow only ERROR and WARN messages
./logtail --follow -f "ERROR|WARN" demo.log

# Add mixed log levels
echo "2024-09-30T10:04:00.000Z DEBUG Processing request 123" >> demo.log
echo "2024-09-30T10:05:00.000Z ERROR Authentication failed" >> demo.log
echo "2024-09-30T10:06:00.000Z INFO Request processed successfully" >> demo.log
echo "2024-09-30T10:07:00.000Z WARN Performance degradation detected" >> demo.log
```

### Test 3: Multi-file Follow Mode

```bash
# Create multiple log files
echo "2024-09-30T10:00:00.000Z INFO Service A started" > service_a.log
echo "2024-09-30T10:00:00.000Z INFO Service B started" > service_b.log

# Follow multiple files simultaneously
./logtail --follow service_a.log service_b.log

# Add logs to both files
echo "2024-09-30T10:08:00.000Z ERROR Service A database error" >> service_a.log
echo "2024-09-30T10:08:30.000Z WARN Service B high CPU usage" >> service_b.log
echo "2024-09-30T10:09:00.000Z INFO Service A recovered" >> service_a.log
```

### Test 4: Follow Mode with Line Numbers

```bash
# Follow with line numbers enabled
./logtail --follow --line-numbers -f "ERROR" demo.log

# Add more errors
echo "2024-09-30T10:10:00.000Z ERROR Connection timeout" >> demo.log
echo "2024-09-30T10:11:00.000Z INFO Connection retry" >> demo.log
echo "2024-09-30T10:12:00.000Z ERROR Max retries exceeded" >> demo.log
```

## Expected Behavior

1. **Immediate Display**: Shows existing file content immediately
2. **Real-time Updates**: New lines appear as they're written to the file
3. **Filtering Works**: Only lines matching regex patterns are shown
4. **Multi-file Support**: Displays content from multiple files with prefixes
5. **Continuous Monitoring**: Runs indefinitely until interrupted (Ctrl+C)
6. **Memory Efficient**: Streams content without loading entire files

## Use Cases

- **Development**: Monitor application logs during development
- **Debugging**: Follow error logs in real-time
- **Operations**: Monitor multiple service logs simultaneously
- **Analysis**: Filter specific patterns while following logs
- **CI/CD**: Monitor deployment logs during automated deployments

## Performance Notes

- **Polling Interval**: 100ms between checks for new content
- **Memory Usage**: Minimal - only buffers one line at a time
- **CPU Usage**: Low - sleeps when no new content available
- **File Handling**: Properly handles file rotation and growth