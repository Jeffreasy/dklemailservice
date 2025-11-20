package models

import (
	"time"
)

// Session representeert een gebruikerssessie voor multi-session management
// Sessions zijn gekoppeld aan access tokens voor gedetailleerde tracking
type Session struct {
	ID           string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OwnerID      string     `json:"owner_id" gorm:"column:owner_id;not null;type:uuid;index"` // Kan gebruiker of participant ID zijn
	AccessToken  string     `json:"access_token" gorm:"not null;uniqueIndex"`                 // Gekoppeld access token
	DeviceInfo   DeviceInfo `json:"device_info" gorm:"column:device_info;type:jsonb"`         // Device informatie
	IPAddress    string     `json:"ip_address" gorm:"not null"`
	UserAgent    string     `json:"user_agent" gorm:"type:text"`
	LoginTime    time.Time  `json:"login_time" gorm:"autoCreateTime"`
	LastActivity time.Time  `json:"last_activity" gorm:"autoUpdateTime"`
	ExpiresAt    time.Time  `json:"expires_at" gorm:"not null;index"`
	IsActive     bool       `json:"is_active" gorm:"default:true;index"`
	IsCurrent    bool       `json:"is_current" gorm:"default:false"` // Voor identificatie huidige sessie
}

// DeviceInfo bevat informatie over het apparaat en browser
type DeviceInfo struct {
	Browser        string `json:"browser" gorm:"size:100"`
	BrowserVersion string `json:"browser_version" gorm:"size:50"`
	OS             string `json:"os" gorm:"size:100"`
	OSVersion      string `json:"os_version" gorm:"size:50"`
	DeviceType     string `json:"device_type" gorm:"size:50"` // desktop, mobile, tablet
	Platform       string `json:"platform" gorm:"size:100"`   // Windows, macOS, Linux, iOS, Android
}

// TableName specificeert de tabelnaam voor GORM
func (Session) TableName() string {
	return "sessions"
}

// IsExpired controleert of de sessie is verlopen
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// UpdateActivity werkt de laatste activiteit bij
func (s *Session) UpdateActivity() {
	s.LastActivity = time.Now()
}

// GetDisplayName geeft een leesbare naam voor de sessie
func (s *Session) GetDisplayName() string {
	if s.DeviceInfo.Browser != "" && s.DeviceInfo.OS != "" {
		return s.DeviceInfo.Browser + " op " + s.DeviceInfo.OS
	}
	if s.DeviceInfo.DeviceType != "" {
		return s.DeviceInfo.DeviceType + " apparaat"
	}
	return "Onbekend apparaat"
}

// GetLocationInfo geeft locatie informatie gebaseerd op IP (placeholder voor toekomstige geo-IP integratie)
func (s *Session) GetLocationInfo() string {
	// Placeholder - zou kunnen worden uitgebreid met geo-IP service
	return "Locatie onbekend"
}
