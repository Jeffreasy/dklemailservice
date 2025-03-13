package tests

import (
	"dklautomationgo/logger"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// TestMain is de hoofdfunctie voor alle tests
func TestMain(m *testing.M) {
	// Voorbereidingen voor tests

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
