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

	// RBAC repositories
	RBACRole       RBACRoleRepository
	Permission     PermissionRepository
	RolePermission RolePermissionRepository
	UserRole       UserRoleRepository
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

		// RBAC repositories
		RBACRole:       NewRBACRoleRepository(db),
		Permission:     NewPermissionRepository(db),
		RolePermission: NewRolePermissionRepository(db),
		UserRole:       NewUserRoleRepository(db),
	}

	return repo
}
