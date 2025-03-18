package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/repository"
	"os"
	"strconv"
	"sync"
	"time"
)

// EmailAutoFetcher is een service die automatisch emails ophaalt op regelmatige intervallen
type EmailAutoFetcher struct {
	mailFetcher     *MailFetcher
	emailRepository repository.IncomingEmailRepository
	interval        time.Duration
	running         bool
	stopChan        chan struct{}
	mutex           sync.Mutex
	lastRunTime     time.Time
}

// NewEmailAutoFetcher maakt een nieuwe EmailAutoFetcher service met de opgegeven dependencies en configureert het interval
func NewEmailAutoFetcher(mailFetcher *MailFetcher, emailRepo repository.IncomingEmailRepository) *EmailAutoFetcher {
	// Standaard interval is 15 minuten
	defaultInterval := 15 * time.Minute

	// Lees interval uit omgevingsvariabele (in minuten)
	intervalStr := os.Getenv("EMAIL_FETCH_INTERVAL")
	if intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			defaultInterval = time.Duration(interval) * time.Minute
			logger.Info("Email auto fetcher interval ingesteld op custom waarde", "minutes", interval)
		}
	} else {
		logger.Info("Email auto fetcher interval ingesteld op standaard waarde", "minutes", 15)
	}

	return &EmailAutoFetcher{
		mailFetcher:     mailFetcher,
		emailRepository: emailRepo,
		interval:        defaultInterval,
		running:         false,
		stopChan:        make(chan struct{}),
		lastRunTime:     time.Time{}, // Zero time
	}
}

// Start begint het periodiek ophalen van emails
func (f *EmailAutoFetcher) Start() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.running {
		logger.Info("Email auto fetcher is al actief")
		return
	}

	f.running = true
	f.stopChan = make(chan struct{})

	logger.Info("Email auto fetcher gestart", "interval", f.interval)

	go f.fetchLoop()
}

// Stop stopt het periodiek ophalen van emails
func (f *EmailAutoFetcher) Stop() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.running {
		logger.Info("Email auto fetcher is niet actief")
		return
	}

	logger.Info("Email auto fetcher stoppen...")
	close(f.stopChan)
	f.running = false
}

// IsRunning controleert of de auto fetcher actief is
func (f *EmailAutoFetcher) IsRunning() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.running
}

// GetLastRunTime geeft de laatste keer terug dat emails werden opgehaald
func (f *EmailAutoFetcher) GetLastRunTime() time.Time {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.lastRunTime
}

// fetchLoop is een interne functie die het ophalen van emails op regelmatige intervallen verzorgt
func (f *EmailAutoFetcher) fetchLoop() {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	// Direct eerste keer uitvoeren
	f.fetchOnce()

	for {
		select {
		case <-ticker.C:
			f.fetchOnce()
		case <-f.stopChan:
			logger.Info("Email auto fetcher gestopt")
			return
		}
	}
}

// fetchOnce haalt eenmalig emails op en slaat ze op
func (f *EmailAutoFetcher) fetchOnce() {
	f.mutex.Lock()
	f.lastRunTime = time.Now()
	f.mutex.Unlock()

	logger.Info("Automatisch emails ophalen...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Haal emails op van alle accounts
	emails, err := f.mailFetcher.FetchMails()
	if err != nil {
		logger.Error("Fout bij automatisch ophalen van emails", "error", err)
		return
	}

	if len(emails) == 0 {
		logger.Info("Geen nieuwe emails gevonden")
		return
	}

	logger.Info("Nieuwe emails gevonden", "count", len(emails))

	// Sla nieuwe emails op in de database
	savedCount := 0
	for _, email := range emails {
		// Controleer eerst of de email al bestaat om dubbele emails te voorkomen
		existing, err := f.emailRepository.FindByUID(ctx, email.UID)
		if err == nil && existing != nil {
			logger.Debug("Email overgeslagen (bestaat al)", "uid", email.UID)
			continue
		}

		// Sla de nieuwe email op
		if err := f.emailRepository.Create(ctx, email); err != nil {
			logger.Error("Fout bij opslaan van nieuwe email", "error", err, "uid", email.UID)
			continue
		}

		savedCount++
	}

	logger.Info("Nieuwe emails opgeslagen", "count", savedCount)
}
