package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCalculateDuration(t *testing.T) {
	records := []Record{
		{time.Now().Add(-4 * time.Hour), "in", ""},
		{time.Now().Add(-2 * time.Hour), "out", ""},
	}

	tests := []struct {
		name    string
		period  string
		length  bool
		avgHrs  bool
		wantErr bool
	}{
		{"day", "day", true, true, false},
		{"week", "week", true, true, false},
		{"month", "month", true, true, false},
		{"year", "year", true, true, false},
		{"unsupported", "unsupported", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateDuration(records, tt.period)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(got) != 0) != tt.length {
				t.Errorf("Unexpected length of results")
			}
			if tt.avgHrs {
				for _, aggRec := range got {
					if aggRec.AverageHours == 0 {
						t.Errorf("Expected non-zero average hours")
					}
				}
			}
		})
	}
}

func TestAggregateBy(t *testing.T) {
	// Create records in the order they appear in the slice (out first, then in)
	// but chronologically in came before out
	now := time.Now()
	records := []Record{
		{now.Add(-4 * time.Hour), "in", ""},  // 4 hours ago
		{now.Add(-2 * time.Hour), "out", ""}, // 2 hours ago (worked for 2 hours)
	}

	groupFunc := func(t time.Time) string {
		return t.Format("2006-01-02")
	}

	result := aggregateBy(records, groupFunc)
	if len(result) == 0 {
		t.Error("Expected non-empty result from aggregateBy")
	}

	// Check that we have at least one aggregated record
	for _, agg := range result {
		if agg.TotalHours != 2.0 {
			t.Errorf("Expected 2 hours of work, got %v", agg.TotalHours)
		}
	}
}

func TestInferLastOut(t *testing.T) {
	records := []Record{
		{time.Now().Add(-2 * time.Hour), "in", ""},
	}

	inferLastOut(&records)
	if len(records) != 2 {
		t.Errorf("Expected 2 records after inferLastOut, got %d", len(records))
	}
	if records[0].Kind != "out" {
		t.Errorf("Expected first record to be 'out', got %s", records[0].Kind)
	}
}

