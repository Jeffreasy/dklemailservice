package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TelegramBotService beheert de Telegram bot functionaliteit
type TelegramBotService struct {
	botToken        string
	chatID          string
	client          *http.Client
	contactRepo     repository.ContactRepository
	aanmeldingRepo  repository.AanmeldingRepository
	commandHandlers map[string]CommandHandlerFunc
	updateOffset    int
	polling         bool
	pollingDone     chan struct{}
	mutex           sync.Mutex
}

// CommandHandlerFunc is een functie die een Telegram commando afhandelt
type CommandHandlerFunc func(update *TelegramUpdate) (string, error)

// TelegramUpdate bevat de gegevens van een Telegram update
type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			Type      string `json:"type"`
			Title     string `json:"title,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			Username  string `json:"username,omitempty"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// TelegramResponse bevat de respons van de Telegram API
type TelegramResponse struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"description,omitempty"`
}

// NewTelegramBotService maakt een nieuwe TelegramBotService
func NewTelegramBotService(
	contactRepo repository.ContactRepository,
	aanmeldingRepo repository.AanmeldingRepository,
) *TelegramBotService {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		logger.Error("Telegram bot configuratie ontbreekt")
		return nil
	}

	service := &TelegramBotService{
		botToken:        botToken,
		chatID:          chatID,
		client:          &http.Client{Timeout: 10 * time.Second},
		contactRepo:     contactRepo,
		aanmeldingRepo:  aanmeldingRepo,
		commandHandlers: make(map[string]CommandHandlerFunc),
		pollingDone:     make(chan struct{}),
	}

	// Registreer command handlers
	service.registerCommandHandlers()

	return service
}

// registerCommandHandlers registreert de beschikbare commando's
func (s *TelegramBotService) registerCommandHandlers() {
	s.commandHandlers = map[string]CommandHandlerFunc{
		"/start":         s.handleStartCommand,
		"/help":          s.handleHelpCommand,
		"/contact":       s.handleContactCommand,
		"/contactnew":    s.handleNewContactCommand,
		"/aanmelding":    s.handleAanmeldingCommand,
		"/aanmeldingnew": s.handleNewAanmeldingCommand,
		"/status":        s.handleStatusCommand,
	}
}

// StartPolling start het pollen voor nieuwe berichten
func (s *TelegramBotService) StartPolling() {
	s.mutex.Lock()
	if s.polling {
		s.mutex.Unlock()
		return
	}
	s.polling = true
	s.mutex.Unlock()

	logger.Info("Starting Telegram bot polling")

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Set up commands
		s.setMyCommands()

		for {
			select {
			case <-ticker.C:
				updates, err := s.getUpdates()
				if err != nil {
					logger.Error("Error getting updates", "error", err)
					continue
				}

				for _, update := range updates {
					s.processUpdate(&update)
				}
			case <-s.pollingDone:
				logger.Info("Stopping Telegram bot polling")
				return
			}
		}
	}()
}

// StopPolling stopt het pollen voor nieuwe berichten
func (s *TelegramBotService) StopPolling() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.polling {
		return
	}

	close(s.pollingDone)
	s.polling = false
}

// getUpdates haalt updates op van de Telegram API
func (s *TelegramBotService) getUpdates() ([]TelegramUpdate, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", s.botToken)

	// Parameters voor het request
	params := url.Values{}
	params.Add("offset", strconv.Itoa(s.updateOffset))
	params.Add("timeout", "1")

	// HTTP GET request
	resp, err := s.client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	// Lees de response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse de response
	var response TelegramResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("telegram API error: %s", response.Error)
	}

	// Parse de updates
	var updates []TelegramUpdate
	if err := json.Unmarshal(response.Result, &updates); err != nil {
		return nil, fmt.Errorf("failed to parse updates: %w", err)
	}

	// Update de offset om dubbele updates te voorkomen
	if len(updates) > 0 {
		s.updateOffset = updates[len(updates)-1].UpdateID + 1
	}

	return updates, nil
}

// processUpdate verwerkt een Telegram update
func (s *TelegramBotService) processUpdate(update *TelegramUpdate) {
	// Controleer of het bericht een commando is
	if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
		// Haal het commando uit het bericht
		command := strings.Split(update.Message.Text, " ")[0]

		// Behandel het commando
		if handler, ok := s.commandHandlers[command]; ok {
			response, err := handler(update)
			if err != nil {
				logger.Error("Error handling command",
					"command", command,
					"error", err)
				s.SendMessage("❌ Er is een fout opgetreden bij het uitvoeren van het commando. Probeer het later opnieuw.")
				return
			}

			// Stuur het antwoord
			s.SendMessage(response)
		} else {
			s.SendMessage(fmt.Sprintf("Onbekend commando: %s\nType /help voor een lijst met beschikbare commando's.", command))
		}
	}
}

// setMyCommands registreert de bot commando's bij Telegram
func (s *TelegramBotService) setMyCommands() error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setMyCommands", s.botToken)

	commands := []map[string]string{
		{"command": "start", "description": "Start de bot"},
		{"command": "help", "description": "Toon beschikbare commando's"},
		{"command": "contact", "description": "Toon recente contactformulieren"},
		{"command": "contactnew", "description": "Toon nieuwe contactformulieren"},
		{"command": "aanmelding", "description": "Toon recente aanmeldingen"},
		{"command": "aanmeldingnew", "description": "Toon onverwerkte aanmeldingen"},
		{"command": "status", "description": "Toon status van de service"},
	}

	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	// Parameters voor het request
	params := url.Values{}
	params.Add("commands", string(commandsJSON))

	// HTTP POST request
	resp, err := s.client.PostForm(apiURL, params)
	if err != nil {
		return fmt.Errorf("failed to set commands: %w", err)
	}
	defer resp.Body.Close()

	// Controleer de response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error, status code: %d", resp.StatusCode)
	}

	return nil
}

// SendMessage stuurt een bericht naar de Telegram chat
func (s *TelegramBotService) SendMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	// Parameters voor het bericht
	params := url.Values{}
	params.Add("chat_id", s.chatID)
	params.Add("text", message)
	params.Add("parse_mode", "HTML")

	// HTTP POST request
	resp, err := s.client.PostForm(apiURL, params)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	// Controleer de response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error, status code: %d", resp.StatusCode)
	}

	return nil
}

// handleStartCommand behandelt het /start commando
func (s *TelegramBotService) handleStartCommand(update *TelegramUpdate) (string, error) {
	return "👋 Welkom bij de DKL Email Service Bot!\n\n" +
		"Deze bot stelt je in staat om contactformulieren en aanmeldingen direct vanuit Telegram te bekijken.\n\n" +
		"Type /help voor een lijst met beschikbare commando's.", nil
}

// handleHelpCommand behandelt het /help commando
func (s *TelegramBotService) handleHelpCommand(update *TelegramUpdate) (string, error) {
	return "📋 <b>Beschikbare commando's:</b>\n\n" +
		"<b>/contact</b> - Toon recente contactformulieren\n" +
		"<b>/contactnew</b> - Toon nieuwe contactformulieren\n" +
		"<b>/aanmelding</b> - Toon recente aanmeldingen\n" +
		"<b>/aanmeldingnew</b> - Toon onverwerkte aanmeldingen\n" +
		"<b>/status</b> - Toon status van de service", nil
}

// handleContactCommand behandelt het /contact commando
func (s *TelegramBotService) handleContactCommand(update *TelegramUpdate) (string, error) {
	// Haal contactformulieren op
	ctx := context.Background()
	contacts, err := s.contactRepo.List(ctx, 5, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get contact forms: %w", err)
	}

	if len(contacts) == 0 {
		return "Geen contactformulieren gevonden.", nil
	}

	// Bouw het antwoord op
	var response strings.Builder
	response.WriteString("📬 <b>Recente contactformulieren:</b>\n\n")

	for i, contact := range contacts {
		status := contact.Status
		if status == "" {
			status = "nieuw"
		}

		response.WriteString(fmt.Sprintf("<b>%d. %s</b> (%s)\n", i+1, contact.Naam, status))
		response.WriteString(fmt.Sprintf("Email: %s\n", contact.Email))
		response.WriteString(fmt.Sprintf("Datum: %s\n", contact.CreatedAt.Format("02-01-2006 15:04")))

		// Verkort het bericht als het te lang is
		bericht := contact.Bericht
		if len(bericht) > 100 {
			bericht = bericht[:97] + "..."
		}
		response.WriteString(fmt.Sprintf("Bericht: %s\n\n", bericht))
	}

	return response.String(), nil
}

// handleNewContactCommand behandelt het /contactnew commando
func (s *TelegramBotService) handleNewContactCommand(update *TelegramUpdate) (string, error) {
	// Haal nieuwe contactformulieren op
	ctx := context.Background()
	contacts, err := s.contactRepo.FindByStatus(ctx, "nieuw")
	if err != nil {
		return "", fmt.Errorf("failed to get new contact forms: %w", err)
	}

	if len(contacts) == 0 {
		return "Er zijn geen nieuwe contactformulieren.", nil
	}

	// Bouw het antwoord op
	var response strings.Builder
	response.WriteString(fmt.Sprintf("🆕 <b>%d nieuwe contactformulieren:</b>\n\n", len(contacts)))

	for i, contact := range contacts {
		response.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, contact.Naam))
		response.WriteString(fmt.Sprintf("Email: %s\n", contact.Email))
		response.WriteString(fmt.Sprintf("Datum: %s\n", contact.CreatedAt.Format("02-01-2006 15:04")))

		// Verkort het bericht als het te lang is
		bericht := contact.Bericht
		if len(bericht) > 100 {
			bericht = bericht[:97] + "..."
		}
		response.WriteString(fmt.Sprintf("Bericht: %s\n\n", bericht))
	}

	return response.String(), nil
}

// handleAanmeldingCommand behandelt het /aanmelding commando
func (s *TelegramBotService) handleAanmeldingCommand(update *TelegramUpdate) (string, error) {
	// Haal aanmeldingen op
	ctx := context.Background()
	aanmeldingen, err := s.aanmeldingRepo.List(ctx, 5, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get registrations: %w", err)
	}

	if len(aanmeldingen) == 0 {
		return "Geen aanmeldingen gevonden.", nil
	}

	// Bouw het antwoord op
	var response strings.Builder
	response.WriteString("👥 <b>Recente aanmeldingen:</b>\n\n")

	for i, aanmelding := range aanmeldingen {
		response.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, aanmelding.Naam))
		response.WriteString(fmt.Sprintf("Email: %s\n", aanmelding.Email))
		response.WriteString(fmt.Sprintf("Datum: %s\n", aanmelding.CreatedAt.Format("02-01-2006 15:04")))
		response.WriteString(fmt.Sprintf("Rol: %s\n", aanmelding.Rol))
		response.WriteString(fmt.Sprintf("Afstand: %s\n", aanmelding.Afstand))

		if aanmelding.Telefoon != "" {
			response.WriteString(fmt.Sprintf("Telefoon: %s\n", aanmelding.Telefoon))
		}

		if aanmelding.Ondersteuning != "" && aanmelding.Ondersteuning != "Nee" {
			response.WriteString(fmt.Sprintf("Ondersteuning: %s\n", aanmelding.Ondersteuning))
		}

		if aanmelding.Bijzonderheden != "" {
			// Verkort bijzonderheden als het te lang is
			bijzonderheden := aanmelding.Bijzonderheden
			if len(bijzonderheden) > 100 {
				bijzonderheden = bijzonderheden[:97] + "..."
			}
			response.WriteString(fmt.Sprintf("Bijzonderheden: %s\n", bijzonderheden))
		}

		response.WriteString("\n")
	}

	return response.String(), nil
}

// handleNewAanmeldingCommand behandelt het /aanmeldingnew commando
func (s *TelegramBotService) handleNewAanmeldingCommand(update *TelegramUpdate) (string, error) {
	// Haal onverwerkte aanmeldingen op op basis van status
	ctx := context.Background()
	aanmeldingen, err := s.aanmeldingRepo.FindByStatus(ctx, "nieuw")
	if err != nil {
		return "", fmt.Errorf("failed to get unprocessed registrations: %w", err)
	}

	if len(aanmeldingen) == 0 {
		return "Er zijn geen onverwerkte aanmeldingen.", nil
	}

	// Bouw het antwoord op
	var response strings.Builder
	response.WriteString(fmt.Sprintf("🆕 <b>%d onverwerkte aanmeldingen:</b>\n\n", len(aanmeldingen)))

	for i, aanmelding := range aanmeldingen {
		response.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, aanmelding.Naam))
		response.WriteString(fmt.Sprintf("Email: %s\n", aanmelding.Email))
		response.WriteString(fmt.Sprintf("Datum: %s\n", aanmelding.CreatedAt.Format("02-01-2006 15:04")))
		response.WriteString(fmt.Sprintf("Rol: %s\n", aanmelding.Rol))
		response.WriteString(fmt.Sprintf("Afstand: %s\n", aanmelding.Afstand))

		if aanmelding.Telefoon != "" {
			response.WriteString(fmt.Sprintf("Telefoon: %s\n", aanmelding.Telefoon))
		}

		if aanmelding.Ondersteuning != "" && aanmelding.Ondersteuning != "Nee" {
			response.WriteString(fmt.Sprintf("Ondersteuning: %s\n", aanmelding.Ondersteuning))
		}

		if aanmelding.Bijzonderheden != "" {
			// Verkort bijzonderheden als het te lang is
			bijzonderheden := aanmelding.Bijzonderheden
			if len(bijzonderheden) > 100 {
				bijzonderheden = bijzonderheden[:97] + "..."
			}
			response.WriteString(fmt.Sprintf("Bijzonderheden: %s\n", bijzonderheden))
		}

		response.WriteString("\n")
	}

	return response.String(), nil
}

// handleStatusCommand behandelt het /status commando
func (s *TelegramBotService) handleStatusCommand(update *TelegramUpdate) (string, error) {
	ctx := context.Background()

	// Haal statistieken op met bestaande repository methoden
	contactList, err := s.contactRepo.List(ctx, 1000, 0)
	if err != nil {
		contactList = []*models.ContactFormulier{}
	}
	contactCount := int64(len(contactList))

	newContacts, err := s.contactRepo.FindByStatus(ctx, "nieuw")
	if err != nil {
		newContacts = []*models.ContactFormulier{}
	}
	newContactCount := int64(len(newContacts))

	aanmeldingList, err := s.aanmeldingRepo.List(ctx, 1000, 0)
	if err != nil {
		aanmeldingList = []*models.Aanmelding{}
	}
	aanmeldingCount := int64(len(aanmeldingList))

	// Gebruik status "nieuw" voor onverwerkte aanmeldingen
	unprocessedAanmeldingen, err := s.aanmeldingRepo.FindByStatus(ctx, "nieuw")
	if err != nil {
		unprocessedAanmeldingen = []*models.Aanmelding{}
	}
	unprocessedAanmeldingCount := int64(len(unprocessedAanmeldingen))

	return fmt.Sprintf("📊 <b>DKL Email Service Status</b>\n\n"+
		"<b>Contactformulieren:</b>\n"+
		"- Totaal: %d\n"+
		"- Nieuw: %d\n\n"+
		"<b>Aanmeldingen:</b>\n"+
		"- Totaal: %d\n"+
		"- Onverwerkt: %d\n\n"+
		"<b>Systeem:</b>\n"+
		"- Uptime: %s\n"+
		"- Bot actief: %t",
		contactCount, newContactCount,
		aanmeldingCount, unprocessedAanmeldingCount,
		formatDuration(time.Since(time.Now().Add(-24*time.Hour))), // Placeholder
		s.polling), nil
}

// formatDuration formatteert een duration naar een leesbare string
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
