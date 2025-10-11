package repository

import (
	"gorm.io/gorm"
)

// Repository is een overkoepelende struct die alle specifieke repositories bevat
type Repository struct {
	Contact                ContactRepository
	ContactAntwoord        ContactAntwoordRepository
	Aanmelding             AanmeldingRepository
	AanmeldingAntwoord     AanmeldingAntwoordRepository
	Gebruiker              GebruikerRepository
	VerzondEmail           VerzondEmailRepository
	EmailTemplate          EmailTemplateRepository
	Migratie               MigratieRepository
	IncomingEmail          IncomingEmailRepository
	Notification           NotificationRepository
	ChatChannel            ChatChannelRepository
	ChatChannelParticipant ChatChannelParticipantRepository
	ChatMessage            ChatMessageRepository
	ChatMessageReaction    ChatMessageReactionRepository
	ChatUserPresence       ChatUserPresenceRepository
	Newsletter             NewsletterRepository
	UploadedImage          UploadedImageRepository
	Partner                PartnerRepository
	RadioRecording         RadioRecordingRepository
	Photo                  PhotoRepository
	Album                  AlbumRepository
	Video                  VideoRepository
	Sponsor                SponsorRepository
	ProgramSchedule        ProgramScheduleRepository
	SocialEmbed            SocialEmbedRepository
	SocialLink             SocialLinkRepository
	UnderConstruction      UnderConstructionRepository

	// RBAC repositories
	RBACRole       RBACRoleRepository
	Permission     PermissionRepository
	RolePermission RolePermissionRepository
	UserRole       UserRoleRepository
	RefreshToken   RefreshTokenRepository
}

// NewRepository maakt een nieuwe Repository met concrete implementaties
func NewRepository(db *gorm.DB) *Repository {
	baseRepo := NewPostgresRepository(db)

	repo := &Repository{
		Contact:                NewPostgresContactRepository(baseRepo),
		ContactAntwoord:        NewPostgresContactAntwoordRepository(baseRepo),
		Aanmelding:             NewPostgresAanmeldingRepository(baseRepo),
		AanmeldingAntwoord:     NewPostgresAanmeldingAntwoordRepository(baseRepo),
		Gebruiker:              NewPostgresGebruikerRepository(baseRepo),
		VerzondEmail:           NewPostgresVerzondEmailRepository(baseRepo),
		EmailTemplate:          NewPostgresEmailTemplateRepository(baseRepo),
		Migratie:               NewPostgresMigratieRepository(baseRepo),
		IncomingEmail:          NewPostgresIncomingEmailRepository(db),
		Notification:           NewPostgresNotificationRepository(baseRepo),
		ChatChannel:            NewPostgresChatChannelRepository(baseRepo),
		ChatChannelParticipant: NewPostgresChatChannelParticipantRepository(baseRepo),
		ChatMessage:            NewPostgresChatMessageRepository(baseRepo),
		ChatMessageReaction:    NewPostgresChatMessageReactionRepository(baseRepo),
		ChatUserPresence:       NewPostgresChatUserPresenceRepository(baseRepo),
		Newsletter:             NewPostgresNewsletterRepository(baseRepo),
		UploadedImage:          NewPostgresUploadedImageRepository(baseRepo),
		Partner:                NewPostgresPartnerRepository(db),
		RadioRecording:         NewPostgresRadioRecordingRepository(db),
		Photo:                  NewPostgresPhotoRepository(db),
		Album:                  NewPostgresAlbumRepository(db),
		Video:                  NewPostgresVideoRepository(db),
		Sponsor:                NewPostgresSponsorRepository(db),
		ProgramSchedule:        NewPostgresProgramScheduleRepository(db),
		SocialEmbed:            NewPostgresSocialEmbedRepository(db),
		SocialLink:             NewPostgresSocialLinkRepository(db),
		UnderConstruction:      NewPostgresUnderConstructionRepository(db),

		// RBAC repositories
		RBACRole:       NewRBACRoleRepository(db),
		Permission:     NewPermissionRepository(db),
		RolePermission: NewRolePermissionRepository(db),
		UserRole:       NewUserRoleRepository(db),
		RefreshToken:   NewPostgresRefreshTokenRepository(baseRepo),
	}

	return repo
}
