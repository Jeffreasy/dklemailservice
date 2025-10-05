package repository

import (
	"context"
	"dklautomationgo/models"
)

// ContactRepository definieert de interface voor contact formulier operaties
type ContactRepository interface {
	// Create slaat een nieuw contactformulier op
	Create(ctx context.Context, contact *models.ContactFormulier) error

	// GetByID haalt een contactformulier op basis van ID
	GetByID(ctx context.Context, id string) (*models.ContactFormulier, error)

	// List haalt een lijst van contactformulieren op
	List(ctx context.Context, limit, offset int) ([]*models.ContactFormulier, error)

	// Update werkt een bestaand contactformulier bij
	Update(ctx context.Context, contact *models.ContactFormulier) error

	// Delete verwijdert een contactformulier
	Delete(ctx context.Context, id string) error

	// FindByEmail zoekt contactformulieren op basis van email
	FindByEmail(ctx context.Context, email string) ([]*models.ContactFormulier, error)

	// FindByStatus zoekt contactformulieren op basis van status
	FindByStatus(ctx context.Context, status string) ([]*models.ContactFormulier, error)
}

// ContactAntwoordRepository definieert de interface voor contact antwoord operaties
type ContactAntwoordRepository interface {
	// Create slaat een nieuw contact antwoord op
	Create(ctx context.Context, antwoord *models.ContactAntwoord) error

	// GetByID haalt een contact antwoord op basis van ID
	GetByID(ctx context.Context, id string) (*models.ContactAntwoord, error)

	// ListByContactID haalt alle antwoorden voor een contact op
	ListByContactID(ctx context.Context, contactID string) ([]*models.ContactAntwoord, error)

	// Update werkt een bestaand contact antwoord bij
	Update(ctx context.Context, antwoord *models.ContactAntwoord) error

	// Delete verwijdert een contact antwoord
	Delete(ctx context.Context, id string) error
}

// AanmeldingRepository definieert de interface voor aanmelding operaties
type AanmeldingRepository interface {
	// Create slaat een nieuwe aanmelding op
	Create(ctx context.Context, aanmelding *models.Aanmelding) error

	// GetByID haalt een aanmelding op basis van ID
	GetByID(ctx context.Context, id string) (*models.Aanmelding, error)

	// List haalt een lijst van aanmeldingen op
	List(ctx context.Context, limit, offset int) ([]*models.Aanmelding, error)

	// Update werkt een bestaande aanmelding bij
	Update(ctx context.Context, aanmelding *models.Aanmelding) error

	// Delete verwijdert een aanmelding
	Delete(ctx context.Context, id string) error

	// FindByEmail zoekt aanmeldingen op basis van email
	FindByEmail(ctx context.Context, email string) ([]*models.Aanmelding, error)

	// FindByStatus zoekt aanmeldingen op basis van status
	FindByStatus(ctx context.Context, status string) ([]*models.Aanmelding, error)
}

// AanmeldingAntwoordRepository definieert de interface voor aanmelding antwoord operaties
type AanmeldingAntwoordRepository interface {
	// Create slaat een nieuw aanmelding antwoord op
	Create(ctx context.Context, antwoord *models.AanmeldingAntwoord) error

	// GetByID haalt een aanmelding antwoord op basis van ID
	GetByID(ctx context.Context, id string) (*models.AanmeldingAntwoord, error)

	// ListByAanmeldingID haalt alle antwoorden voor een aanmelding op
	ListByAanmeldingID(ctx context.Context, aanmeldingID string) ([]*models.AanmeldingAntwoord, error)

	// Update werkt een bestaand aanmelding antwoord bij
	Update(ctx context.Context, antwoord *models.AanmeldingAntwoord) error

	// Delete verwijdert een aanmelding antwoord
	Delete(ctx context.Context, id string) error
}

