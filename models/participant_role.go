package models

// ParticipantRole definieert de rollen (deelnemer, vrijwilliger, etc.)
type ParticipantRole struct {
	Name        string `json:"name" gorm:"primaryKey"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"` // Required by repository queries
}

// TableName specificeert de tabelnaam voor GORM
func (ParticipantRole) TableName() string {
	return "participant_roles"
}
