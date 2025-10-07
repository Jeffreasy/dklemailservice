package models

import "time"

type NewsItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pub_date"`
	Category    string    `json:"category"`
}

type Newsletter struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Subject   string     `json:"subject"`
	Content   string     `json:"content"` // HTML
	SentAt    *time.Time `json:"sent_at"`
	BatchID   string     `json:"batch_id"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Newsletter) TableName() string { return "newsletters" }

// ProcessedNews bevat gefilterde en samengevatte nieuwsitems voor rendering
type ProcessedNews struct {
	Items   []NewsItem `json:"items"`
	Summary string     `json:"summary"`
}
