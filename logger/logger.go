package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// Beschikbare log niveaus
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// Voeg dit toe aan logger.go
type logFunc func(msg string, keysAndValues ...interface{})

type loggerFunctions struct {
	debugFunc logFunc
	infoFunc  logFunc
	warnFunc  logFunc
	errorFunc logFunc
	fatalFunc logFunc
}

var logFunctions *loggerFunctions

var logWriters []LogWriter

// Setup de default functies
func setupDefaultLogFunctions() {
	logFunctions = &loggerFunctions{
		debugFunc: defaultDebug,
		infoFunc:  defaultInfo,
		warnFunc:  defaultWarn,
		errorFunc: defaultError,
		fatalFunc: defaultFatal,
	}
}

// De oorspronkelijke implementaties met een andere naam
func defaultDebug(msg string, keysAndValues ...interface{}) {
	if log == nil {
		Setup(InfoLevel)
	}

	// Log naar zap logger (stdout/stderr)
	log.Debugw(msg, keysAndValues...)

	// Log naar alle extra writers
	if len(logWriters) > 0 {
		// Converteer naar map voor de writers
		entry := buildLogEntry(DebugLevel, msg, keysAndValues...)
		for _, writer := range logWriters {
			_ = writer.Write(entry)
		}
	}
}

func defaultInfo(msg string, keysAndValues ...interface{}) {
	if log == nil {
		Setup(InfoLevel)
	}

	// Log naar zap logger (stdout/stderr)
	log.Infow(msg, keysAndValues...)

	// Log naar alle extra writers
	if len(logWriters) > 0 {
		// Converteer naar map voor de writers
		entry := buildLogEntry(InfoLevel, msg, keysAndValues...)
		for _, writer := range logWriters {
			_ = writer.Write(entry)
		}
	}
}

func defaultWarn(msg string, keysAndValues ...interface{}) {
	if log == nil {
		Setup(InfoLevel)
	}

	// Log naar zap logger (stdout/stderr)
	log.Warnw(msg, keysAndValues...)

	// Log naar alle extra writers
	if len(logWriters) > 0 {
		// Converteer naar map voor de writers
		entry := buildLogEntry(WarnLevel, msg, keysAndValues...)
		for _, writer := range logWriters {
			_ = writer.Write(entry)
		}
	}
}

func defaultError(msg string, keysAndValues ...interface{}) {
	if log == nil {
		Setup(InfoLevel)
	}

	// Log naar zap logger (stdout/stderr)
	log.Errorw(msg, keysAndValues...)

	// Log naar alle extra writers
	if len(logWriters) > 0 {
		// Converteer naar map voor de writers
		entry := buildLogEntry(ErrorLevel, msg, keysAndValues...)
		for _, writer := range logWriters {
			_ = writer.Write(entry)
		}
	}
}

func defaultFatal(msg string, keysAndValues ...interface{}) {
	if log == nil {
		Setup(InfoLevel)
	}

	// Log naar zap logger (stdout/stderr)
	log.Fatalw(msg, keysAndValues...)

	// Log naar alle extra writers (alhoewel na Fatal niet veel meer gebeurt)
	if len(logWriters) > 0 {
		// Converteer naar map voor de writers
		entry := buildLogEntry("fatal", msg, keysAndValues...)
		for _, writer := range logWriters {
			_ = writer.Write(entry)
			_ = writer.Flush() // Direct flushen bij fatal logs
		}
	}
}

// Setup initialiseert de logger met het opgegeven niveau
func Setup(level string) {
	// Configuratie voor de encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "tijd",
		LevelKey:       "niveau",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "bericht",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Bepaal het log niveau
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Maak core aan die naar console schrijft
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// Maak de logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log = logger.Sugar()

	log.Infof("Logger geïnitialiseerd op niveau: %s", level)
}

// Vervang de oorspronkelijke publieke functies door wrappers
func Debug(msg string, keysAndValues ...interface{}) {
	if logFunctions == nil {
		setupDefaultLogFunctions()
	}
	logFunctions.debugFunc(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	if logFunctions == nil {
		setupDefaultLogFunctions()
	}
	logFunctions.infoFunc(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	if logFunctions == nil {
		setupDefaultLogFunctions()
	}
	logFunctions.warnFunc(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	if logFunctions == nil {
		setupDefaultLogFunctions()
	}
	logFunctions.errorFunc(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	if logFunctions == nil {
		setupDefaultLogFunctions()
	}
	logFunctions.fatalFunc(msg, keysAndValues...)
}

// Sync zorgt dat alle logs geschreven zijn
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// AddWriter voegt een extra log destination toe
func AddWriter(writer LogWriter) {
	logWriters = append(logWriters, writer)
}

// Helper om een log entry als map te bouwen
func buildLogEntry(level, msg string, keysAndValues ...interface{}) map[string]interface{} {
	entry := map[string]interface{}{
		"level":   level,
		"message": msg,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}

	// Parse extra velden
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if ok {
				// Maskeer gevoelige informatie
				value := keysAndValues[i+1]
				if isSensitiveKey(key) && value != nil {
					if strValue, ok := value.(string); ok && strValue != "" {
						entry[key] = maskSensitiveValue(strValue)
					} else {
						entry[key] = value
					}
				} else {
					entry[key] = value
				}
			}
		}
	}

	return entry
}

// isSensitiveKey controleert of een sleutel gevoelige informatie bevat
func isSensitiveKey(key string) bool {
	// Lijst met sleutelwoorden die gevoelige data kunnen bevatten
	sensitiveKeywords := []string{
		"password", "passwd", "secret", "token", "key", "auth",
		"credential", "cred", "connection_string", "dsn", "pwd",
	}

	lowerKey := strings.ToLower(key)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lowerKey, keyword) {
			return true
		}
	}

	return false
}

// maskSensitiveValue vervangt gevoelige waarden door asterisks
func maskSensitiveValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}

	// Toon alleen de eerste en laatste twee tekens
	return value[:2] + "******" + value[len(value)-2:]
}

// Toevoegen aan Setup of als aparte functie
func SetupELK(config ELKConfig) {
	elkWriter := NewELKWriter(config)
	AddWriter(elkWriter)
}

// Voeg een nieuwe functie toe voor het afhandelen van writers bij afsluiten
func CloseWriters() {
	Info("Applicatie wordt afgesloten, logs worden verzonden...")
	for _, writer := range logWriters {
		writer.Flush()
		writer.Close()
	}
}

// Shutdown sluit de logger en alle writers netjes af
func Shutdown() {
	// Stuur laatste logs
	Info("Applicatie wordt afgesloten, logs worden verzonden...")

	// Geef tijd om logs te verzenden
	time.Sleep(100 * time.Millisecond)

	// Sluit alle writers
	for _, writer := range logWriters {
		writer.Flush() // Eerst flushen om de buffer te legen
		writer.Close() // Dan sluiten
	}

	// Wacht op afsluiting
	time.Sleep(100 * time.Millisecond)
}
