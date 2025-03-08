package main

import (
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
	now := time.Now()
	records := []Record{
		{now.Add(-23 * time.Hour), "out", "Note1"},
		{now.Add(-24 * time.Hour), "in", "Note1"},
	}

	labeler := func(t time.Time) string {
		return t.Format("2006-01-02")
	}

	got := aggregateBy(records, labeler)

	if len(got) != 1 {
		t.Errorf("Expected 1 aggregation record, got %d", len(got))
	}
	for k, v := range got {
		if v.TotalHours != 1 {
			t.Errorf("For key %s, expected 1 hours, got %.2f \n %+v", k, v.TotalHours, got)
		}
	}
}

func TestInferLastOut(t *testing.T) {
	records := []Record{
		{time.Now().Add(-2 * time.Hour), "in", ""},
	}

	n := inferLastOut(&records)

	if n != 1 {
		t.Errorf("Expected to infer 1 'out' record, inferred %d", n)
	}

	if records[0].Kind != "out" {
		t.Errorf("Expected first record to be 'out', got '%s'", records[0].Kind)
	}
}

func TestReadRecords(t *testing.T) {
	dummyCSV := `timestamp,kind,notes
2020-02-01T00:00:00Z,in,Note2
2020-01-01T01:00:00Z,out,Note1
2020-01-01T00:00:00Z,in,Note1
`
	// Create test file
	err := os.WriteFile(FileName, []byte(dummyCSV), 0644)
	if err != nil {
		t.Fatalf("Failed to create file for testing: %v", err)
	}
	defer os.Remove(FileName) // clean up

	records, err := readRecords(-1)
	if err != nil {
		t.Fatalf("Failed to read records: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}
}

// TODO: improve this
func TestValidateRecord(t *testing.T) {
    now := time.Now()
    future := now.Add(24 * time.Hour)
    
    tests := []struct {
        name    string
        record  Record
        wantErr bool
    }{
        {
            name: "valid record",
            record: Record{
                Timestamp: now,
                Kind:     "in",
                Notes:    "test",
            },
            wantErr: false,
        },
        {
            name: "zero timestamp",
            record: Record{
                Timestamp: time.Time{},
                Kind:     "in",
                Notes:    "test",
            },
            wantErr: true,
        },
        {
            name: "invalid kind",
            record: Record{
                Timestamp: now,
                Kind:     "invalid",
                Notes:    "test",
            },
            wantErr: true,
        },
        {
            name: "future timestamp",
            record: Record{
                Timestamp: future,
                Kind:     "in",
                Notes:    "test",
            },
            wantErr: true,
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
    // Create a temporary test file
    testContent := "timestamp,kind,notes\n2023-01-01T12:00:00Z,in,test note"
    tmpFile, err := os.CreateTemp("", "takt-test-*.csv")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())
    
    if err := os.WriteFile(tmpFile.Name(), []byte(testContent), 0644); err != nil {
        t.Fatalf("Failed to write test content: %v", err)
    }

    // Test backup creation
    if err := backupFile(tmpFile.Name()); err != nil {
        t.Errorf("backupFile() failed: %v", err)
    }

    // Verify backup file exists
    backupName := tmpFile.Name() + ".bak"
    if _, err := os.Stat(backupName); os.IsNotExist(err) {
        t.Errorf("Backup file was not created")
    }
    defer os.Remove(backupName)

    // Test recovery
    records, err := recoverFromBackup(tmpFile.Name())
    if err != nil {
        t.Errorf("recoverFromBackup() failed: %v", err)
    }
    if len(records) != 1 {
        t.Errorf("Expected 1 record from backup, got %d", len(records))
    }
}

func TestWriteValidRecords(t *testing.T) {
    // Create test records
    now := time.Now()
    records := []Record{
        {now, "in", "test note 1"},
        {now.Add(time.Hour), "out", "test note 2"},
    }

    // Create temporary file
    tmpFile, err := os.CreateTemp("", "takt-test-*.csv")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    // Test writing valid records
    if err := writeValidRecords(tmpFile.Name(), records); err != nil {
        t.Errorf("writeValidRecords() failed: %v", err)
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
	// Create a temporary file to simulate the CSV records file
	// You may need to specify a unique temp file for concurrent tests.
	tempFile, err := os.CreateTemp("", "test_checkAction.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // clean up

	// Write initial record to the temp file
	csvContent := "timestamp,kind,notes\n"
	if _, err := tempFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	checkAction(tempFile.Name(), "Test Note")

	// Read the modified file content
	modifiedFile, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to open modified temp file: %v", err)
	}
	defer modifiedFile.Close()

	var _gotCsvContent []byte
	if _gotCsvContent, err = os.ReadFile(modifiedFile.Name()); err != nil {
		t.Fatalf("Failed to read modified temp file: %v", err)
	}
	gotCsvContentList := strings.Split(string(_gotCsvContent), "\n")

	if len(gotCsvContentList) == len(csvContent) {
		t.Errorf("Expected modified file content, but no changes detected \nGot:\n%+v\n\nExpected:\n%+v", gotCsvContentList, csvContent)
	}
}
