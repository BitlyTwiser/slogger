package slogger_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/BitlyTwiser/slogger"
)

// Just look at some output, this is nothing more than showing how the values may look
func TestSloggerOutput(t *testing.T) {
	stdoutLogger := slogger.NewLogger(os.Stdout)
	stdoutLogger.LogEvent("info", "Test one", map[string]any{"one": false})
	stdoutLogger.LogEvent("info", "Test two", "one", "two", "Another", false, "four", true, "bob")
	stdoutLogger.LogEvent("info", "Test three", "one", "two", "Another", false, "four", true)
	stdoutLogger.LogEvent("info", "Test four", "one", "two", "Another", false, "four", true, map[string]any{"one": false})
	stdoutLogger.LogEvent("info", "Test five", "key", "value", "AnotherKey", false, "four", 123123, map[string]any{"MORETEST": 123123, "TestAgain": false}, "more", 123123)
	stdoutLogger.LogEvent("warn", "Test six", "key", "value", "AnotherKey", false, "four", 123123, map[string]any{"testOne": 42069, "testTwo": false})
}

func TestSloggerError(t *testing.T) {
	stdoutLogger := slogger.NewLogger(os.Stdout)
	err := stdoutLogger.LogError("error", fmt.Errorf("something died")) // Test error logger
	if err == nil {
		t.Fatalf("expected error, got nil return")
	}

	err = stdoutLogger.LogError("error", fmt.Errorf("Test 2 died"), 1, 2, 3, "masdasd")

	if err == nil {
		t.Fatalf("expected error, got nil return")
	}
}

// Test logging JSON output to file
func TestSloggerToFile(t *testing.T) {
	testFile, err := filepath.Abs("log.json")
	if err != nil {
		return
	}

	f, err := os.OpenFile(testFile, os.O_RDWR|os.O_TRUNC, os.ModeAppend)

	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(testFile)
			if err != nil {
				t.Fatalf("could not write to log file. %v", err.Error())
			}

		}
	}

	defer func() {
		f.Close()
		err := os.Remove(testFile)
		if err != nil {
			if os.IsNotExist(err) {
				panic(fmt.Sprintf("test file was never created, the following error occurred: %v", err))
			}
			panic(fmt.Sprintf("error deleting test file, delete file manually. err: %v", err))
		}
	}()

	tests := []struct {
		name      string
		msg       string
		file      *os.File
		have      map[string]any
		want      map[string]any
		eventType string
		eventData []any
		clean     func(f *os.File)
	}{
		{
			name:      "General Test of Values",
			msg:       "A Test",
			file:      f,
			have:      map[string]any{"miscFields": "bob", "one": "two", "Another": false, "four": true},
			want:      make(map[string]any),
			eventType: "info",
			eventData: []any{"one", "two", "Another", false, "four", true, "bob"},
			clean:     wipeFileData,
		},
		{
			name:      "Embedded map test",
			msg:       "Another test",
			file:      f,
			have:      map[string]any{"key": "value", "AnotherKey": false, "four": float64(123123), "testOne": float64(42069), "testTwo": false},
			want:      make(map[string]any),
			eventType: "warn",
			eventData: []any{"key", "value", "AnotherKey", false, "four", 123123, map[string]any{"testOne": 42069, "testTwo": false}},
			clean:     wipeFileData,
		},
		{
			name:      "Many Values test",
			msg:       "a third test",
			file:      f,
			have:      map[string]any{"AnotherKey": false, "four": float64(123123), "more": "", "MORETEST": float64(123123), "TestAgain": false, "miscFields": float64(123123), "key": "value"},
			want:      make(map[string]any),
			eventType: "warn",
			eventData: []any{"key", "value", "AnotherKey", false, "four", 123123, map[string]any{"MORETEST": 123123, "TestAgain": false}, "more", 123123},
			clean:     wipeFileData,
		},
		{
			name:      "No Values test",
			msg:       "A fourth test",
			file:      f,
			have:      map[string]any{},
			want:      make(map[string]any),
			eventType: "info",
			eventData: []any{},
			clean:     wipeFileData,
		},
	}

	fileLogger := slogger.NewLogger(f)
	for _, test := range tests {
		fileLogger.LogEvent(test.eventType, test.msg, test.eventData...)
		if !mapTester(test.have, readLogFileData(f, test.want)) {
			t.Fatalf("test %s failed. Have: %v Want: %v", test.name, test.have, test.want)
		}

		// Clean log file between each run
		test.clean(f)
	}
}

// replace test file with 0 bytes
func wipeFileData(f *os.File) {
	if err := os.Truncate(f.Name(), 0); err != nil {
		log.Println("Error removing file contents")
	}
}

// Map json data from file into a struct.
func readLogFileData(f *os.File, s map[string]any) map[string]any {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Printf("Error seeking: %v", err.Error())
	}
	b, err := os.ReadFile(f.Name())
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	json.Unmarshal(b, &s)

	return s
}

// Validates incoming map from test matches our desired map
func mapTester(m1, m2 map[string]any) bool {
	for k, v := range m1 {
		// Do not check time or level
		if m2Val, found := m2[k]; found {
			if v != m2Val {
				return false
			}
		} else {
			return false
		}
	}

	return true
}
