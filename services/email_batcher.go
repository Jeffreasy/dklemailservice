package services

import (
	"dklautomationgo/logger"
	"sync"
	"time"
)

// EmailBatch representeert een groep emails die samen worden verwerkt
type EmailBatch struct {
	EmailType    string
	Recipients   []string
	Subject      string
	TemplateData map[string]interface{}
	TemplateName string
	FromAddress  string // Optional custom from address
	BatchID      string
	CreatedAt    time.Time
}

// EmailBatcher verzamelt emails in batches en verwerkt ze periodiek
type EmailBatcher struct {
	batchMap          map[string]*EmailBatch
	mutex             sync.Mutex
	batchSize         int
	batchWindow       time.Duration
	emailSvc          *EmailService
	ticker            *time.Ticker
	doneChan          chan bool
	prometheusMetrics *PrometheusMetrics
}

// NewEmailBatcher creëert een nieuwe email batcher
func NewEmailBatcher(emailSvc *EmailService, batchSize int, batchWindow time.Duration) *EmailBatcher {
	batcher := &EmailBatcher{
		batchMap:    make(map[string]*EmailBatch),
		batchSize:   batchSize,
		batchWindow: batchWindow,
		emailSvc:    emailSvc,
		doneChan:    make(chan bool),
	}

	// Start periodieke verwerking
	batcher.ticker = time.NewTicker(batchWindow)
	go batcher.processBatchesRoutine()

	return batcher
}

// AddToBatch voegt een email toe aan een batch
func (b *EmailBatcher) AddToBatch(batchKey, recipient, subject string,
	templateName string, templateData map[string]interface{}, fromAddress ...string) {

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Extract from address if provided
	var fromAddr string
	if len(fromAddress) > 0 {
		fromAddr = fromAddress[0]
	}

	// Maak een nieuwe batch als deze niet bestaat
	batch, exists := b.batchMap[batchKey]
	if !exists {
		batch = &EmailBatch{
			EmailType:    batchKey,
			Recipients:   make([]string, 0),
			Subject:      subject,
			TemplateData: templateData,
			TemplateName: templateName,
			FromAddress:  fromAddr,
			BatchID:      "batch-" + time.Now().Format("20060102-150405"),
			CreatedAt:    time.Now(),
		}
		b.batchMap[batchKey] = batch
	} else {
		// Update from address if batch already exists and fromAddress is provided
		if fromAddr != "" && batch.FromAddress == "" {
			batch.FromAddress = fromAddr
		}
	}

	// Voeg ontvanger toe
	batch.Recipients = append(batch.Recipients, recipient)

	// Verwerk meteen als we de batch size bereiken
	if len(batch.Recipients) >= b.batchSize {
		go b.processBatch(batchKey, batch)
		delete(b.batchMap, batchKey)
	}

	b.updateBatchCount()
}

// processBatchesRoutine verwerkt periodiek alle openstaande batches
func (b *EmailBatcher) processBatchesRoutine() {
	for {
		select {
		case <-b.ticker.C:
			b.processAllBatches()
		case <-b.doneChan:
			return
		}
	}
}

// processAllBatches verwerkt alle openstaande batches
func (b *EmailBatcher) processAllBatches() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now()
	for key, batch := range b.batchMap {
		// Verwerk batches die de maximum leeftijd hebben bereikt of niet leeg zijn
		if now.Sub(batch.CreatedAt) >= b.batchWindow && len(batch.Recipients) > 0 {
			go b.processBatch(key, batch)
			delete(b.batchMap, key)
		}
	}

	b.updateBatchCount()
}

// processBatch verwerkt één batch
func (b *EmailBatcher) processBatch(batchKey string, batch *EmailBatch) {
	if len(batch.Recipients) == 0 {
		return
	}

	logger.Info("Verwerken van email batch",
		"batch_key", batchKey,
		"batch_id", batch.BatchID,
		"recipients", len(batch.Recipients),
	)

	// Verwerk elke email sequentieel
	for _, recipient := range batch.Recipients {
		var err error
		if batch.FromAddress != "" {
			err = b.emailSvc.SendTemplateEmail(recipient, batch.Subject, batch.TemplateName, batch.TemplateData, batch.FromAddress)
		} else {
			err = b.emailSvc.SendTemplateEmail(recipient, batch.Subject, batch.TemplateName, batch.TemplateData)
		}

		if err != nil {
			logger.Error("Fout bij verzenden batch email",
				"error", err,
				"recipient", recipient,
				"batch_id", batch.BatchID,
			)
			continue // Ga door met de volgende email
		}
	}

	b.updateBatchCount()
}

// Shutdown stopt de batcher en verwerkt eventueel resterende batches
func (b *EmailBatcher) Shutdown() {
	b.ticker.Stop()
	b.doneChan <- true
	b.processAllBatches()
}

// Update count bijhouden in Prometheus
func (b *EmailBatcher) updateBatchCount() {
	if b.prometheusMetrics != nil {
		b.prometheusMetrics.UpdateActiveBatches(len(b.batchMap))
	}
}
