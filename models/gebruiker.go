package models

import (
	"time"
)

// RBACRole is imported here to avoid circular imports
// The actual definition is in role_rbac.go

// Gebruiker representeert een gebruiker van het systeem
type Gebruiker struct {
	ID             string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Naam           string `json:"naam" gorm:"not null"`
	Email          string `json:"email" gorm:"not null;uniqueIndex"`
	WachtwoordHash string `json:"-" gorm:"not null"` // Niet zichtbaar in JSON

	IsActief             bool       `json:"is_actief" gorm:"default:true"`
	NewsletterSubscribed bool       `json:"newsletter_subscribed" gorm:"default:false;index"`
	LaatsteLogin         *time.Time `json:"laatste_login"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Roles []RBACRole `gorm:"many2many:user_roles;" json:"roles,omitempty"` // Many-to-many with RBAC roles
}

// TableName specificeert de tabelnaam voor GORM
func (Gebruiker) TableName() string {
	return "gebruikers"
}

// GebruikerLogin representeert de login gegevens voor een gebruiker
type GebruikerLogin struct {
	Email      string `json:"email" binding:"required"`
	Wachtwoord string `json:"wachtwoord" binding:"required"`
}

// AuthLoginResponse representeert de response van een login verzoek
type AuthLoginResponse struct {
	Success      bool             `json:"success"`
	Token        string           `json:"token"`
	RefreshToken string           `json:"refresh_token"`
	User         AuthUserResponse `json:"user"`
}

// AuthUserResponse representeert de user data in auth responses
type AuthUserResponse struct {
	ID          string                   `json:"id"`
	Email       string                   `json:"email"`
	Naam        string                   `json:"naam"`
	Permissions []map[string]string      `json:"permissions"`
	Roles       []map[string]interface{} `json:"roles"`
	IsActief    bool                     `json:"is_actief"`
}

// AuthRefreshRequest representeert een refresh token verzoek
type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthRefreshResponse representeert de response van een refresh verzoek
type AuthRefreshResponse struct {
	Success      bool   `json:"success"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthForgotPasswordRequest representeert een forgot password verzoek
type AuthForgotPasswordRequest struct {
	Email string `json:"email"`
}

// AuthResetPasswordWithTokenRequest representeert een reset password met token verzoek
type AuthResetPasswordWithTokenRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// AuthSendVerificationRequest representeert een send verification verzoek
type AuthSendVerificationRequest struct {
	Email string `json:"email"`
}

// AuthVerifyEmailRequest representeert een verify email verzoek
type AuthVerifyEmailRequest struct {
	Token string `json:"token"`
}

// AuthProfileResponse representeert de response van een profile verzoek
type AuthProfileResponse struct {
	ID           string                   `json:"id"`
	Naam         string                   `json:"naam"`
	Email        string                   `json:"email"`
	Permissions  []map[string]string      `json:"permissions"`
	Roles        []map[string]interface{} `json:"roles"`
	IsActief     bool                     `json:"is_actief"`
	LaatsteLogin *string                  `json:"laatste_login"`
	CreatedAt    string                   `json:"created_at"`
}

// AuthResetPasswordRequest representeert een reset password verzoek
type AuthResetPasswordRequest struct {
	HuidigWachtwoord string `json:"huidig_wachtwoord"`
	NieuwWachtwoord  string `json:"nieuw_wachtwoord"`
}

// AuthDeleteAccountRequest representeert een account deletion verzoek
type AuthDeleteAccountRequest struct {
	Password string `json:"password"`
	Reason   string `json:"reason,omitempty"`
}

// AuthSessionResponse representeert een session in de response
type AuthSessionResponse struct {
	ID           string `json:"id"`
	DeviceInfo   string `json:"device_info"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	LoginTime    string `json:"login_time"`
	LastActivity string `json:"last_activity"`
	IsCurrent    bool   `json:"is_current"`
	DisplayName  string `json:"display_name"`
	LocationInfo string `json:"location_info"`
}
