package models

import (
	"time"
)

// LeaderboardEntry representeert één entry in de leaderboard
type LeaderboardEntry struct {
	ID                string    `json:"id" db:"id"`
	Naam              string    `json:"naam" db:"naam"`
	Route             string    `json:"route" db:"route"`
	Steps             int       `json:"steps" db:"steps"`
	AchievementPoints int       `json:"achievement_points" db:"achievement_points"`
	TotalScore        int       `json:"total_score" db:"total_score"`
	Rank              int       `json:"rank" db:"rank"`
	BadgeCount        int       `json:"badge_count" db:"badge_count"`
	JoinedAt          time.Time `json:"joined_at" db:"joined_at"`
}

// LeaderboardResponse bevat leaderboard data met metadata
type LeaderboardResponse struct {
	Entries      []LeaderboardEntry `json:"entries"`
	TotalEntries int                `json:"total_entries"`
	CurrentPage  int                `json:"current_page,omitempty"`
	TotalPages   int                `json:"total_pages,omitempty"`
	Limit        int                `json:"limit,omitempty"`
}

// LeaderboardFilters bevat filters voor leaderboard queries
type LeaderboardFilters struct {
	Route    *string `json:"route,omitempty"`
	Year     *int    `json:"year,omitempty"`
	MinSteps *int    `json:"min_steps,omitempty"`
	TopN     *int    `json:"top_n,omitempty"`
	Page     int     `json:"page,omitempty"`
	Limit    int     `json:"limit,omitempty"`
}

// ParticipantRankInfo bevat rank informatie voor een specifieke deelnemer
type ParticipantRankInfo struct {
	ParticipantID     string            `json:"participant_id"`
	Naam              string            `json:"naam"`
	Rank              int               `json:"rank"`
	TotalScore        int               `json:"total_score"`
	Steps             int               `json:"steps"`
	AchievementPoints int               `json:"achievement_points"`
	BadgeCount        int               `json:"badge_count"`
	AboveMe           *LeaderboardEntry `json:"above_me,omitempty"`
	BelowMe           *LeaderboardEntry `json:"below_me,omitempty"`
}
