package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// NotificationThrottleKey is een struct voor het bijhouden van throttling
type NotificationThrottleKey struct {
	Type    models.NotificationType
	Title   string
	Message string
}

// NotificationThrottleValue is een struct voor het bijhouden van throttling
type NotificationThrottleValue struct {
	LastSent time.Time
	Count    int
}

// NotificationClient is een interface voor het verzenden van notificaties
type NotificationClient interface {
	// SendMessage verstuurt een bericht
	SendMessage(title, message string) error
}

// TelegramClient implementeert NotificationClient met Telegram
type TelegramClient struct {
	BotToken string
	ChatID   string
	client   *http.Client
}

// NewTelegramClient maakt een nieuwe Telegram client
func NewTelegramClient(botToken, chatID string) *TelegramClient {
	return &TelegramClient{
		BotToken: botToken,
		ChatID:   chatID,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SendMessage stuurt een bericht naar Telegram
func (t *TelegramClient) SendMessage(title, message string) error {
	if t.BotToken == "" || t.ChatID == "" {
		return fmt.Errorf("telegram not configured (bot token: %v, chat id: %v)",
			t.BotToken != "", t.ChatID != "")
	}

	fullMessage := fmt.Sprintf("%s\n\n%s", title, message)

	// Telegram Bot API URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	// Parameters voor het bericht
	params := url.Values{}
	params.Add("chat_id", t.ChatID)
	params.Add("text", fullMessage)
	params.Add("parse_mode", "HTML")

	// HTTP POST request
	resp, err := t.client.PostForm(apiURL, params)
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

// NotificationServiceImpl implementeert NotificationService
type NotificationServiceImpl struct {
	notificationRepo repository.NotificationRepository
	client           NotificationClient
	throttleMap      map[NotificationThrottleKey]*NotificationThrottleValue
	throttleDuration time.Duration
	minPriority      models.NotificationPriority
	ticker           *time.Ticker
	running          bool
	mutex            sync.Mutex
	startupDone      bool
}

// NewNotificationService maakt een nieuwe notificatie service
func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	client NotificationClient,
	throttleDuration time.Duration,
	minPriority models.NotificationPriority,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		notificationRepo: notificationRepo,
		client:           client,
		throttleMap:      make(map[NotificationThrottleKey]*NotificationThrottleValue),
		throttleDuration: throttleDuration,
		minPriority:      minPriority,
		mutex:            sync.Mutex{},
		startupDone:      false,
	}
}

// CreateNotification maakt een nieuwe notificatie aan
func (s *NotificationServiceImpl) CreateNotification(
	ctx context.Context,
	notificationType models.NotificationType,
	priority models.NotificationPriority,
	title, message string,
) (*models.Notification, error) {
	notification := &models.Notification{
		Type:     notificationType,
		Priority: priority,
		Title:    title,
		Message:  message,
		Sent:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Als de prioriteit hoger is dan of gelijk aan de minimaal ingestelde prioriteit,
	// probeer de notificatie meteen te verzenden
	if isPriorityHighEnough(priority, s.minPriority) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// We negeren de fout hier, omdat de notificatie alsnog later verzonden kan worden
			_ = s.SendNotification(ctx, notification)
		}()
	}

	return notification, nil
}

// SendNotification verstuurt een notificatie
func (s *NotificationServiceImpl) SendNotification(ctx context.Context, notification *models.Notification) error {
	// Controleer of de notificatie al verzonden is
	if notification.Sent {
		return nil
	}

	// Controleer of de prioriteit hoog genoeg is
	if !isPriorityHighEnough(notification.Priority, s.minPriority) {
		logger.Info("Notificatie overgeslagen vanwege lage prioriteit",
			"id", notification.ID,
			"type", notification.Type,
			"priority", notification.Priority,
			"min_priority", s.minPriority)
		return nil
	}

	// Controleer throttling
	if !s.shouldSendNotification(notification) {
		logger.Info("Notificatie overgeslagen vanwege throttling",
			"id", notification.ID,
			"type", notification.Type,
			"title", notification.Title)
		return nil
	}

	// Voeg emoji toe op basis van prioriteit
	formattedTitle := formatTitleWithEmoji(notification.Priority, notification.Title)

	// Verstuur de notificatie
	err := s.client.SendMessage(formattedTitle, notification.Message)
	if err != nil {
		logger.Error("Fout bij het verzenden van notificatie",
			"id", notification.ID,
			"error", err)
		return err
	}

	// Update de notificatie in de database
	now := time.Now()
	notification.Sent = true
	notification.SentAt = &now
	if err := s.notificationRepo.Update(ctx, notification); err != nil {
		logger.Error("Fout bij het updaten van notificatie status",
			"id", notification.ID,
			"error", err)
		return err
	}

	logger.Info("Notificatie succesvol verzonden",
		"id", notification.ID,
		"type", notification.Type,
		"priority", notification.Priority)

	return nil
}

// GetNotification haalt een notificatie op basis van ID
func (s *NotificationServiceImpl) GetNotification(ctx context.Context, id string) (*models.Notification, error) {
	return s.notificationRepo.GetByID(ctx, id)
}

