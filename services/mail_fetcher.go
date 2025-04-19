package services

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/quotedprintable"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/microcosm-cc/bluemonday"
)

// MailAccount bevat de configuratie voor een mail account
type MailAccount struct {
	Username string
	Password string
	Host     string
	Port     int
	Type     string // "info" of "inschrijving"
}

// MailFetcher is verantwoordelijk voor het ophalen van e-mails uit inboxen
type MailFetcher struct {
	accounts  []*MailAccount
	metrics   *EmailMetrics
	lastFetch time.Time
	mu        sync.RWMutex
}

// NewMailFetcher maakt een nieuwe MailFetcher
func NewMailFetcher(metrics *EmailMetrics) *MailFetcher {
	return &MailFetcher{
		accounts:  make([]*MailAccount, 0),
		metrics:   metrics,
		lastFetch: time.Now().Add(-24 * time.Hour), // Start met ophalen vanaf 24 uur geleden
	}
}

// AddAccount voegt een mail account toe aan de fetcher
func (f *MailFetcher) AddAccount(username, password, host string, port int, accountType string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	account := &MailAccount{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Type:     accountType,
	}

	f.accounts = append(f.accounts, account)
	logger.Info("Mail account toegevoegd", "username", username, "host", host, "type", accountType)
}

// FetchMails haalt e-mails op van alle geconfigureerde accounts
func (f *MailFetcher) FetchMails() ([]*models.IncomingEmail, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var allMails []*models.IncomingEmail
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)

	// Maak een timestamp van nu om de lastFetch bij te werken
	fetchTime := time.Now()

	// Begin FetchMails
	currentLastFetch := f.lastFetch
	logger.Info("FetchMails start", "huidige_lastFetch", currentLastFetch)

	for _, account := range f.accounts {
		wg.Add(1)
		go func(acc *MailAccount) {
			defer wg.Done()

			mails, err := f.fetchFromAccount(acc, f.lastFetch)
			if err != nil {
				logger.Error("Fout bij ophalen e-mails", "error", err, "account", acc.Username)
				mu.Lock()
				errors = append(errors, fmt.Errorf("account %s: %w", acc.Username, err))
				mu.Unlock()
				return
			}

			if len(mails) > 0 {
				mu.Lock()
				allMails = append(allMails, mails...)
				mu.Unlock()
				logger.Info("E-mails opgehaald", "count", len(mails), "account", acc.Username)
			}
		}(account)
	}

	wg.Wait()

	// Update de lastFetch timestamp voor de volgende keer
	f.lastFetch = fetchTime

	// Vlak voor return allMails, nil of de error return
	logger.Info("FetchMails einde", "nieuwe_lastFetch", f.lastFetch)

	if len(errors) > 0 {
		errMessages := make([]string, len(errors))
		for i, err := range errors {
			errMessages[i] = err.Error()
		}
		return allMails, fmt.Errorf("errors tijdens ophalen e-mails: %s", strings.Join(errMessages, "; "))
	}

	return allMails, nil
}

// fetchFromAccount haalt e-mails op van één specifiek account
func (f *MailFetcher) fetchFromAccount(account *MailAccount, since time.Time) ([]*models.IncomingEmail, error) {
	// Begin fetchFromAccount
	logger.Info("Start fetchFromAccount", "account", account.Username, "since_parameter", since)
	// Verbind met de IMAP server
	imapAddr := fmt.Sprintf("%s:%d", account.Host, account.Port)
	c, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("kan niet verbinden met IMAP server: %w", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(account.Username, account.Password); err != nil {
		return nil, fmt.Errorf("login mislukt: %w", err)
	}

	// Selecteer INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("kan inbox niet selecteren: %w", err)
	}

	// Geen berichten in inbox
	if mbox.Messages == 0 {
		return []*models.IncomingEmail{}, nil
	}

	// Zoek berichten van na de laatste keer ophalen
	criteria := imap.NewSearchCriteria()
	criteria.Since = since

	// Voer zoekopdracht uit
	uids, err := c.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("zoeken mislukt: %w", err)
	}

	logger.Info("IMAP Search resultaat", "account", account.Username, "aantal_uids", len(uids))

	if len(uids) == 0 {
		return []*models.IncomingEmail{}, nil
	}

	// Maak een sequence set voor berichten
	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	// Items om op te halen
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid}

	// Haal berichten op
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	var emails []*models.IncomingEmail
	for msg := range messages {
		email, err := processMessage(msg, section, account.Type)
		if err != nil {
			logger.Warn("Fout bij verwerken bericht", "error", err, "uid", msg.Uid)
			continue
		}
		emails = append(emails, email)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("ophalen berichten mislukt: %w", err)
	}

	return emails, nil
}

