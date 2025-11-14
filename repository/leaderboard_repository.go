package repository

import (
	"context"
	"dklautomationgo/models"
	"fmt"

	"gorm.io/gorm"
)

// LeaderboardRepository interface voor leaderboard operaties
type LeaderboardRepository interface {
	GetLeaderboard(ctx context.Context, filters models.LeaderboardFilters) (*models.LeaderboardResponse, error)
	GetParticipantRank(ctx context.Context, participantID string) (*models.ParticipantRankInfo, error)
}

// PostgresLeaderboardRepository implementatie voor PostgreSQL
type PostgresLeaderboardRepository struct {
	*PostgresRepository
}

// NewLeaderboardRepository maakt een nieuwe leaderboard repository
func NewLeaderboardRepository(db *gorm.DB) LeaderboardRepository {
	return &PostgresLeaderboardRepository{
		PostgresRepository: NewPostgresRepository(db),
	}
}

// GetLeaderboard haalt de leaderboard op met filters
func (r *PostgresLeaderboardRepository) GetLeaderboard(ctx context.Context, filters models.LeaderboardFilters) (*models.LeaderboardResponse, error) {
	var entries []models.LeaderboardEntry

	// Build query met filters - FIXED: Join with event_registrations for steps and distance_route
	query := `
		SELECT
			p.id,
			p.naam,
			COALESCE(er.distance_route, 'Unknown') as route,
			COALESCE(er.steps, 0) as steps,
			COALESCE(SUM(b.points), 0) as achievement_points,
			COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0) as total_score,
			RANK() OVER (ORDER BY (COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0)) DESC) as rank,
			COUNT(pa.id) as badge_count,
			COALESCE(er.registered_at, p.created_at) as joined_at
		FROM participants p
		LEFT JOIN event_registrations er ON p.id = er.participant_id
			AND er.event_id = (SELECT get_active_event())
		LEFT JOIN participant_achievements pa ON p.id = pa.participant_id
		LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
		WHERE 1=1
	`

	args := []interface{}{}

	// Filter op route
	if filters.Route != nil {
		query += " AND er.distance_route = ?"
		args = append(args, *filters.Route)
	}

	// Filter op year (als je created_at gebruikt voor jaar filtering)
	if filters.Year != nil {
		query += " AND EXTRACT(YEAR FROM COALESCE(er.registered_at, p.created_at)) = ?"
		args = append(args, *filters.Year)
	}

	query += `
		GROUP BY p.id, p.naam, er.distance_route, er.steps, er.registered_at, p.created_at
	`

	// Filter op minimum steps
	if filters.MinSteps != nil {
		query += " HAVING COALESCE(er.steps, 0) >= ?"
		args = append(args, *filters.MinSteps)
	}

	query += " ORDER BY total_score DESC, steps DESC"

	// Limit voor top N
	if filters.TopN != nil {
		query += " LIMIT ?"
		args = append(args, *filters.TopN)
	} else if filters.Limit > 0 {
		offset := 0
		if filters.Page > 0 {
			offset = (filters.Page - 1) * filters.Limit
		}
		query += " LIMIT ? OFFSET ?"
		args = append(args, filters.Limit, offset)
	}

	err := r.DB().WithContext(ctx).Raw(query, args...).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	// Tel totaal entries (zonder limit)
	var totalCount int64
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM participants p
		LEFT JOIN event_registrations er ON p.id = er.participant_id
			AND er.event_id = (SELECT get_active_event())
		WHERE 1=1
	`
	countArgs := []interface{}{}

	if filters.Route != nil {
		countQuery += " AND er.distance_route = ?"
		countArgs = append(countArgs, *filters.Route)
	}

	if filters.Year != nil {
		countQuery += " AND EXTRACT(YEAR FROM COALESCE(er.registered_at, p.created_at)) = ?"
		countArgs = append(countArgs, *filters.Year)
	}

	err = r.DB().WithContext(ctx).Raw(countQuery, countArgs...).Scan(&totalCount).Error
	if err != nil {
		return nil, err
	}

	response := &models.LeaderboardResponse{
		Entries:      entries,
		TotalEntries: int(totalCount),
	}

	// Bereken pagination info
	if filters.Limit > 0 {
		response.CurrentPage = filters.Page
		if filters.Page == 0 {
			response.CurrentPage = 1
		}
		response.Limit = filters.Limit
		response.TotalPages = (int(totalCount) + filters.Limit - 1) / filters.Limit
	}

	return response, nil
}

// GetParticipantRank haalt rank informatie op voor een specifieke deelnemer
func (r *PostgresLeaderboardRepository) GetParticipantRank(ctx context.Context, participantID string) (*models.ParticipantRankInfo, error) {
	var info models.ParticipantRankInfo

	// Haal rank en scores op voor participant - FIXED: Join with event_registrations
	query := `
		WITH ranked_participants AS (
			SELECT
				p.id,
				p.naam,
				COALESCE(er.steps, 0) as steps,
				COALESCE(SUM(b.points), 0) as achievement_points,
				COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0) as total_score,
				COUNT(pa.id) as badge_count,
				RANK() OVER (ORDER BY (COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0)) DESC) as rank
			FROM participants p
			LEFT JOIN event_registrations er ON p.id = er.participant_id
				AND er.event_id = (SELECT get_active_event())
			LEFT JOIN participant_achievements pa ON p.id = pa.participant_id
			LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
			GROUP BY p.id, p.naam, er.steps
		)
		SELECT
			id as participant_id,
			naam,
			rank,
			total_score,
			steps,
			achievement_points,
			badge_count
		FROM ranked_participants
		WHERE id = ?
	`

	err := r.DB().WithContext(ctx).Raw(query, participantID).Scan(&info).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("deelnemer niet gevonden")
	}
	if err != nil {
		return nil, err
	}

	// Haal participant boven me op - FIXED: Join with event_registrations
	var aboveMe models.LeaderboardEntry
	aboveQuery := `
		WITH ranked_participants AS (
			SELECT
				p.id,
				p.naam,
				COALESCE(er.distance_route, 'Unknown') as route,
				COALESCE(er.steps, 0) as steps,
				COALESCE(SUM(b.points), 0) as achievement_points,
				COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0) as total_score,
				COUNT(pa.id) as badge_count,
				COALESCE(er.registered_at, p.created_at) as joined_at,
				RANK() OVER (ORDER BY (COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0)) DESC) as rank
			FROM participants p
			LEFT JOIN event_registrations er ON p.id = er.participant_id
				AND er.event_id = (SELECT get_active_event())
			LEFT JOIN participant_achievements pa ON p.id = pa.participant_id
			LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
			GROUP BY p.id, p.naam, er.distance_route, er.steps, er.registered_at, p.created_at
		)
		SELECT *
		FROM ranked_participants
		WHERE rank = ? - 1
		LIMIT 1
	`

	err = r.DB().WithContext(ctx).Raw(aboveQuery, info.Rank).Scan(&aboveMe).Error
	if err == nil {
		info.AboveMe = &aboveMe
	}

	// Haal participant onder me op - FIXED: Join with event_registrations
	var belowMe models.LeaderboardEntry
	belowQuery := `
		WITH ranked_participants AS (
			SELECT
				p.id,
				p.naam,
				COALESCE(er.distance_route, 'Unknown') as route,
				COALESCE(er.steps, 0) as steps,
				COALESCE(SUM(b.points), 0) as achievement_points,
				COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0) as total_score,
				COUNT(pa.id) as badge_count,
				COALESCE(er.registered_at, p.created_at) as joined_at,
				RANK() OVER (ORDER BY (COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0)) DESC) as rank
			FROM participants p
			LEFT JOIN event_registrations er ON p.id = er.participant_id
				AND er.event_id = (SELECT get_active_event())
			LEFT JOIN participant_achievements pa ON p.id = pa.participant_id
			LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
			GROUP BY p.id, p.naam, er.distance_route, er.steps, er.registered_at, p.created_at
		)
		SELECT *
		FROM ranked_participants
		WHERE rank = ? + 1
		LIMIT 1
	`

	err = r.DB().WithContext(ctx).Raw(belowQuery, info.Rank).Scan(&belowMe).Error
	if err == nil {
		info.BelowMe = &belowMe
	}

	return &info, nil
}
