package main

// Final test of simplified CI/CD workflow with path filters
import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	Version            = "2025-01-09"
	DefaultHead        = 10
	DefaultTargetHours = 8.0

	// Grid constants
	DaysPerYear = 365
	GridColumns = 9
	GridRows    = 54

	// Time format constants
	TimeFormat = time.RFC3339
	DateFormat = "2006-01-02"

	// Hour thresholds for grid display
	LowHours      = 1.0
	MediumHours   = 4.0
	HighHours     = 8.0
	VeryHighHours = 12.0

	// Grid symbols
	SymbolMinimal   = "󰋣 " // 0-1 hours
	SymbolLight     = "▪ " // 1-4 hours
	SymbolNormal    = "▮ " // 4-8 hours
	SymbolHeavy     = "󰈸 " // 8-12 hours
	SymbolVeryHeavy = "󰯆 " // 12+ hours
)

// Config holds application configuration
type Config struct {
	Editor      string
	FileName    string
	TargetHours float64
}

// LoadConfig initializes configuration from environment variables
func LoadConfig() (*Config, error) {
	fileName, err := getFileName("TAKT_FILE", "~/takt.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to get file name: %w", err)
	}

	targetHours, err := getTargetHours("TAKT_TARGET_HOURS", 8.0)
	if err != nil {
		return nil, fmt.Errorf("failed to get target hours: %w", err)
	}

	return &Config{
		Editor:      os.Getenv("TAKT_EDITOR"),
		FileName:    fileName,
		TargetHours: targetHours,
	}, nil
}

// Global configuration
var config *Config

// CSV Header
var Header = []string{"timestamp", "kind", "notes"}

type Record struct {
	Timestamp time.Time
	Kind      string
	Notes     string
}

type AggregatedRecord struct {
	Group        string
	TotalHours   float64
	Dates        []string
	Notes        []string
	AverageHours float64
}

// printGrid prints the grid of the records.
func printGrid(year string, legend bool) error {
	records, err := readRecords(1)
	if err != nil {
		return fmt.Errorf("failed to read latest record: %w", err)
	}

	if len(records) == 0 {
		return errors.New("no records found")
	}

	lastDay := records[0].Timestamp.Format(DateFormat)

	records, err = readRecords(-1)
	if err != nil {
		return fmt.Errorf("failed to read all records: %w", err)
	}

	agg, err := calculateDuration(records, "day")
	if err != nil {
		return fmt.Errorf("error calculating duration: %w", err)
	}

	daysAgg := make(map[string]AggregatedRecord)
	for _, a := range agg {
		daysAgg[a.Group] = a
	}

	grid := make([][GridColumns]string, GridRows)
	for i := range grid {
		for j := range grid[i] {
			grid[i][j] = "  "
		}
	}

	startDate := year + "-01-01"
	t, err := time.Parse(DateFormat, startDate)
	if err != nil {
		return fmt.Errorf("invalid year format: %w", err)
	}

	value := ""
	lastIdx := -1

	for i := 0; i < DaysPerYear; i++ {
		day := t.Format(DateFormat)
		currentYear, week := t.ISOWeek()
		if day > lastDay {
			lastIdx = week
			break
		}
		if t.Year() != currentYear {
			continue
		}
		dayOfWeek := t.Weekday()

		item := daysAgg[day]
		hours := item.TotalHours
		if hours < LowHours {
			value = SymbolMinimal
		} else if hours < MediumHours {
			value = SymbolLight
		} else if hours < HighHours {
			value = SymbolNormal
		} else if hours < VeryHighHours {
			value = SymbolHeavy
		} else {
			value = SymbolVeryHeavy
		}

		grid[week][0] = day
		grid[week][1] = fmt.Sprintf("%02d", week)
		grid[week][dayOfWeek+2] = value
		t = t.Add(time.Hour * 24)
	}

	if lastIdx > 0 {
		grid = grid[:lastIdx+1]
	}

	printGridOutput(grid, lastIdx, legend)
	return nil
}

