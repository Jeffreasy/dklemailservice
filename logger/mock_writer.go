package logger

// MockLogWriter implementeert de LogWriter interface voor tests
type MockLogWriter struct {
	Logs []map[string]interface{}
}

// NewMockLogWriter maakt een nieuwe mock writer
func NewMockLogWriter() *MockLogWriter {
	return &MockLogWriter{
		Logs: make([]map[string]interface{}, 0),
	}
}

// Write voegt een log toe aan de interne opslag
func (w *MockLogWriter) Write(entry map[string]interface{}) error {
	// Maak een kopie van de entry om wijzigingen te isoleren
	entryCopy := make(map[string]interface{})
	for k, v := range entry {
		entryCopy[k] = v
	}
	w.Logs = append(w.Logs, entryCopy)
	return nil
}

// Flush doet niets in de mock implementatie
func (w *MockLogWriter) Flush() error {
	return nil
}

// Close doet niets in de mock implementatie
func (w *MockLogWriter) Close() error {
	return nil
}

// GetLogs geeft alle opgeslagen logs terug
func (w *MockLogWriter) GetLogs() []map[string]interface{} {
	return w.Logs
}

// Reset verwijdert alle opgeslagen logs
func (w *MockLogWriter) Reset() {
	w.Logs = make([]map[string]interface{}, 0)
}
