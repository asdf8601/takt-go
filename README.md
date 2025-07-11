# takt-go

> [!NOTE]
> This is a Go implementation of [`takt`](https://github.com/asdf8601/takt), which was originally written in Python.

Takt-go is a command-line tool for tracking time, inspired by the original Takt project. It aims to provide similar functionality while leveraging the strengths of the Go programming language.

## Features

- Simple and efficient time tracking
- Command-line interface for easy use
- Data stored in CSV format for compatibility and ease of use
- **Overtime/Undertime tracking** - Compare actual hours worked against target hours
- **Flexible target hours** - Support for both decimal (7.5) and time format (7:30)
- **Comprehensive reporting** - Daily, weekly, monthly, and yearly summaries with overtime calculations

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

**Example output with overtime tracking:**
```
Date          Total	Days	   Avg	 Balance
2024-07-26    9h00m	   1	 9h00m	  +1h00m
2024-07-25    7h30m	   1	 7h30m	  -0h30m
2024-07-24    8h00m	   1	 8h00m	  00h00m
```

The **Balance** column shows:
- `+1h00m` - 1 hour overtime
- `-0h30m` - 30 minutes undertime  
- `00h00m` - exactly on target

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

# Set target daily hours (default: 8 hours)
export TAKT_TARGET_HOURS=8          # decimal format
export TAKT_TARGET_HOURS=7:30       # time format (7h 30m)
export TAKT_TARGET_HOURS=8:15       # time format (8h 15m)
```

#### Target Hours Format

The `TAKT_TARGET_HOURS` environment variable supports two formats:

**Decimal Format:**
- `7.5` → 7.5 hours
- `8.25` → 8.25 hours (8 hours 15 minutes)

**Time Format:**
- `7:30` → 7.5 hours (7 hours 30 minutes)
- `8:15` → 8.25 hours (8 hours 15 minutes)
- `10:00` → 10 hours

**Validation:**
- Hours must be non-negative integers
- Minutes must be between 0-59
- Invalid formats fall back to default 8 hours

#### How Overtime/Undertime is Calculated

- **Daily**: Actual hours - Target hours
- **Weekly/Monthly/Yearly**: Actual hours - (Target hours × Working days)

**Examples:**
```bash
# With 8-hour target
TAKT_TARGET_HOURS=8 takt d 1
# 9 hours worked = +1h00m overtime

# With 7.5-hour target  
TAKT_TARGET_HOURS=7:30 takt m 1
# 30 hours worked in 4 days = 30 - (7.5 × 4) = 0h00m
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
- [x] Create tests
- [x] Add overtime/undertime tracking
- [x] Support flexible target hours format

## License

(Add license information here)

[takt]: https://github.com/asdf8601/takt