// printGridOutput prints the formatted grid with improved formatting
func printGridOutput(grid [][GridColumns]string, lastIdx int, legend bool) {
	pad := "    "

	// ANSI color codes
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorYellow = "\033[33m"
		colorGreen  = "\033[32m"
		colorBlue   = "\033[34m"
		colorOrange = "\033[38;5;208m"
		colorPurple = "\033[35m"
		colorCyan   = "\033[36m"
		colorGray   = "\033[37m"
		colorBold   = "\033[1m"
	)

	// Print header with better formatting
	fmt.Printf("%s%s%s%-10s W  M  T  W  T  F  S  S%s\n", pad, colorBold, colorBlue, "Date", colorReset)
	fmt.Printf("%s%s══════════════════════════════════%s", pad, colorBlue, colorReset)

	var stats struct {
		totalDays     int
		activeDays    int
		totalHours    float64
		lightDays     int
		normalDays    int
		heavyDays     int
		veryHeavyDays int
	}

	var currentMonth string

	for idx, week := range grid {
		if week[0] == "" {
			continue
		}

		// Extract month from date for separators
		if len(week[0]) >= 7 {
			weekMonth := week[0][:7] // "2025-01"
			if currentMonth != "" && currentMonth != weekMonth {
				fmt.Printf("%s%s──────────────────────────────────%s\n", pad, colorBlue, colorReset)
			}
			currentMonth = weekMonth
		}

		// Color code the symbols and collect stats
		coloredWeek := make([]string, len(week))
		copy(coloredWeek, week[:])

		for i := 2; i < len(week); i++ {
			symbol := week[i]
			if symbol != "  " {
				stats.totalDays++
				if symbol != SymbolMinimal {
					stats.activeDays++
				}

				switch symbol {
				case SymbolMinimal:
					coloredWeek[i] = colorGray + symbol + colorReset
					stats.lightDays++
				case SymbolLight:
					coloredWeek[i] = colorYellow + symbol + colorReset
					stats.normalDays++
				case SymbolNormal:
					coloredWeek[i] = colorGreen + symbol + colorReset
					stats.normalDays++
				case SymbolHeavy:
					coloredWeek[i] = colorOrange + symbol + colorReset
					stats.heavyDays++
				case SymbolVeryHeavy:
					coloredWeek[i] = colorRed + symbol + colorReset
					stats.veryHeavyDays++
				}
			}
		}
		// Print the week with better alignment
		fmt.Printf("%s%s%s%s %s %s %s %s %s %s %s %s\n",
			pad, colorBold, coloredWeek[0], colorReset,
			coloredWeek[1], coloredWeek[3], coloredWeek[4],
			coloredWeek[5], coloredWeek[6], coloredWeek[7],
			coloredWeek[8], coloredWeek[2])

		if (lastIdx > 0) && (idx == lastIdx) {
			break
		}
	}

	// Print summary statistics
	if stats.totalDays > 0 {
		fmt.Printf("\n%s%s📊 Summary:%s\n", pad, colorBold, colorReset)
		fmt.Printf("%s%s├─ Total tracked days: %d%s\n", pad, colorBlue, stats.totalDays, colorReset)
		fmt.Printf("%s%s├─ Active work days: %d%s\n", pad, colorBlue, stats.activeDays, colorReset)
		if stats.activeDays > 0 {
			activePercent := float64(stats.activeDays) / float64(stats.totalDays) * 100
			fmt.Printf("%s%s└─ Activity rate: %.1f%%%s\n", pad, colorBlue, activePercent, colorReset)
		}
	}

	if legend {
		fmt.Printf("\n%s%s🎨 Legend:%s\n", pad, colorBold, colorReset)
		fmt.Printf("%s%s├─ %s%s%s 0h00m - 1h00m   (Minimal work)%s\n", pad, colorGray, colorGray, SymbolMinimal, colorReset, colorReset)
		fmt.Printf("%s%s├─ %s%s%s 1h00m - 4h00m   (Light work)%s\n", pad, colorYellow, colorYellow, SymbolLight, colorReset, colorReset)
		fmt.Printf("%s%s├─ %s%s%s 4h00m - 8h00m   (Normal work)%s\n", pad, colorGreen, colorGreen, SymbolNormal, colorReset, colorReset)
		fmt.Printf("%s%s├─ %s%s%s 8h00m - 12h00m  (Heavy work)%s\n", pad, colorBlue, colorBlue, SymbolHeavy, colorReset, colorReset)
		fmt.Printf("%s%s└─ %s%s%s 12h00m or more  (Very heavy work)%s\n", pad, colorPurple, colorPurple, SymbolVeryHeavy, colorReset, colorReset)
	}
}

// findGitRoot finds the git root directory starting from the config file directory.
func findGitRoot() (string, error) {
	if config == nil {
		return "", errors.New("config not initialized")
	}

	dir := filepath.Dir(config.FileName)
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("couldn't get absolute path: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		if dir == "/" {
			return "", errors.New("not in a git repository")
		}
		dir = filepath.Join(dir, "..")
	}
}

