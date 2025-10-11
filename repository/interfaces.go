package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
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

	// GetNewsletterSubscribers haalt actieve subscribers op
	GetNewsletterSubscribers(ctx context.Context) ([]*models.Gebruiker, error)
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

// NewsletterRepository definieert de interface voor nieuwsbrief operaties
type NewsletterRepository interface {
	// Create slaat een nieuwe nieuwsbrief op
	Create(ctx context.Context, nl *models.Newsletter) error

	// GetByID haalt een nieuwsbrief op basis van ID
	GetByID(ctx context.Context, id string) (*models.Newsletter, error)

	// List haalt een lijst van nieuwsbrieven op
	List(ctx context.Context, limit, offset int) ([]*models.Newsletter, error)

	// Update werkt een bestaande nieuwsbrief bij
	Update(ctx context.Context, nl *models.Newsletter) error

	// Delete verwijdert een nieuwsbrief
	Delete(ctx context.Context, id string) error

	// UpdateBatchID werkt de batch ID bij na batching
	UpdateBatchID(ctx context.Context, id, batchID string) error

	// MarkSent markeert een nieuwsbrief als verzonden
	MarkSent(ctx context.Context, id string, sentAt time.Time) error
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

// UploadedImageRepository defines the interface for uploaded image operations
type UploadedImageRepository interface {
	// Create saves a new uploaded image record
	Create(ctx context.Context, image *models.UploadedImage) error

	// GetByID retrieves an uploaded image by ID
	GetByID(ctx context.Context, id string) (*models.UploadedImage, error)

	// GetByPublicID retrieves an uploaded image by Cloudinary public ID
	GetByPublicID(ctx context.Context, publicID string) (*models.UploadedImage, error)

	// GetByUserID retrieves uploaded images for a user with pagination
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.UploadedImage, error)

	// List retrieves a paginated list of all uploaded images
	List(ctx context.Context, limit, offset int) ([]*models.UploadedImage, error)

	// Update updates an existing uploaded image record
	Update(ctx context.Context, image *models.UploadedImage) error

	// Delete removes an uploaded image record
	Delete(ctx context.Context, id string) error

	// SoftDelete marks an image as deleted (for GDPR compliance)
	SoftDelete(ctx context.Context, id string) error

	// GetByFolder retrieves images by folder
	GetByFolder(ctx context.Context, folder string, limit, offset int) ([]*models.UploadedImage, error)
}

// PartnerRepository defines the interface for partner operations
type PartnerRepository interface {
	// Create saves a new partner
	Create(ctx context.Context, partner *models.Partner) error

	// GetByID retrieves a partner by ID
	GetByID(ctx context.Context, id string) (*models.Partner, error)

	// List retrieves a paginated list of partners
	List(ctx context.Context, limit, offset int) ([]*models.Partner, error)

	// ListVisible retrieves only visible partners ordered by order_number
	ListVisible(ctx context.Context) ([]*models.Partner, error)

	// Update updates an existing partner
	Update(ctx context.Context, partner *models.Partner) error

	// Delete removes a partner
	Delete(ctx context.Context, id string) error
}

// RadioRecordingRepository defines the interface for radio recording operations
type RadioRecordingRepository interface {
	// Create saves a new radio recording
	Create(ctx context.Context, recording *models.RadioRecording) error

	// GetByID retrieves a radio recording by ID
	GetByID(ctx context.Context, id string) (*models.RadioRecording, error)

	// List retrieves a paginated list of radio recordings
	List(ctx context.Context, limit, offset int) ([]*models.RadioRecording, error)

	// ListVisible retrieves only visible radio recordings ordered by order_number
	ListVisible(ctx context.Context) ([]*models.RadioRecording, error)

	// Update updates an existing radio recording
	Update(ctx context.Context, recording *models.RadioRecording) error

	// Delete removes a radio recording
	Delete(ctx context.Context, id string) error
}

// PhotoRepository defines the interface for photo operations
type PhotoRepository interface {
	// Create saves a new photo
	Create(ctx context.Context, photo *models.Photo) error

	// GetByID retrieves a photo by ID
	GetByID(ctx context.Context, id string) (*models.Photo, error)

	// List retrieves a paginated list of photos
	List(ctx context.Context, limit, offset int) ([]*models.Photo, error)

	// ListVisible retrieves only visible photos
	ListVisible(ctx context.Context) ([]*models.Photo, error)

	// ListVisibleFiltered retrieves visible photos with filtering
	ListVisibleFiltered(ctx context.Context, filters map[string]interface{}) ([]*models.Photo, error)

	// ListByAlbumID retrieves photos for a specific album
	ListByAlbumID(ctx context.Context, albumID string) ([]*models.Photo, error)

	// Update updates an existing photo
	Update(ctx context.Context, photo *models.Photo) error

	// Delete removes a photo
	Delete(ctx context.Context, id string) error
}

// AlbumRepository defines the interface for album operations
type AlbumRepository interface {
	// Create saves a new album
	Create(ctx context.Context, album *models.Album) error

	// GetByID retrieves an album by ID
	GetByID(ctx context.Context, id string) (*models.Album, error)

	// List retrieves a paginated list of albums
	List(ctx context.Context, limit, offset int) ([]*models.Album, error)

	// ListVisible retrieves only visible albums ordered by order_number
	ListVisible(ctx context.Context) ([]*models.Album, error)

	// ListVisibleWithCovers retrieves visible albums with cover photo information
	ListVisibleWithCovers(ctx context.Context) ([]*models.AlbumWithCover, error)

	// Update updates an existing album
	Update(ctx context.Context, album *models.Album) error

	// Delete removes an album
	Delete(ctx context.Context, id string) error
}

