package logger

import (
	"sync"
)

// LogEntry representeert een enkel logbericht
type LogEntry struct {
	Level   string
	Message string
	Fields  map[string]interface{}
}

// TestLogger implementeert een in-memory logger voor tests
type TestLogger struct {
	entries []LogEntry
	mu      sync.Mutex // Voor thread-safety
}

// NewTestLogger maakt een nieuwe TestLogger aan
func NewTestLogger() *TestLogger {
	return &TestLogger{
		entries: make([]LogEntry, 0),
	}
}

// GetEntries geeft alle geregistreerde logberichten terug
func (l *TestLogger) GetEntries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Maak een kopie van de entries om thread-safety te waarborgen
	result := make([]LogEntry, len(l.entries))
	copy(result, l.entries)
	return result
}

// AddEntry voegt een logbericht toe aan de lijst
func (l *TestLogger) AddEntry(level, message string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = append(l.entries, LogEntry{
		Level:   level,
		Message: message,
		Fields:  fields,
	})
}

// Reset verwijdert alle opgeslagen logberichten
func (l *TestLogger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = make([]LogEntry, 0)
}

// Implementatie van de log functionaliteit voor tests
func (l *TestLogger) logMessage(level, msg string, keysAndValues ...interface{}) {
	// Converteer key-value pairs naar een map
	fields := make(map[string]interface{})
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if ok {
				fields[key] = keysAndValues[i+1]
			}
		}
	}

	l.AddEntry(level, msg, fields)
}

// Testlogger Debug implementatie
func (l *TestLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logMessage(DebugLevel, msg, keysAndValues...)
}

// Testlogger Info implementatie
func (l *TestLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logMessage(InfoLevel, msg, keysAndValues...)
}

// Testlogger Warn implementatie
func (l *TestLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logMessage(WarnLevel, msg, keysAndValues...)
}

// Testlogger Error implementatie
func (l *TestLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logMessage(ErrorLevel, msg, keysAndValues...)
}

// Testlogger Fatal implementatie (niet fataal in tests)
func (l *TestLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logMessage("fatal", msg, keysAndValues...)
}

// Globale instance voor tests
var testLoggerInstance *TestLogger

// Initialiseer de testlogger voor gebruik in tests
func UseTestLogger() *TestLogger {
	// Maak een nieuwe testlogger als die nog niet bestaat
	if testLoggerInstance == nil {
		testLoggerInstance = NewTestLogger()
	} else {
		testLoggerInstance.Reset()
	}

	// Sla de originele logger op als nodig

	// Vervang de functies in het logger package
	logFunctions = &loggerFunctions{
		debugFunc: testLoggerInstance.Debug,
		infoFunc:  testLoggerInstance.Info,
		warnFunc:  testLoggerInstance.Warn,
		errorFunc: testLoggerInstance.Error,
		fatalFunc: testLoggerInstance.Fatal,
	}

	return testLoggerInstance
}

// RestoreLogger herstelt de originele logger
func RestoreDefaultLogger() {
	// Reset naar de standaard logger implementatie
	log = nil
	setupDefaultLogFunctions()
}