// processMessage verwerkt een imap bericht naar een IncomingEmail model
func processMessage(msg *imap.Message, section imap.BodySectionName, accountType string) (*models.IncomingEmail, error) {
	var contentType string

	// Haal body op
	bodyReader := msg.GetBody(&section)
	if bodyReader == nil {
		return nil, fmt.Errorf("geen body gevonden")
	}

	// Parse het bericht
	m, err := mail.ReadMessage(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("kan bericht niet lezen: %w", err)
	}

	// Lees de headers
	contentType = m.Header.Get("Content-Type")
	contentEncoding := m.Header.Get("Content-Transfer-Encoding")
	from := m.Header.Get("From")
	subject := m.Header.Get("Subject")
	date := m.Header.Get("Date")
	messageId := m.Header.Get("Message-ID")

	// Maak de juiste reader aan op basis van de encoding
	var finalBodyReader io.Reader = m.Body
	switch strings.ToLower(contentEncoding) {
	case "quoted-printable":
		finalBodyReader = quotedprintable.NewReader(m.Body)
		logger.Debug("Using quoted-printable decoder", "message_id", messageId, "uid", msg.Uid)
	case "base64":
		finalBodyReader = base64.NewDecoder(base64.StdEncoding, m.Body)
		logger.Debug("Using base64 decoder", "message_id", messageId, "uid", msg.Uid)
	case "7bit", "8bit", "binary", "":
		// Geen decoding nodig, finalBodyReader blijft m.Body
		logger.Debug("No transfer decoding needed", "encoding", contentEncoding, "message_id", messageId, "uid", msg.Uid)
	default:
		logger.Warn("Onbekende Content-Transfer-Encoding", "encoding", contentEncoding, "message_id", messageId, "uid", msg.Uid)
		// Probeer toch direct te lezen
	}

	// Lees de (mogelijk gedecodeerde) body
	bodyBytes, err := ioutil.ReadAll(finalBodyReader)
	if err != nil {
		return nil, fmt.Errorf("kan body niet lezen na decoding: %w", err)
	}
	decodedBody := string(bodyBytes) // Dit is nu de gedecodeerde body

	// DEBUG: Log de gedecodeerde body
	logger.Debug("Decoded email body content", "message_id", messageId, "uid", msg.Uid, "decoded_body_preview", getFirstNChars(decodedBody, 200)) // Log first 200 chars

	// TODO: Overweeg charset conversie hier indien nodig, gebaseerd op contentType

	// Sanitize body if content type is text/html
	finalBody := decodedBody // Gebruik de gedecodeerde body als basis
	if strings.Contains(strings.ToLower(contentType), "text/html") {
		p := bluemonday.UGCPolicy()         // User Generated Content policy
		finalBody = p.Sanitize(decodedBody) // Sanitize de gedecodeerde body
		logger.Debug("HTML body gesanitized", "message_id", messageId, "uid", msg.Uid)
	}

	// Parse de datum
	var receivedAt time.Time
	if date != "" {
		parsedTime, err := mail.ParseDate(date)
		if err == nil {
			receivedAt = parsedTime
		} else {
			receivedAt = time.Now() // Fallback naar huidige tijd
		}
	} else {
		receivedAt = time.Now()
	}

	// Maak een IncomingEmail model met de CORRECTE veldnamen
	email := &models.IncomingEmail{
		MessageID:   messageId,
		From:        from,
		To:          accountType + "@dekoninklijkeloop.nl", // Houd de originele logica voor To aan
		Subject:     subject,
		Body:        finalBody, // Sla de gedecodeerde en mogelijk gesanitized body op
		ContentType: contentType,
		ReceivedAt:  receivedAt,
		AccountType: accountType,
		UID:         strconv.FormatUint(uint64(msg.Uid), 10),
		IsProcessed: false,
		ProcessedAt: nil, // Gebruik nil, niet false
	}

	return email, nil
}

// GetLastFetchTime retourneert de laatste fetch tijd
func (f *MailFetcher) GetLastFetchTime() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.lastFetch
}

// SetLastFetchTime stelt handmatig de laatste fetch tijd in
func (f *MailFetcher) SetLastFetchTime(t time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastFetch = t
}

// Helper function to get first N characters of a string
func getFirstNChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