// VideoRepository defines the interface for video operations
type VideoRepository interface {
	// Create saves a new video
	Create(ctx context.Context, video *models.Video) error

	// GetByID retrieves a video by ID
	GetByID(ctx context.Context, id string) (*models.Video, error)

	// List retrieves a paginated list of videos
	List(ctx context.Context, limit, offset int) ([]*models.Video, error)

	// ListVisible retrieves only visible videos ordered by order_number
	ListVisible(ctx context.Context) ([]*models.Video, error)

	// Update updates an existing video
	Update(ctx context.Context, video *models.Video) error

	// Delete removes a video
	Delete(ctx context.Context, id string) error
}

// SponsorRepository defines the interface for sponsor operations
type SponsorRepository interface {
	// Create saves a new sponsor
	Create(ctx context.Context, sponsor *models.Sponsor) error

	// GetByID retrieves a sponsor by ID
	GetByID(ctx context.Context, id string) (*models.Sponsor, error)

	// List retrieves a paginated list of sponsors
	List(ctx context.Context, limit, offset int) ([]*models.Sponsor, error)

	// ListVisible retrieves only visible sponsors ordered by order_number
	ListVisible(ctx context.Context) ([]*models.Sponsor, error)

	// Update updates an existing sponsor
	Update(ctx context.Context, sponsor *models.Sponsor) error

	// Delete removes a sponsor
	Delete(ctx context.Context, id string) error
}

// ProgramScheduleRepository defines the interface for program schedule operations
type ProgramScheduleRepository interface {
	// Create saves a new program schedule
	Create(ctx context.Context, schedule *models.ProgramSchedule) error

	// GetByID retrieves a program schedule by ID
	GetByID(ctx context.Context, id string) (*models.ProgramSchedule, error)

	// List retrieves a paginated list of program schedules
	List(ctx context.Context, limit, offset int) ([]*models.ProgramSchedule, error)

	// ListVisible retrieves only visible program schedules ordered by order_number
	ListVisible(ctx context.Context) ([]*models.ProgramSchedule, error)

	// Update updates an existing program schedule
	Update(ctx context.Context, schedule *models.ProgramSchedule) error

	// Delete removes a program schedule
	Delete(ctx context.Context, id string) error
}

// SocialEmbedRepository defines the interface for social embed operations
type SocialEmbedRepository interface {
	// Create saves a new social embed
	Create(ctx context.Context, embed *models.SocialEmbed) error

	// GetByID retrieves a social embed by ID
	GetByID(ctx context.Context, id string) (*models.SocialEmbed, error)

	// List retrieves a paginated list of social embeds
	List(ctx context.Context, limit, offset int) ([]*models.SocialEmbed, error)

	// ListVisible retrieves only visible social embeds ordered by order_number
	ListVisible(ctx context.Context) ([]*models.SocialEmbed, error)

	// Update updates an existing social embed
	Update(ctx context.Context, embed *models.SocialEmbed) error

	// Delete removes a social embed
	Delete(ctx context.Context, id string) error
}

// SocialLinkRepository defines the interface for social link operations
type SocialLinkRepository interface {
	// Create saves a new social link
	Create(ctx context.Context, link *models.SocialLink) error

	// GetByID retrieves a social link by ID
	GetByID(ctx context.Context, id string) (*models.SocialLink, error)

	// List retrieves a paginated list of social links
	List(ctx context.Context, limit, offset int) ([]*models.SocialLink, error)

	// ListVisible retrieves only visible social links ordered by order_number
	ListVisible(ctx context.Context) ([]*models.SocialLink, error)

	// Update updates an existing social link
	Update(ctx context.Context, link *models.SocialLink) error

	// Delete removes a social link
	Delete(ctx context.Context, id string) error
}

// UnderConstructionRepository defines the interface for under construction operations
type UnderConstructionRepository interface {
	// Create saves a new under construction record
	Create(ctx context.Context, uc *models.UnderConstruction) error

	// GetByID retrieves an under construction record by ID
	GetByID(ctx context.Context, id int) (*models.UnderConstruction, error)

	// List retrieves a paginated list of under construction records
	List(ctx context.Context, limit, offset int) ([]*models.UnderConstruction, error)

	// GetActive retrieves the active under construction record
	GetActive(ctx context.Context) (*models.UnderConstruction, error)

	// Update updates an existing under construction record
	Update(ctx context.Context, uc *models.UnderConstruction) error

	// Delete removes an under construction record
	Delete(ctx context.Context, id int) error
}

// AlbumPhotoRepository defines the interface for album photo operations
type AlbumPhotoRepository interface {
	// Create adds a photo to an album
	Create(ctx context.Context, albumPhoto *models.AlbumPhoto) error

	// Delete removes a photo from an album
	Delete(ctx context.Context, albumID, photoID string) error

	// GetByAlbumAndPhoto retrieves a specific album-photo relationship
	GetByAlbumAndPhoto(ctx context.Context, albumID, photoID string) (*models.AlbumPhoto, error)

	// ListByAlbum retrieves all photos for an album ordered by order_number
	ListByAlbum(ctx context.Context, albumID string) ([]*models.AlbumPhoto, error)

	// UpdateOrder updates the order number of a photo in an album
	UpdateOrder(ctx context.Context, albumID, photoID string, orderNumber int) error

	// DeleteByAlbum removes all photos from an album
	DeleteByAlbum(ctx context.Context, albumID string) error

	// DeleteByPhoto removes a photo from all albums
	DeleteByPhoto(ctx context.Context, photoID string) error
}