// ListUnsentNotifications haalt alle niet verzonden notificaties op
func (s *NotificationServiceImpl) ListUnsentNotifications(ctx context.Context) ([]*models.Notification, error) {
	return s.notificationRepo.ListUnsent(ctx)
}

// ProcessUnsentNotifications verwerkt alle niet verzonden notificaties
func (s *NotificationServiceImpl) ProcessUnsentNotifications(ctx context.Context) error {
	notifications, err := s.notificationRepo.ListUnsent(ctx)
	if err != nil {
		return fmt.Errorf("failed to list unsent notifications: %w", err)
	}

	if len(notifications) == 0 {
		return nil
	}

	logger.Info("Verwerken van onverzonden notificaties", "count", len(notifications))

	for _, notification := range notifications {
		if err := s.SendNotification(ctx, notification); err != nil {
			logger.Error("Fout bij het verzenden van notificatie",
				"id", notification.ID,
				"error", err)
			// Ga door met andere notificaties
			continue
		}
	}

	return nil
}

// Start begint het periodiek verzenden van notificaties
func (s *NotificationServiceImpl) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return
	}

	// Stuur een opstart notificatie
	if !s.startupDone {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			hostname, _ := os.Hostname()
			environment := getEnvWithDefault("ENVIRONMENT", "unknown")

			_, _ = s.CreateNotification(
				ctx,
				models.NotificationTypeSystem,
				models.NotificationPriorityLow,
				"Service Gestart",
				fmt.Sprintf("DKL Email Service is gestart op %s in omgeving %s.",
					hostname, environment),
			)

			s.startupDone = true
		}()
	}

	// Start een ticker voor periodiek checken van onverzonden notificaties
	s.ticker = time.NewTicker(1 * time.Minute)
	s.running = true

	go func() {
		for range s.ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := s.ProcessUnsentNotifications(ctx); err != nil {
				logger.Error("Fout bij het verwerken van onverzonden notificaties", "error", err)
			}
			cancel()
		}
	}()

	logger.Info("Notificatie service gestart",
		"throttle", s.throttleDuration.String(),
		"min_priority", s.minPriority)
}

// Stop stopt het periodiek verzenden van notificaties
func (s *NotificationServiceImpl) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return
	}

	// Stuur een shutdown notificatie
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		hostname, _ := os.Hostname()
		environment := getEnvWithDefault("ENVIRONMENT", "unknown")

		_, _ = s.CreateNotification(
			ctx,
			models.NotificationTypeSystem,
			models.NotificationPriorityLow,
			"Service Gestopt",
			fmt.Sprintf("DKL Email Service is gestopt op %s in omgeving %s.",
				hostname, environment),
		)
	}()

	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.running = false
	logger.Info("Notificatie service gestopt")
}

// IsRunning controleert of de service actief is
func (s *NotificationServiceImpl) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}

// shouldSendNotification controleert of een notificatie verzonden mag worden op basis van throttling
func (s *NotificationServiceImpl) shouldSendNotification(notification *models.Notification) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Kritieke notificaties altijd verzenden
	if notification.Priority == models.NotificationPriorityCritical {
		return true
	}

	key := NotificationThrottleKey{
		Type:    notification.Type,
		Title:   notification.Title,
		Message: notification.Message,
	}

	now := time.Now()
	value, exists := s.throttleMap[key]

	if !exists {
		// Nieuwe notificatie, nog niet eerder gezien
		s.throttleMap[key] = &NotificationThrottleValue{
			LastSent: now,
			Count:    1,
		}
		return true
	}

	// Als de laatste notificatie langer geleden is dan de throttle duur,
	// stuur dan opnieuw en reset de teller
	if now.Sub(value.LastSent) > s.throttleDuration {
		value.LastSent = now
		value.Count = 1
		return true
	}

	// Verhoog de teller
	value.Count++

	// Voor sommige notificaties willen we misschien toch een update sturen,
	// bijvoorbeeld als er erg veel zijn in korte tijd
	if value.Count >= 10 && notification.Priority == models.NotificationPriorityHigh {
		value.LastSent = now
		value.Count = 0
		return true
	}

	return false
}

// isPriorityHighEnough controleert of een prioriteit hoog genoeg is
func isPriorityHighEnough(priority, minPriority models.NotificationPriority) bool {
	priorityMap := map[models.NotificationPriority]int{
		models.NotificationPriorityLow:      1,
		models.NotificationPriorityMedium:   2,
		models.NotificationPriorityHigh:     3,
		models.NotificationPriorityCritical: 4,
	}

	return priorityMap[priority] >= priorityMap[minPriority]
}

// formatTitleWithEmoji voegt emoji toe aan een titel op basis van prioriteit
func formatTitleWithEmoji(priority models.NotificationPriority, title string) string {
	var emoji string
	switch priority {
	case models.NotificationPriorityLow:
		emoji = "ℹ️"
	case models.NotificationPriorityMedium:
		emoji = "⚠️"
	case models.NotificationPriorityHigh:
		emoji = "🔴"
	case models.NotificationPriorityCritical:
		emoji = "🚨"
	default:
		emoji = "📩"
	}

	return fmt.Sprintf("%s %s", emoji, title)
}
