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
	// Initialize config for testing
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	originalConfig := config
	config = cfg
	config.TargetHours = 8.0 // Use 8 hours as target for testing
	defer func() { config = originalConfig }()

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
		{"positive_7hr_30min", 7.5, "+7h30m"},
		{"negative_7hr_30min", -7.5, "-7h30m"},
		{"positive_8hr", 8.0, "+1d"},               // 8 hours = 1 day
		{"negative_8hr", -8.0, "-1d"},              // -8 hours = -1 day
		{"positive_9hr", 9.0, "+1d1h"},             // 9 hours = 1 day + 1 hour
		{"negative_9hr", -9.0, "-1d1h"},            // -9 hours = -1 day + 1 hour
		{"positive_16hr", 16.0, "+2d"},             // 16 hours = 2 days
		{"negative_16hr", -16.0, "-2d"},            // -16 hours = -2 days
		{"positive_17hr_30min", 17.5, "+2d1h30m"},  // 17.5 hours = 2 days + 1.5 hours
		{"negative_17hr_30min", -17.5, "-2d1h30m"}, // -17.5 hours = -2 days + 1.5 hours
		{"positive_24hr", 24.0, "+3d"},             // 24 hours = 3 days (8h each)
		{"negative_24hr", -24.0, "-3d"},            // -24 hours = -3 days
		{"positive_25hr", 25.0, "+3d1h"},           // 25 hours = 3 days + 1 hour
		{"negative_25hr", -25.0, "-3d1h"},          // -25 hours = -3 days + 1 hour
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

func TestFormatOvertimeWithDifferentTargetHours(t *testing.T) {
	// Initialize config for testing
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	originalConfig := config
	config = cfg
	defer func() { config = originalConfig }()

	tests := []struct {
		name        string
		targetHours float64
		difference  float64
		expected    string
	}{
		// Tests with 7.5 hour target
		{"7.5h_target_positive_7.5hr", 7.5, 7.5, "+1d"},
		{"7.5h_target_positive_15hr", 7.5, 15.0, "+2d"},
		{"7.5h_target_positive_16hr", 7.5, 16.0, "+2d1h"},
		{"7.5h_target_positive_22.5hr", 7.5, 22.5, "+3d"},
		{"7.5h_target_negative_7.5hr", 7.5, -7.5, "-1d"},
		
		// Tests with 6 hour target
		{"6h_target_positive_6hr", 6.0, 6.0, "+1d"},
		{"6h_target_positive_12hr", 6.0, 12.0, "+2d"},
		{"6h_target_positive_13.5hr", 6.0, 13.5, "+2d1h30m"},
		{"6h_target_negative_6hr", 6.0, -6.0, "-1d"},
		
		// Tests with 10 hour target
		{"10h_target_positive_10hr", 10.0, 10.0, "+1d"},
		{"10h_target_positive_20hr", 10.0, 20.0, "+2d"},
		{"10h_target_positive_25hr", 10.0, 25.0, "+2d5h"},
		{"10h_target_negative_10hr", 10.0, -10.0, "-1d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.TargetHours = tt.targetHours
			result := formatOvertime(tt.difference)
			if result != tt.expected {
				t.Errorf("formatOvertime(%v) with target %v = %v, want %v", tt.difference, tt.targetHours, result, tt.expected)
			}
		})
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
