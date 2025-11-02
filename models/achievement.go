package models

import (
	"time"
)

// ParticipantAchievement representeert een verdiende badge door een deelnemer
type ParticipantAchievement struct {
	ID            string    `json:"id" db:"id"`
	ParticipantID string    `json:"participant_id" db:"participant_id"`
	BadgeID       string    `json:"badge_id" db:"badge_id"`
	EarnedAt      time.Time `json:"earned_at" db:"earned_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// AchievementWithBadge bevat achievement informatie met badge details
type AchievementWithBadge struct {
	ID            string    `json:"id"`
	ParticipantID string    `json:"participant_id"`
	EarnedAt      time.Time `json:"earned_at"`
	Badge         Badge     `json:"badge"`
}

// AchievementRequest representeert een request voor het toekennen van een achievement
type AchievementRequest struct {
	ParticipantID string `json:"participant_id" validate:"required,uuid"`
	BadgeID       string `json:"badge_id" validate:"required,uuid"`
}

// ParticipantAchievementsSummary bevat een overzicht van achievements voor een deelnemer
type ParticipantAchievementsSummary struct {
	ParticipantID   string                 `json:"participant_id"`
	ParticipantName string                 `json:"participant_name"`
	TotalBadges     int                    `json:"total_badges"`
	TotalPoints     int                    `json:"total_points"`
	Achievements    []AchievementWithBadge `json:"achievements"`
	AvailableBadges []BadgeWithStats       `json:"available_badges"`
}
