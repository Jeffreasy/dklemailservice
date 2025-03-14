package tests

import (
	"dklautomationgo/logger"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// TestMain is de hoofdfunctie voor alle tests
func TestMain(m *testing.M) {
	// Voorbereidingen voor tests
	setupTemplateDir()

	// Voer de tests uit met een timeout
	done := make(chan int, 1)
	go func() {
		result := m.Run()
		done <- result
	}()

	// Wacht op voltooiing met timeout
	select {
	case result := <-done:
		os.Exit(result)
	case <-time.After(60 * time.Second): // 60 seconden timeout voor alle tests
		// Tests zijn vastgelopen, dump stack en exit
		// Het kan geen kwaad om de stack te dumpen voor debug doeleinden
		// debug.PrintStack()
		os.Exit(1)
	}
}

// setupTemplateDir zorgt ervoor dat de templates directory correct is ingesteld voor tests
func setupTemplateDir() {
	// Bepaal de huidige werkdirectory
	wd, err := os.Getwd()
	if err != nil {
		logger.Error("Kon werkdirectory niet bepalen", "error", err)
		return
	}

	// Ga naar de parent directory (van tests/ naar project root)
	projectRoot := filepath.Dir(wd)

	// Controleer of de templates directory bestaat
	templatesDir := filepath.Join(projectRoot, "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		logger.Error("Templates directory niet gevonden", "path", templatesDir)
		return
	}

	// Kopieer de templates naar de tests directory
	testTemplatesDir := filepath.Join(wd, "templates")

	// Maak de templates directory in de tests directory als deze nog niet bestaat
	if _, err := os.Stat(testTemplatesDir); os.IsNotExist(err) {
		if err := os.Mkdir(testTemplatesDir, 0755); err != nil {
			logger.Error("Kon templates directory niet aanmaken", "error", err)
			return
		}
	}

	// Kopieer alle template bestanden
	templateFiles, err := os.ReadDir(templatesDir)
	if err != nil {
		logger.Error("Kon templates directory niet lezen", "error", err)
		return
	}

	for _, file := range templateFiles {
		if file.IsDir() {
			continue
		}

		// Lees het bronbestand
		srcPath := filepath.Join(templatesDir, file.Name())
		content, err := os.ReadFile(srcPath)
		if err != nil {
			logger.Error("Kon template bestand niet lezen", "file", file.Name(), "error", err)
			continue
		}

		// Schrijf naar het doelbestand
		dstPath := filepath.Join(testTemplatesDir, file.Name())
		if err := os.WriteFile(dstPath, content, 0644); err != nil {
			logger.Error("Kon template bestand niet schrijven", "file", file.Name(), "error", err)
			continue
		}

		logger.Info("Template bestand gekopieerd", "file", file.Name())
	}

	logger.Info("Templates directory ingesteld", "path", testTemplatesDir)
}

func init() {
	// Zet de SKIP_ELK_TEST op true voor normale test runs
	// Alleen expliciet uitvoeren van elk_integration_test.go zou deze test draaien
	if os.Getenv("SKIP_ELK_TEST") == "" {
		os.Setenv("SKIP_ELK_TEST", "true")
	}

	// Zet de logger op error-only tijdens tests om ruis te verminderen
	logger.Setup(logger.ErrorLevel)
}

// NewTestRegistry maakt een nieuwe Prometheus registry voor tests
func NewTestRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

// GetTemplatesDir zoekt naar de templates directory en geeft het pad terug
func GetTemplatesDir() (string, error) {
	// Probeer eerst de templates directory in de huidige directory
	if _, err := os.Stat("templates"); err == nil {
		absPath, err := filepath.Abs("templates")
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	// Probeer één niveau omhoog (voor als we in de tests directory zijn)
	if _, err := os.Stat("../templates"); err == nil {
		absPath, err := filepath.Abs("../templates")
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	// Probeer in de project root (als we in een subdirectory zijn)
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Zoek naar de templates directory in de huidige directory en alle parent directories
	for {
		templatesDir := filepath.Join(workDir, "templates")
		if _, err := os.Stat(templatesDir); err == nil {
			return templatesDir, nil
		}

		// Ga één niveau omhoog
		parent := filepath.Dir(workDir)
		if parent == workDir {
			// We hebben de root bereikt zonder de templates te vinden
			break
		}
		workDir = parent
	}

	return "", fmt.Errorf("templates directory niet gevonden")
}
