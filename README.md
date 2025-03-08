# takt-go

> [!NOTE]
> This is a Go implementation of [`takt`](https://github.com/asdf8601/takt), which was originally written in Python.

Takt-go is a command-line tool for tracking time, inspired by the original Takt project. It aims to provide similar functionality while leveraging the strengths of the Go programming language.

## Features

- Simple and efficient time tracking
- Command-line interface for easy use
- Data stored in CSV format for compatibility and ease of use

## Installation

To install takt-go, make sure you have Go installed on your system, then run:

```
go install github.com/yourusername/takt-go@latest
```

Replace `yourusername` with the actual GitHub username or organization where this project is hosted.

## Usage

Basic commands to get started with takt:

### Check In/Out

```bash
# Check in with a note
takt check "Starting work on feature X"

# Check out 
takt check "Completed feature X"

# Short form
takt c "Meeting with team"
```

### View Records

```bash
# Show last 10 records (default)
takt cat

# Show last 5 records
takt cat 5

# Edit records manually
takt edit   # or takt e
```

### Summary Reports

```bash
# Daily summary (last 10 days)
takt day    # or takt d

# Weekly summary
takt week   # or takt w

# Monthly summary
takt month  # or takt m

# Yearly summary
takt year   # or takt y

# Show more entries
takt month 12  # show last 12 months
```

### Grid View

```bash
# Show current year's grid
takt grid

# Show specific year with legend
takt grid 2023 true
```

### Git Integration

```bash
# Commit and push changes
takt commit  # or takt cm
```

### Configuration

Takt can be configured using environment variables:

```bash
# Set custom file location
export TAKT_FILE=~/my-time.csv

# Set preferred editor
export TAKT_EDITOR=vim
```

The grid view uses these symbols:
- 󰋣 : 0h00m - 1h00m
-  : 1h00m - 4h00m
-  : 4h00m - 8h00m
- 󰈸 : 8h00m - 12h00m
-  : 12h00m or more

## About the Name

The name "Takt" is derived from the German word "Taktzeit," which means cycle time. It is a key principle in lean manufacturing, describing the pace of production that aligns with customer demand. This tool aims to help you track and manage your time with similar precision.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## TODO

- [ ] Add demo
- [ ] Implement plugin support
- [ ] Add more detailed usage instructions
- [ ] Create tests

## License

(Add license information here)

[takt]: https://github.com/asdf8601/takt