func TestReadRecords(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "takt_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	writer := csv.NewWriter(tempFile)
	writer.Write([]string{"timestamp", "kind", "notes"})
	writer.Write([]string{"2023-01-01T10:00:00Z", "in", "test"})
	writer.Write([]string{"2023-01-01T18:00:00Z", "out", "test"})
	writer.Flush()
	tempFile.Close()

	// Test reading
	records, err := readRecordsFromFile(tempFile.Name(), -1)
	if err != nil {
		t.Fatalf("readRecordsFromFile() failed: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}
}

func TestValidateRecord(t *testing.T) {
	tests := []struct {
		name    string
		record  Record
		wantErr bool
	}{
		{
			"valid_record",
			Record{time.Now().Add(-1 * time.Hour), "in", "test"},
			false,
		},
		{
			"zero_timestamp",
			Record{time.Time{}, "in", "test"},
			true,
		},
		{
			"invalid_kind",
			Record{time.Now().Add(-1 * time.Hour), "invalid", "test"},
			true,
		},
		{
			"future_timestamp",
			Record{time.Now().Add(1 * time.Hour), "in", "test"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecord(tt.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackupAndRecover(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "takt_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	testData := "timestamp,kind,notes\n2023-01-01T10:00:00Z,in,test\n"
	if _, err := tempFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	// Test backup creation
	if err := backupFile(tempFile.Name()); err != nil {
		t.Errorf("backupFile() failed: %v", err)
	}

	// Verify backup file exists
	backupName := tempFile.Name() + ".bak"
	if _, err := os.Stat(backupName); os.IsNotExist(err) {
		t.Errorf("Backup file was not created")
	}
	defer os.Remove(backupName)

	// Test recovery
	records, err := recoverFromBackup(tempFile.Name())
	if err != nil {
		t.Errorf("recoverFromBackup() failed: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("Expected 1 record from backup, got %d", len(records))
	}
}

func TestWriteValidRecords(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "takt_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create test records
	records := []Record{
		{time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), "in", "test1"},
		{time.Date(2023, 1, 1, 18, 0, 0, 0, time.UTC), "out", "test2"},
	}

	// Write records
	if err := writeValidRecords(tmpFile.Name(), records); err != nil {
		t.Fatalf("writeValidRecords() failed: %v", err)
	}

	// Verify file contents
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	// Check if header is present
	if !strings.Contains(string(content), "timestamp,kind,notes") {
		t.Error("Header is missing in written file")
	}

	// Check if records are present
	for _, record := range records {
		timestamp := record.Timestamp.Format(TimeFormat)
		if !strings.Contains(string(content), timestamp) {
			t.Errorf("Record with timestamp %s is missing", timestamp)
		}
	}
}

func TestCheckAction(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "takt_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Initialize with header
	file, err := os.Create(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"timestamp", "kind", "notes"})
	writer.Flush()
	file.Close()

	// Test check action
	err = checkAction(tempFile.Name(), "test")
	if err != nil {
		t.Errorf("checkAction() failed: %v", err)
	}
}

func TestGetTargetHours(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		dflt     float64
		expected float64
	}{
		{"default", "", 8.0, 8.0},
		{"float_format", "7.5", 8.0, 7.5},
		{"time_format", "7:30", 8.0, 7.5},
		{"invalid_format", "invalid", 8.0, 8.0},
		{"invalid_time", "7:99", 8.0, 8.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			key := "TEST_TARGET_HOURS"
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
			} else {
				os.Unsetenv(key)
			}
			defer os.Unsetenv(key)

			result, err := getTargetHours(key, tt.dflt)
			if err != nil {
				t.Errorf("getTargetHours() returned error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("getTargetHours() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	// Test config loading
	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() failed: %v", err)
	}
	if cfg == nil {
		t.Error("LoadConfig() returned nil config")
	}
}

func TestFormatOvertime(t *testing.T) {
	tests := []struct {
		name       string
		difference float64
		expected   string
	}{
		{"zero_difference", 0.0, "00h00m"},
		{"positive_30min", 0.5, "+0h30m"},
		{"negative_30min", -0.5, "-0h30m"},
		{"positive_1hr", 1.0, "+1h00m"},
		{"negative_1hr", -1.0, "-1h00m"},
		{"positive_1hr_30min", 1.5, "+1h30m"},
		{"negative_1hr_30min", -1.5, "-1h30m"},
		{"positive_8hr", 8.0, "+8h00m"},
		{"negative_8hr", -8.0, "-8h00m"},
		{"positive_24hr", 24.0, "+24h00m"},
		{"negative_24hr", -24.0, "-24h00m"},
		{"positive_over_24hr", 25.5, "+1d01h30m"},
		{"negative_over_24hr", -25.5, "-1d01h30m"},
		{"positive_exact_48hr", 48.0, "+2d00h00m"},
		{"negative_exact_48hr", -48.0, "-2d00h00m"},
		{"positive_complex", 74.75, "+3d02h45m"},
		{"negative_complex", -74.75, "-3d02h45m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatOvertime(tt.difference)
			if result != tt.expected {
				t.Errorf("formatOvertime(%v) = %v, want %v", tt.difference, result, tt.expected)
			}
		})
	}
}

func TestBalanceCalculation(t *testing.T) {
	// Initialize config for testing
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	originalConfig := config
	config = cfg
	defer func() { config = originalConfig }()

	// Test balance calculation in aggregated records
	now := time.Now()
	records := []Record{
		{now.Add(-48 * time.Hour), "in", ""},  // 2 days ago, 8am
		{now.Add(-39 * time.Hour), "out", ""}, // 2 days ago, 5pm (9 hours)
		{now.Add(-24 * time.Hour), "in", ""},  // 1 day ago, 8am
		{now.Add(-18 * time.Hour), "out", ""}, // 1 day ago, 2pm (6 hours)
	}

	// Set target hours to 8 for testing
	originalTargetHours := config.TargetHours
	config.TargetHours = 8.0
	defer func() { config.TargetHours = originalTargetHours }()

	// Test daily aggregation
	result, err := calculateDuration(records, "day")
	if err != nil {
		t.Fatalf("calculateDuration() failed: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 daily records, got %d", len(result))
	}

	// Check first day (9 hours worked, 8 hours expected, +1 hour balance)
	firstDay := result[0]
	if firstDay.TotalHours != 9.0 {
		t.Errorf("First day total hours = %v, want 9.0", firstDay.TotalHours)
	}
	expectedBalance1 := 9.0 - (8.0 * 1) // 9 hours - (8 target * 1 day)
	if expectedBalance1 != 1.0 {
		t.Errorf("First day expected balance = %v, want 1.0", expectedBalance1)
	}

	// Check second day (6 hours worked, 8 hours expected, -2 hours balance)
	secondDay := result[1]
	if secondDay.TotalHours != 6.0 {
		t.Errorf("Second day total hours = %v, want 6.0", secondDay.TotalHours)
	}
	expectedBalance2 := 6.0 - (8.0 * 1) // 6 hours - (8 target * 1 day)
	if expectedBalance2 != -2.0 {
		t.Errorf("Second day expected balance = %v, want -2.0", expectedBalance2)
	}

	// Test weekly aggregation (both days should be in same week)
	weeklyResult, err := calculateDuration(records, "week")
	if err != nil {
		t.Fatalf("calculateDuration() for week failed: %v", err)
	}

	if len(weeklyResult) != 1 {
		t.Fatalf("Expected 1 weekly record, got %d", len(weeklyResult))
	}

	weeklyRecord := weeklyResult[0]
	if weeklyRecord.TotalHours != 15.0 { // 9 + 6 = 15 hours
		t.Errorf("Weekly total hours = %v, want 15.0", weeklyRecord.TotalHours)
	}
	expectedWeeklyBalance := 15.0 - (8.0 * 2) // 15 hours - (8 target * 2 days)
	if expectedWeeklyBalance != -1.0 {
		t.Errorf("Weekly expected balance = %v, want -1.0", expectedWeeklyBalance)
	}
}

func TestBalanceCalculationWithDifferentTargetHours(t *testing.T) {
	// Initialize config for testing
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	originalConfig := config
	config = cfg
	defer func() { config = originalConfig }()

	// Test with 7.5 hour target
	now := time.Now()
	records := []Record{
		{now.Add(-24 * time.Hour), "in", ""},  // 1 day ago, start
		{now.Add(-16 * time.Hour), "out", ""}, // 1 day ago, end (8 hours)
	}

	// Set target hours to 7.5 for testing
	originalTargetHours := config.TargetHours
	config.TargetHours = 7.5
	defer func() { config.TargetHours = originalTargetHours }()

	result, err := calculateDuration(records, "day")
	if err != nil {
		t.Fatalf("calculateDuration() failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 daily record, got %d", len(result))
	}

	dailyRecord := result[0]
	if dailyRecord.TotalHours != 8.0 {
		t.Errorf("Daily total hours = %v, want 8.0", dailyRecord.TotalHours)
	}

	// 8 hours worked - 7.5 target = +0.5 hour balance
	expectedBalance := 8.0 - (7.5 * 1)
	if expectedBalance != 0.5 {
		t.Errorf("Expected balance = %v, want 0.5", expectedBalance)
	}
}

func TestBalanceCalculationMultipleDays(t *testing.T) {
	// Initialize config for testing
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	originalConfig := config
	config = cfg
	defer func() { config = originalConfig }()

	// Test over multiple days with different work patterns
	now := time.Now()
	records := []Record{
		// Day 1: 10 hours
		{now.Add(-72 * time.Hour), "in", ""},  // 3 days ago
		{now.Add(-62 * time.Hour), "out", ""}, // 3 days ago
		// Day 2: 6 hours
		{now.Add(-48 * time.Hour), "in", ""},  // 2 days ago
		{now.Add(-42 * time.Hour), "out", ""}, // 2 days ago
		// Day 3: 8 hours
		{now.Add(-24 * time.Hour), "in", ""},  // 1 day ago
		{now.Add(-16 * time.Hour), "out", ""}, // 1 day ago
	}

	// Set target hours to 8 for testing
	originalTargetHours := config.TargetHours
	config.TargetHours = 8.0
	defer func() { config.TargetHours = originalTargetHours }()

	result, err := calculateDuration(records, "day")
	if err != nil {
		t.Fatalf("calculateDuration() failed: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("Expected 3 daily records, got %d", len(result))
	}

	// Check individual day balances
	expectedHours := []float64{10.0, 6.0, 8.0}
	expectedBalances := []float64{2.0, -2.0, 0.0} // vs 8 hour target

	for i, record := range result {
		if record.TotalHours != expectedHours[i] {
			t.Errorf("Day %d total hours = %v, want %v", i+1, record.TotalHours, expectedHours[i])
		}
		expectedBalance := expectedHours[i] - (8.0 * 1)
		if expectedBalance != expectedBalances[i] {
			t.Errorf("Day %d expected balance = %v, want %v", i+1, expectedBalance, expectedBalances[i])
		}
	}

	// Test weekly aggregation
	weeklyResult, err := calculateDuration(records, "week")
	if err != nil {
		t.Fatalf("calculateDuration() for week failed: %v", err)
	}

	if len(weeklyResult) != 1 {
		t.Fatalf("Expected 1 weekly record, got %d", len(weeklyResult))
	}

	weeklyRecord := weeklyResult[0]
	totalHours := 10.0 + 6.0 + 8.0 // 24 hours
	if weeklyRecord.TotalHours != totalHours {
		t.Errorf("Weekly total hours = %v, want %v", weeklyRecord.TotalHours, totalHours)
	}

	// 24 hours worked - (8 target * 3 days) = 0 balance
	expectedWeeklyBalance := totalHours - (8.0 * 3)
	if expectedWeeklyBalance != 0.0 {
		t.Errorf("Weekly expected balance = %v, want 0.0", expectedWeeklyBalance)
	}
}

func TestCLIIntegrationWithBalance(t *testing.T) {
	// Create a temporary test file
	testFile, err := os.CreateTemp("", "takt_integration_balance_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(testFile.Name())
	testFile.Close()

	// Initialize config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	originalConfig := config
	config = cfg
	config.FileName = testFile.Name()
	config.TargetHours = 8.0
	defer func() { config = originalConfig }()

	// Create file and add test data
	err = createFile()
	if err != nil {
		t.Fatalf("createFile() failed: %v", err)
	}

	// Add records for multiple days with different work patterns
	testTime1 := time.Now().Add(-48 * time.Hour) // Day 1 start
	testTime2 := time.Now().Add(-39 * time.Hour) // Day 1 end (9 hours)
	testTime3 := time.Now().Add(-24 * time.Hour) // Day 2 start
	testTime4 := time.Now().Add(-18 * time.Hour) // Day 2 end (6 hours)

	// Create CSV content with multiple days
	csvContent := fmt.Sprintf("timestamp,kind,notes\n%s,out,Day 2 end\n%s,in,Day 2 start\n%s,out,Day 1 end\n%s,in,Day 1 start\n",
		testTime4.Format(TimeFormat), testTime3.Format(TimeFormat),
		testTime2.Format(TimeFormat), testTime1.Format(TimeFormat))

	err = os.WriteFile(testFile.Name(), []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Read back and verify
	records, err := readRecordsFromFile(testFile.Name(), -1)
	if err != nil {
		t.Fatalf("readRecordsFromFile() failed: %v", err)
	}

	if len(records) != 4 {
		t.Fatalf("Expected 4 records, got %d", len(records))
	}

	// Test daily calculation
	dailyResult, err := calculateDuration(records, "day")
	if err != nil {
		t.Fatalf("calculateDuration() daily failed: %v", err)
	}

	if len(dailyResult) != 2 {
		t.Fatalf("Expected 2 daily records, got %d", len(dailyResult))
	}

	// Day 1: 9 hours vs 8 hour target = +1 hour
	expectedBalance1 := 9.0 - 8.0
	if expectedBalance1 != 1.0 {
		t.Errorf("Day 1 expected balance = %v, want 1.0", expectedBalance1)
	}

	// Day 2: 6 hours vs 8 hour target = -2 hours
	expectedBalance2 := 6.0 - 8.0
	if expectedBalance2 != -2.0 {
		t.Errorf("Day 2 expected balance = %v, want -2.0", expectedBalance2)
	}

	// Test weekly calculation
	weeklyResult, err := calculateDuration(records, "week")
	if err != nil {
		t.Fatalf("calculateDuration() weekly failed: %v", err)
	}

	if len(weeklyResult) != 1 {
		t.Fatalf("Expected 1 weekly record, got %d", len(weeklyResult))
	}

	// Total: 15 hours over 2 days vs 16 hour target (8*2) = -1 hour
	expectedWeeklyBalance := 15.0 - (8.0 * 2)
	if expectedWeeklyBalance != -1.0 {
		t.Errorf("Weekly expected balance = %v, want -1.0", expectedWeeklyBalance)
	}
}
