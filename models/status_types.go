package models

// RegistrationStatusType definieert de statussen voor een registratie (V27)
type RegistrationStatusType struct {
	Status      string `json:"status" gorm:"primaryKey;column:status"` // Database column: status
	Description string `json:"description"`
}

// TableName specificeert de tabelnaam voor GORM
func (RegistrationStatusType) TableName() string {
	return "registration_status_types"
}

// EventStatusType representeert event status lookup table (V27)
type EventStatusType struct {
	Status      string `json:"status" gorm:"primaryKey;column:status"` // Database column: status
	Description string `json:"description" gorm:"type:text"`
}

func (EventStatusType) TableName() string {
	return "event_status_types"
}

// NotificationType representeert notification type lookup table (V27)
type NotificationType struct {
	Name        string `json:"name" gorm:"primaryKey;column:name"` // Database column: name
	Description string `json:"description" gorm:"type:text"`
}

func (NotificationType) TableName() string {
	return "notification_types"
}

// NotificationPriorityType representeert notification priority lookup table (V27)
type NotificationPriorityType struct {
	Name        string `json:"name" gorm:"primaryKey;column:name"` // Database column: name
	Description string `json:"description" gorm:"type:text"`
}

func (NotificationPriorityType) TableName() string {
	return "notification_priority_types"
}

// ContactStatusType representeert contact status lookup table (V27)
type ContactStatusType struct {
	Status      string `json:"status" gorm:"primaryKey;column:status"` // Database column: status
	Description string `json:"description" gorm:"type:text"`
}

func (ContactStatusType) TableName() string {
	return "contact_status_types"
}

// EmailStatusType representeert email status lookup table (V27)
type EmailStatusType struct {
	Status      string `json:"status" gorm:"primaryKey;column:status"` // Database column: status
	Description string `json:"description" gorm:"type:text"`
}

func (EmailStatusType) TableName() string {
	return "email_status_types"
}

// ChatChannelType representeert chat channel type lookup table (V27)
type ChatChannelType struct {
	Type        string `json:"type" gorm:"primaryKey;column:type"` // Database column: type (correct)
	Description string `json:"description" gorm:"type:text"`
}

func (ChatChannelType) TableName() string {
	return "chat_channel_types"
}