// EmailTemplateRepository definieert de interface voor email template operaties
type EmailTemplateRepository interface {
	// Create slaat een nieuwe email template op
	Create(ctx context.Context, template *models.EmailTemplate) error

	// GetByID haalt een email template op basis van ID
	GetByID(ctx context.Context, id string) (*models.EmailTemplate, error)

	// GetByNaam haalt een email template op basis van naam
	GetByNaam(ctx context.Context, naam string) (*models.EmailTemplate, error)

	// List haalt een lijst van email templates op
	List(ctx context.Context, limit, offset int) ([]*models.EmailTemplate, error)

	// Update werkt een bestaande email template bij
	Update(ctx context.Context, template *models.EmailTemplate) error

	// Delete verwijdert een email template
	Delete(ctx context.Context, id string) error

	// FindActive haalt alle actieve email templates op
	FindActive(ctx context.Context) ([]*models.EmailTemplate, error)
}

// VerzondEmailRepository definieert de interface voor verzonden email operaties
type VerzondEmailRepository interface {
	// Create slaat een nieuwe verzonden email op
	Create(ctx context.Context, email *models.VerzondEmail) error

	// GetByID haalt een verzonden email op basis van ID
	GetByID(ctx context.Context, id string) (*models.VerzondEmail, error)

	// List haalt een lijst van verzonden emails op
	List(ctx context.Context, limit, offset int) ([]*models.VerzondEmail, error)

	// Update werkt een bestaande verzonden email bij
	Update(ctx context.Context, email *models.VerzondEmail) error

	// FindByContactID haalt verzonden emails op basis van contact ID
	FindByContactID(ctx context.Context, contactID string) ([]*models.VerzondEmail, error)

	// FindByAanmeldingID haalt verzonden emails op basis van aanmelding ID
	FindByAanmeldingID(ctx context.Context, aanmeldingID string) ([]*models.VerzondEmail, error)

	// FindByOntvanger haalt verzonden emails op basis van ontvanger
	FindByOntvanger(ctx context.Context, ontvanger string) ([]*models.VerzondEmail, error)
}

// GebruikerRepository definieert de interface voor gebruiker operaties
type GebruikerRepository interface {
	// Create slaat een nieuwe gebruiker op
	Create(ctx context.Context, gebruiker *models.Gebruiker) error

	// GetByID haalt een gebruiker op basis van ID
	GetByID(ctx context.Context, id string) (*models.Gebruiker, error)

	// GetByEmail haalt een gebruiker op basis van email
	GetByEmail(ctx context.Context, email string) (*models.Gebruiker, error)

	// List haalt een lijst van gebruikers op
	List(ctx context.Context, limit, offset int) ([]*models.Gebruiker, error)

	// Update werkt een bestaande gebruiker bij
	Update(ctx context.Context, gebruiker *models.Gebruiker) error

	// Delete verwijdert een gebruiker
	Delete(ctx context.Context, id string) error

	// UpdateLastLogin werkt de laatste login tijd van een gebruiker bij
	UpdateLastLogin(ctx context.Context, id string) error
}

// MigratieRepository definieert de interface voor migratie operaties
type MigratieRepository interface {
	// Create slaat een nieuwe migratie op
	Create(ctx context.Context, migratie *models.Migratie) error

	// GetByVersie haalt een migratie op basis van versie
	GetByVersie(ctx context.Context, versie string) (*models.Migratie, error)

	// List haalt een lijst van migraties op
	List(ctx context.Context) ([]*models.Migratie, error)

	// GetLatest haalt de laatste migratie op
	GetLatest(ctx context.Context) (*models.Migratie, error)
}

// IncomingEmailRepository definieert de interface voor inkomende e-mail operaties
type IncomingEmailRepository interface {
	// Create slaat een nieuwe inkomende e-mail op
	Create(ctx context.Context, email *models.IncomingEmail) error

	// GetByID haalt een inkomende e-mail op basis van ID
	GetByID(ctx context.Context, id string) (*models.IncomingEmail, error)

	// List haalt een lijst van inkomende e-mails op
	List(ctx context.Context, limit, offset int) ([]*models.IncomingEmail, error)

	// Update werkt een bestaande inkomende e-mail bij
	Update(ctx context.Context, email *models.IncomingEmail) error

	// Delete verwijdert een inkomende e-mail
	Delete(ctx context.Context, id string) error

	// FindByUID zoekt een inkomende e-mail op basis van UID
	FindByUID(ctx context.Context, uid string) (*models.IncomingEmail, error)

	// FindUnprocessed haalt alle onverwerkte e-mails op
	FindUnprocessed(ctx context.Context) ([]*models.IncomingEmail, error)

	// FindByAccountType zoekt inkomende e-mails op basis van account type
	FindByAccountType(ctx context.Context, accountType string) ([]*models.IncomingEmail, error)

	// ListByAccountTypePaginated haalt een lijst van inkomende e-mails op basis van account type en paginatie
	ListByAccountTypePaginated(ctx context.Context, accountType string, limit, offset int) ([]*models.IncomingEmail, int64, error)
}