// gitPush pushes the file to the git repository.
func gitPush() error {
	gitRoot, _ := findGitRoot()
	gitCmd := exec.Command("git", "-C", gitRoot, "push")
	err := execBashCmd(gitCmd)
	return err
}

// gitCommit commits the file to the git repository.
func gitCommit() error {
	gitRoot, _ := findGitRoot()
	gitCmd := exec.Command("git", "-C", gitRoot, "commit", "-m", "Automatic commit from Takt")
	err := execBashCmd(gitCmd)
	return err
}

// gitAdd adds the file to the git repository.
func gitAdd() error {
	if config == nil {
		return errors.New("config not initialized")
	}

	gitRoot, err := findGitRoot()
	if err != nil {
		return err
	}

	dir := filepath.Dir(config.FileName)
	dir, err = filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("couldn't get absolute path: %w", err)
	}

	fileDirRel, err := filepath.Rel(gitRoot, dir)
	if err != nil {
		return fmt.Errorf("couldn't get relative path: %w", err)
	}

	fileNameAbs := filepath.Join(fileDirRel, filepath.Base(config.FileName))
	gitCmd := exec.Command("git", "-C", gitRoot, "add", fileNameAbs)
	return execBashCmd(gitCmd)
}

// execBashCmd executes a bash command.
func execBashCmd(cmd *exec.Cmd) error {

	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Print("error= " + err.Error())
	}

	slurp, _ := io.ReadAll(stderr)
	if slurp != nil {
		fmt.Printf("%s\n", slurp)
	}

	if err := cmd.Wait(); err != nil {
		if e, ok := err.(interface{ ExitCode() int }); ok {
			if e.ExitCode() != 1 {
				// exit code is neither zero (as we have an error) or one
				fmt.Print("error= " + err.Error())
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// absPath returns the absolute path by expanding the tilde (~) to the user's home directory.
func absPath(path string) (string, error) {
	if path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error: could not get user home directory")
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// getTargetHours returns the target hours from the environment variable or the default value.
// Supports both float format (e.g., "7.5") and time format (e.g., "7:30").
func getTargetHours(key string, dflt float64) (float64, error) {
	value := os.Getenv(key)
	if value == "" {
		return dflt, nil
	}

	// Check if the value contains a colon (HH:MM format)
	if strings.Contains(value, ":") {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return dflt, nil
		}

		hours, err := strconv.Atoi(parts[0])
		if err != nil || hours < 0 {
			return dflt, nil
		}

		minutes, err := strconv.Atoi(parts[1])
		if err != nil || minutes < 0 || minutes >= 60 {
			return dflt, nil
		}

		// Convert to decimal hours
		return float64(hours) + float64(minutes)/60.0, nil
	}

	// Try to parse as float
	hours, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return dflt, nil
	}

	return hours, nil
}

// getFileName returns the file name from the environment variable or the default value.
func getFileName(key, dflt string) (string, error) {
	path := os.Getenv(key)

	if path == "" {
		return absPath(dflt)
	}

	return absPath(path)
}

// sortedKeys returns the keys of a map sorted in descending order.
func sortedKeys(m map[string]AggregatedRecord) []string {
	// Crear un slice para las claves
	keys := make([]string, 0, len(m))

	// Agregar las claves al slice
	for k := range m {
		keys = append(keys, k)
	}

	// Ordenar las claves
	sort.Strings(keys)

	// invert the order
	out := make([]string, 0, len(keys))
	for i := len(keys) - 1; i >= 0; i-- {
		out = append(out, keys[i])
	}

	return out
}

// hoursToText converts hours to a human-readable format.
func hoursToText(totalHours float64) string {
	if totalHours <= 0 {
		return "00h00m"
	} else if totalHours <= 24 {
		hours := int(totalHours)
		minutes := int(math.Round((totalHours - float64(hours)) * 60))
		return fmt.Sprintf("%dh%02dm", hours, minutes)
	} else {
		days := int(totalHours / 24)
		hours := int(totalHours) % 24
		minutes := int(math.Round((totalHours - float64(days*24+hours)) * 60))
		return fmt.Sprintf("%dd%02dh%02dm", days, hours, minutes)
	}
}

// formatOvertime formats the overtime/undertime difference with a sign.
// formatOvertime formats the overtime/undertime difference using TARGET_HOUR as the day unit
func formatOvertime(difference float64) string {
	if difference == 0 {
		return "00h00m"
	}

	sign := ""
	if difference > 0 {
		sign = "+"
	} else {
		sign = "-"
	}

	// Use absolute value for formatting
	absDiff := math.Abs(difference)

	// Calculate days based on TARGET_HOUR
	targetHour := config.TargetHours
	if targetHour == 0 {
		targetHour = DefaultTargetHours // fallback to default if config is not set
	}

	if absDiff >= targetHour {
		days := int(absDiff / targetHour)
		remainingHours := absDiff - (float64(days) * targetHour)
		hours := int(remainingHours)
		minutes := int(math.Round((remainingHours - float64(hours)) * 60))

		// Handle case where minutes round to 60
		if minutes >= 60 {
			hours += minutes / 60
			minutes = minutes % 60
		}

		if hours == 0 && minutes == 0 {
			return fmt.Sprintf("%s%dd", sign, days)
		} else if minutes == 0 {
			return fmt.Sprintf("%s%dd%dh", sign, days, hours)
		} else {
			return fmt.Sprintf("%s%dd%dh%02dm", sign, days, hours, minutes)
		}
	} else {
		hours := int(absDiff)
		minutes := int(math.Round((absDiff - float64(hours)) * 60))

		// Handle case where minutes round to 60
		if minutes >= 60 {
			hours += minutes / 60
			minutes = minutes % 60
		}

		return fmt.Sprintf("%s%dh%02dm", sign, hours, minutes)
	}
}

// summary prints a summary of the records.
func summary(offset string, head int) {
	records, err := readRecords(-1)
	if err != nil {
		log.Fatal(err)
	}
	agg, err := calculateDuration(records, offset)
	if err != nil {
		log.Fatalf("error calculating duration: %v", err)
	}

	if head < 1 || head > len(agg) {
		head = len(agg)
	}

	var outFmt string
	if offset == "day" {
		outFmt = "%-12s %6s\t%4s\t%6s\t%8s\n"
		fmt.Printf(outFmt, "Date", "Total", "Days", "Avg", "Balance")
	} else {
		// wider total hours column for week, month, year
		outFmt = "%-8s %10s\t%4s\t%6s\t%8s\n"
		fmt.Printf(outFmt, "Date", "Total", "Days", "Avg", "Balance")
	}

	for i := 0; i < head; i++ {
		a := agg[i]
		hhmm := hoursToText(a.TotalHours)
		ndays := strconv.Itoa(len(a.Dates))
		avg := hoursToText(a.AverageHours)

		// For all periods, calculate expected hours = target hours * working days
		expectedHours := config.TargetHours * float64(len(a.Dates))
		diff := a.TotalHours - expectedHours
		overtime := formatOvertime(diff)
		fmt.Printf(outFmt, a.Group, hhmm, ndays, avg, overtime)
	}
}

// contains returns true if the item is in the slice.
func contains(items []string, item string) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}

// unique returns a slice with unique items.
func unique(items []string) []string {

	out := []string{}
	for _, it := range items {
		if !contains(out, it) {
			out = append(out, it)
		}
	}
	return out
}

// calculateDuration calculates the duration of the records.
func calculateDuration(records []Record, period string) ([]AggregatedRecord, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to process")
	}

	inferLastOut(&records)

	var aggregations map[string]AggregatedRecord
	var labeler func(time.Time) string

	switch period {
	case "day":
		labeler = func(t time.Time) string {
			return t.Format("2006-01-02")
		}
	case "week":
		labeler = func(t time.Time) string {
			year, week := t.ISOWeek()
			return fmt.Sprintf("%d-W%02d", year, week)
		}
	case "month":
		labeler = func(t time.Time) string {
			return t.Format("2006-01")
		}
	case "year":
		labeler = func(t time.Time) string {
			return t.Format("2006")
		}
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	aggregations = aggregateBy(records, labeler)
	var out []AggregatedRecord
	keys := sortedKeys(aggregations)
	for _, k := range keys {
		v := aggregations[k]
		v.Dates = unique(v.Dates)
		v.AverageHours = v.TotalHours / float64(len(v.Dates))
		out = append(out, v)
	}
	return out, nil
}

// aggregateBy aggregates the records by the groupFunc.
func aggregateBy(records []Record, groupFunc func(time.Time) string) map[string]AggregatedRecord {
	aggregations := make(map[string]AggregatedRecord)

	var lastOutTime time.Time
	for _, record := range records {
		if record.Kind == "out" {
			lastOutTime = record.Timestamp
		} else if record.Kind == "in" && !lastOutTime.IsZero() {
			groupKey := groupFunc(record.Timestamp)
			duration := lastOutTime.Sub(record.Timestamp).Hours()

			if agg, exists := aggregations[groupKey]; exists {
				agg.TotalHours += duration
				agg.Dates = append(agg.Dates, record.Timestamp.Format(DateFormat))
				agg.Notes = append(agg.Notes, record.Notes)
				aggregations[groupKey] = agg
			} else {
				aggregations[groupKey] = AggregatedRecord{
					Group:      groupKey,
					TotalHours: duration,
					Dates:      []string{record.Timestamp.Format(DateFormat)},
					Notes:      []string{record.Notes},
				}
			}
			lastOutTime = time.Time{} // reset
		}
	}

	return aggregations
}

// inferLastOut adds an "out" record at the beginning of the records if the last record is "in".
func inferLastOut(records *[]Record) int {
	if len(*records) > 0 && (*records)[0].Kind == "in" {
		record := []Record{
			{
				Timestamp: time.Now(),
				Kind:      "out",
				Notes:     "Inferred by takt.",
			},
		}
		*records = append(record, *records...)
		return 1
	}
	return 0
}

// printRecords prints the records.
func printRecords(records []Record) {
	fmt.Printf("%-25s %-5s %s\n", Header[0], Header[1], Header[2])
	for _, record := range records {
		fmt.Printf("%-25s %-5s %s\n", record.Timestamp.Format(TimeFormat), record.Kind, record.Notes)
	}
}

// createFile creates a new file with the header using the configured filename.
func createFile() error {
	if config == nil {
		return errors.New("config not initialized")
	}

	file, err := os.Create(config.FileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(Header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	return nil
}

// readRecords reads nrows records from the configured file
func readRecords(head int) ([]Record, error) {
	if config == nil {
		return nil, errors.New("config not initialized")
	}
	return readRecordsFromFile(config.FileName, head)
}

// readRecordsFromFile reads nrows records from the file fileName and returns them.
func readRecordsFromFile(fileName string, head int) ([]Record, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if err := createFile(); err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
	}

	// Create backup before reading
	if err := backupFile(fileName); err != nil {
		return nil, fmt.Errorf("could not create backup: %w", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		// Try to recover from backup
		if records, err := recoverFromBackup(fileName); err == nil {
			return records, nil
		}
		return nil, fmt.Errorf("could not read CSV: %w", err)
	}

	var records []Record
	var validRecords []Record
	var invalidLines []int

	if head == 0 || len(lines) < 2 {
		return records, nil
	}

	// Process all lines except header
	for i, line := range lines[1:] {
		if len(line) != 3 {
			invalidLines = append(invalidLines, i+1)
			continue
		}

		timestamp, err := time.Parse(TimeFormat, line[0])
		if err != nil {
			invalidLines = append(invalidLines, i+1)
			continue
		}

		record := Record{timestamp, line[1], line[2]}
		if err := validateRecord(record); err != nil {
			invalidLines = append(invalidLines, i+1)
			continue
		}

		validRecords = append(validRecords, record)
	}

	if len(invalidLines) > 0 {
		// Log warning about invalid lines
		log.Printf("Warning: found %d invalid records at lines: %v", len(invalidLines), invalidLines)

		// Write only valid records back to file
		if err := writeValidRecords(fileName, validRecords); err != nil {
			log.Printf("Error: could not clean up invalid records: %v", err)
		}
	}

	if head > 0 && len(validRecords) > head {
		return validRecords[:head], nil
	}
	return validRecords, nil
}

// checkAction checks in or out.
func checkAction(filename, notes string) error {
	records, err := readRecordsFromFile(filename, 1)
	if err != nil {
		return fmt.Errorf("failed to read records: %w", err)
	}

	var kind string
	if len(records) == 0 || records[0].Kind == "out" {
		kind = "in"
	} else {
		kind = "out"
	}

	timestamp := time.Now().Format(TimeFormat)
	line := fmt.Sprintf("%s,%s,%s", timestamp, kind, notes)
	if err := writeRecords(filename, line); err != nil {
		return fmt.Errorf("failed to write records: %w", err)
	}

	fmt.Printf("Check %s at %s\n", kind, timestamp)
	return nil
}

// validateRecord checks if a record is valid and returns an error if not
func validateRecord(record Record) error {
	if record.Timestamp.IsZero() {
		return fmt.Errorf("invalid timestamp")
	}
	if record.Kind != "in" && record.Kind != "out" {
		return fmt.Errorf("invalid kind: %s (must be 'in' or 'out')", record.Kind)
	}
	if record.Timestamp.After(time.Now()) {
		return fmt.Errorf("timestamp in future: %v", record.Timestamp)
	}
	return nil
}

// backupFile creates a backup of the file
func backupFile(fileName string) error {
	backupName := fileName + ".bak"
	source, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := source.Close(); err != nil {
			fmt.Printf("Error closing source file: %v\n", err)
		}
	}()

	destination, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer func() {
		if err := destination.Close(); err != nil {
			fmt.Printf("Error closing destination file: %v\n", err)
		}
	}()

	_, err = io.Copy(destination, source)
	return err
}

// recoverFromBackup attempts to recover records from the backup file
func recoverFromBackup(fileName string) ([]Record, error) {
	backupName := fileName + ".bak"
	records, err := readRecordsFromFile(backupName, -1)
	if err != nil {
		return nil, fmt.Errorf("could not recover from backup: %w", err)
	}
	return records, nil
}

// writeValidRecords writes only valid records back to the file
func writeValidRecords(fileName string, records []Record) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(Header); err != nil {
		return err
	}

	// Write records
	for _, record := range records {
		line := []string{
			record.Timestamp.Format(TimeFormat),
			record.Kind,
			record.Notes,
		}
		if err := writer.Write(line); err != nil {
			return err
		}
	}

	return nil
}

// writeRecords writes a new line to the file.
func writeRecords(fileName, newLine string) error {
	prevFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := prevFile.Close(); err != nil {
			fmt.Printf("Error closing previous file: %v\n", err)
		}
	}()

	newFile, err := os.CreateTemp("", "takt_tempfile.csv")
	if err != nil {
		fmt.Printf("Error: could not create temp file")
		return err
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			fmt.Printf("Error closing new file: %v\n", err)
		}
	}()

	newWriter := bufio.NewWriter(newFile)
	defer func() {
		if err := newWriter.Flush(); err != nil {
			fmt.Printf("Error flushing writer: %v\n", err)
		}
	}()
	_, err = fmt.Fprintf(newWriter, "%s,%s,%s\n", Header[0], Header[1], Header[2])
	if err != nil {
		fmt.Printf("Error: could not write to temp file")
		return err
	}
	_, err = newWriter.WriteString(newLine + "\n")
	if err != nil {
		fmt.Printf("Error: could not write to temp file")
		return err
	}

	prevReader := bufio.NewReader(prevFile)

	// drop the header
	_, _, err = prevReader.ReadLine()
	if err != nil {
		return err
	}
	_, err = io.Copy(newWriter, prevReader)
	if err != nil {
		return err
	}

	if err := os.Rename(newFile.Name(), fileName); err != nil {
		return err
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "takt [COMMAND] [ARGS]",
	Short: "CLI Time Tracking Tool",
	Long: `Takt is a simple time tracking tool that allows you to check in and out.

CONFIGURATION:
  Set these environment variables to customize behavior:
  - TAKT_FILE: Path to CSV file (default: ~/takt.csv)
  - TAKT_TARGET_HOURS: Target hours per day (default: 8.0)
  - TAKT_EDITOR: Editor for 'takt edit' command

EXAMPLES:
  # Check in/out (toggles automatically)
  takt check
  takt check "Working on project X"

  # View recent records
  takt cat 5                    # Show last 5 records

  # Daily summary with balance calculation
  takt day 10                   # Show last 10 days
  Output: Date         Total   Days  Avg     Balance
          2025-01-09   8h30m   1     8h30m   +0h30m
          2025-01-08   16h00m  1     16h00m  +1d

  # Weekly/monthly summaries
  takt week 4                   # Show last 4 weeks
  takt month 6                  # Show last 6 months

  # Visual grid for the year
  takt grid 2025 true          # Show 2025 with legend

BALANCE CALCULATION:
  The Balance column shows overtime/undertime based on your TARGET_HOURS:
  - +1d = 1 full working day of overtime
  - +1d2h = 1 day + 2 hours overtime
  - -0h30m = 30 minutes undertime
  - With 8h target: 16h worked = +1d balance
  - With 7.5h target: 16h worked = +1d1h balance

OUTPUT FORMAT:
  - Date: Date or period (2025-01-09, 2025-W02, 2025-01)
  - Total: Total hours worked in period
  - Days: Number of working days in period
  - Avg: Average hours per working day
  - Balance: Overtime/undertime vs target (±days/hours)`,
}

var checkCmd = &cobra.Command{
	Aliases: []string{"c"},
	Use:     "check [NOTE]",
	Short:   "Check in or out (toggles automatically)",
	Long: `Check in or out. The command automatically toggles between 'in' and 'out' states.
If you're currently checked out, it will check you in.
If you're currently checked in, it will check you out.

EXAMPLES:
  takt check                    # Simple check in/out
  takt check "Meeting prep"     # Check in/out with note
  takt c "Lunch break"          # Using alias

OUTPUT:
  Check in at 2025-01-09T14:30:00Z
  Check out at 2025-01-09T17:45:00Z`,
	Run: func(cmd *cobra.Command, args []string) {
		if config == nil {
			fmt.Println("Error: config not initialized")
			return
		}

		notes := ""
		if len(args) > 0 {
			notes = args[0]
		}

		if err := checkAction(config.FileName, notes); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

var catCmd = &cobra.Command{
	Aliases: []string{"display"},
	Use:     "cat [HEAD]",
	Short:   "Show recent records",
	Long: `Show recent time tracking records in chronological order.
Default shows the last 10 records. Use HEAD to specify a different number.

EXAMPLES:
  takt cat                      # Show last 10 records
  takt cat 20                   # Show last 20 records
  takt display 5                # Using alias

OUTPUT FORMAT:
  timestamp                 kind  notes
  2025-01-09T14:30:00Z     in    Meeting prep
  2025-01-09T17:45:00Z     out   End of day`,
	Run: func(cmd *cobra.Command, args []string) {
		head := DefaultHead
		var err error
		if len(args) > 0 {
			head, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}
		records, err := readRecords(head)
		if err != nil {
			log.Fatal(err)
		}
		printRecords(records)
	},
}

var dayCmd = &cobra.Command{
	Aliases: []string{"d"},
	Use:     "day [HEAD]",
	Short:   "Daily summary with balance calculation",
	Long: `Show daily time tracking summary with balance calculation.
Balance shows overtime/undertime based on your TARGET_HOURS setting.
Default shows the last 10 days.

EXAMPLES:
  takt day                      # Show last 10 days
  takt day 30                   # Show last 30 days
  takt d 5                      # Using alias

OUTPUT FORMAT:
  Date         Total   Days  Avg     Balance
  2025-01-09   8h30m   1     8h30m   +0h30m
  2025-01-08   16h00m  1     16h00m  +1d

BALANCE EXPLANATION:
  - +1d = 1 full working day of overtime (based on TARGET_HOURS)
  - +0h30m = 30 minutes overtime
  - -2h00m = 2 hours undertime`,
	Run: func(cmd *cobra.Command, args []string) {
		head := DefaultHead
		var err error
		if len(args) > 0 {
			head, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}
		summary("day", head)
	},
}

var weekCmd = &cobra.Command{
	Aliases: []string{"w"},
	Use:     "week [HEAD]",
	Short:   "Weekly summary with balance calculation",
	Long: `Show weekly time tracking summary with balance calculation.
Weeks are calculated using ISO week numbers (Monday to Sunday).
Balance shows overtime/undertime based on TARGET_HOURS × working days.

EXAMPLES:
  takt week                     # Show last 10 weeks
  takt week 4                   # Show last 4 weeks
  takt w 12                     # Using alias

OUTPUT FORMAT:
  Date      Total     Days  Avg     Balance
  2025-W02  40h15m    5     8h03m   +0h15m
  2025-W01  37h30m    5     7h30m   -2h30m`,
	Run: func(cmd *cobra.Command, args []string) {
		head := DefaultHead
		var err error
		if len(args) > 0 {
			head, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}
		summary("week", head)
	},
}

var monthCmd = &cobra.Command{
	Aliases: []string{"m"},
	Use:     "month [HEAD]",
	Short:   "Monthly summary with balance calculation",
	Long: `Show monthly time tracking summary with balance calculation.
Balance shows overtime/undertime based on TARGET_HOURS × working days.
Use -1 to show all months.

EXAMPLES:
  takt month                    # Show last 10 months
  takt month 6                  # Show last 6 months
  takt month -1                 # Show all months
  takt m 3                      # Using alias

OUTPUT FORMAT:
  Date     Total     Days  Avg     Balance
  2025-01  168h30m   21    8h02m   +0h30m
  2024-12  159h45m   20    7h59m   -0h15m`,
	Run: func(cmd *cobra.Command, args []string) {
		head := DefaultHead
		var err error
		if len(args) > 0 {
			head, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}
		summary("month", head)
	},
}

var yearCmd = &cobra.Command{
	Aliases: []string{"y"},
	Use:     "year [HEAD]",
	Short:   "Yearly summary with balance calculation",
	Long: `Show yearly time tracking summary with balance calculation.
Balance shows overtime/undertime based on TARGET_HOURS × working days.

EXAMPLES:
  takt year                     # Show last 10 years
  takt year 3                   # Show last 3 years
  takt y 5                      # Using alias

OUTPUT FORMAT:
  Date  Total      Days  Avg     Balance
  2025  2080h30m   260   8h00m   +0h30m
  2024  2076h15m   259   8h01m   +4h15m`,
	Run: func(cmd *cobra.Command, args []string) {
		head := DefaultHead
		var err error
		if len(args) > 0 {
			head, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}
		summary("year", head)
	},
}

var gridCmd = &cobra.Command{
	Short: "Visual grid showing daily activity with colors",
	Use:   "grid [YEAR] [LEGEND]",
	Long: `Display a visual grid showing daily activity for the specified year.
Each day is represented by a colored symbol indicating hours worked.
Features month separators, color coding, and activity statistics.
Default shows current year without legend.

EXAMPLES:
  takt grid                     # Show current year
  takt grid 2024                # Show 2024
  takt grid 2025 true           # Show 2025 with legend

GRID SYMBOLS:
  󰋣  = 0-1 hours (minimal work) - Gray
  ▪  = 1-4 hours (light work) - Yellow
  ▮  = 4-8 hours (normal work) - Green
  󰈸 = 8-12 hours (heavy work) - Blue
  █  = 12+ hours (very heavy work) - Purple

FEATURES:
  • Color-coded activity levels
  • Month separators for better readability
  • Activity statistics (total days, active days, activity rate)
  • Improved legend with descriptions
  • Better visual alignment and formatting

OUTPUT FORMAT:
      Date         W  L  M  X  J  V  S  D
      ═══════════════════════════════════
      2025-01-01 01 󰋣       󰈸
      2025-01-06 02     󰈸 󰈸   󰋣
      ───────────────────────────────
      📊 Summary:
      ├─ Total tracked days: 192
      ├─ Active work days: 116
      └─ Activity rate: 60.4%`,
	Run: func(cmd *cobra.Command, args []string) {
		lenArgs := len(args)
		legend := false
		year := ""

		if lenArgs < 1 {
			year = time.Now().Format("2006")
		} else {
			year = args[0]
		}

		if lenArgs > 1 {
			legend = args[1] == "true"
		}

		if err := printGrid(year, legend); err != nil {
			log.Fatalf("Failed to print grid: %v", err)
		}
	},
}

var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit the records file",
	Long: `Open the time tracking CSV file in your configured editor.
Set TAKT_EDITOR environment variable to specify your preferred editor.

EXAMPLES:
  takt edit                     # Open in configured editor
  takt e                        # Using alias

SETUP:
  export TAKT_EDITOR=vim        # Use vim
  export TAKT_EDITOR=code       # Use VS Code
  export TAKT_EDITOR=nano       # Use nano

FILE FORMAT:
  timestamp,kind,notes
  2025-01-09T14:30:00Z,in,Meeting prep
  2025-01-09T17:45:00Z,out,End of day`,
	Run: func(cmd *cobra.Command, args []string) {
		if config == nil {
			fmt.Println("Error: config not initialized")
			return
		}

		if config.Editor == "" {
			fmt.Println("Error: TAKT_EDITOR environment variable not set")
			return
		}

		editCmd := exec.Command(config.Editor, config.FileName)
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		err := editCmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var commitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"cm"},
	Short:   "Commit the records file",
	Run: func(cmd *cobra.Command, args []string) {
		err := gitAdd()
		if err != nil {
			fmt.Println("Error: git add failed")
			return
		}
		fmt.Println("Records added")

		err = gitCommit()
		if err != nil {
			fmt.Println("Error: git commit failed")
			return
		}
		fmt.Println("Records committed")

		err = gitPush()
		if err != nil {
			fmt.Println("Error: git push failed")
			return
		}
		fmt.Println("Records pushed")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information including Go version and build target.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("takt version %s\n", Version)
		fmt.Printf("Built with Go %s\n", runtime.Version())
		fmt.Printf("Target: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(catCmd)
	rootCmd.AddCommand(dayCmd)
	rootCmd.AddCommand(weekCmd)
	rootCmd.AddCommand(monthCmd)
	rootCmd.AddCommand(yearCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(gridCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	Execute()
}