// NotificationRepository definieert de interface voor notificaties
type NotificationRepository interface {
	// Create slaat een nieuwe notificatie op
	Create(ctx context.Context, notification *models.Notification) error

	// GetByID haalt een notificatie op basis van ID
	GetByID(ctx context.Context, id string) (*models.Notification, error)

	// Update werkt een bestaande notificatie bij
	Update(ctx context.Context, notification *models.Notification) error

	// Delete verwijdert een notificatie
	Delete(ctx context.Context, id string) error

	// ListUnsent haalt alle niet verzonden notificaties op
	ListUnsent(ctx context.Context) ([]*models.Notification, error)

	// ListByType haalt alle notificaties op van een bepaald type
	ListByType(ctx context.Context, notificationType models.NotificationType) ([]*models.Notification, error)

	// ListByPriority haalt alle notificaties op met een bepaalde prioriteit
	ListByPriority(ctx context.Context, priority models.NotificationPriority) ([]*models.Notification, error)
}

// ChatChannelRepository defines the interface for chat channel operations
type ChatChannelRepository interface {
	Create(ctx context.Context, channel *models.ChatChannel) error
	GetByID(ctx context.Context, id string) (*models.ChatChannel, error)
	List(ctx context.Context, limit, offset int) ([]*models.ChatChannel, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.ChatChannel, error)
	Update(ctx context.Context, channel *models.ChatChannel) error
	Delete(ctx context.Context, id string) error
	ListPublicChannels(ctx context.Context) ([]*models.ChatChannel, error)
}

// ChatChannelParticipantRepository defines the interface for chat channel participant operations
type ChatChannelParticipantRepository interface {
	Create(ctx context.Context, participant *models.ChatChannelParticipant) error
	GetByID(ctx context.Context, id string) (*models.ChatChannelParticipant, error)
	List(ctx context.Context, limit, offset int) ([]*models.ChatChannelParticipant, error)
	Update(ctx context.Context, participant *models.ChatChannelParticipant) error
	Delete(ctx context.Context, id string) error
	ListByChannelID(ctx context.Context, channelID string) ([]*models.ChatChannelParticipant, error)
	GetByChannelAndUser(ctx context.Context, channelID, userID string) (*models.ChatChannelParticipant, error)
}

// ChatMessageRepository defines the interface for chat message operations
type ChatMessageRepository interface {
	Create(ctx context.Context, message *models.ChatMessage) error
	GetByID(ctx context.Context, id string) (*models.ChatMessage, error)
	List(ctx context.Context, limit, offset int) ([]*models.ChatMessage, error)
	Update(ctx context.Context, message *models.ChatMessage) error
	Delete(ctx context.Context, id string) error
	ListByChannelID(ctx context.Context, channelID string, limit, offset int) ([]*models.MessageWithUser, error)
}

// ChatMessageReactionRepository defines the interface for chat message reaction operations
type ChatMessageReactionRepository interface {
	Create(ctx context.Context, reaction *models.ChatMessageReaction) error
	GetByID(ctx context.Context, id string) (*models.ChatMessageReaction, error)
	List(ctx context.Context, limit, offset int) ([]*models.ChatMessageReaction, error)
	Delete(ctx context.Context, id string) error
	ListByMessageID(ctx context.Context, messageID string) ([]*models.ChatMessageReaction, error)
}

// ChatUserPresenceRepository defines the interface for chat user presence operations
type ChatUserPresenceRepository interface {
	Upsert(ctx context.Context, presence *models.ChatUserPresence) error
	GetByUserID(ctx context.Context, userID string) (*models.ChatUserPresence, error)
	Delete(ctx context.Context, userID string) error
	ListOnlineUsers(ctx context.Context) ([]*models.OnlineUser, error)
}
