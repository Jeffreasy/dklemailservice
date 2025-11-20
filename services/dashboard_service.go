package services

import (
	"context"
	"dklautomationgo/repository"
	"time"
)

// DashboardService aggregates data from various services for dashboard display
type DashboardService struct {
	participantRepo       repository.ParticipantRepository
	eventRepo             repository.EventRepository
	eventRegistrationRepo repository.EventRegistrationRepository
	contactRepo           repository.ContactRepository
	newsletterRepo        repository.NewsletterRepository
	userRepo              repository.GebruikerRepository
	emailMetrics          *EmailMetrics
	rateLimiter           *RateLimiter
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	participantRepo repository.ParticipantRepository,
	eventRepo repository.EventRepository,
	eventRegistrationRepo repository.EventRegistrationRepository,
	contactRepo repository.ContactRepository,
	newsletterRepo repository.NewsletterRepository,
	userRepo repository.GebruikerRepository,
	emailMetrics *EmailMetrics,
	rateLimiter *RateLimiter,
) *DashboardService {
	return &DashboardService{
		participantRepo:       participantRepo,
		eventRepo:             eventRepo,
		eventRegistrationRepo: eventRegistrationRepo,
		contactRepo:           contactRepo,
		newsletterRepo:        newsletterRepo,
		userRepo:              userRepo,
		emailMetrics:          emailMetrics,
		rateLimiter:           rateLimiter,
	}
}

// DashboardData represents the aggregated dashboard data
type DashboardData struct {
	// User statistics
	TotalUsers  int `json:"total_users"`
	ActiveUsers int `json:"active_users"`

	// Participant statistics
	TotalParticipants  int `json:"total_participants"`
	ActiveParticipants int `json:"active_participants"`

	// Event statistics
	TotalEvents        int `json:"total_events"`
	UpcomingEvents     int `json:"upcoming_events"`
	EventRegistrations int `json:"event_registrations"`

	// Contact/Inquiry statistics
	TotalContacts      int `json:"total_contacts"`
	UnansweredContacts int `json:"unanswered_contacts"`

	// Newsletter statistics
	TotalNewsletters int `json:"total_newsletters"`
	SentNewsletters  int `json:"sent_newsletters"`

	// System health
	SystemHealth SystemHealthData `json:"system_health"`

	// Recent activity (last 7 days)
	RecentActivity RecentActivityData `json:"recent_activity"`
}

// SystemHealthData contains system health metrics
type SystemHealthData struct {
	EmailMetrics    map[string]interface{} `json:"email_metrics"`
	RateLimitStatus map[string]int         `json:"rate_limit_status"`
	Uptime          string                 `json:"uptime"`
}

// RecentActivityData contains recent activity metrics
type RecentActivityData struct {
	NewUsers           int `json:"new_users"`
	NewParticipants    int `json:"new_participants"`
	NewContacts        int `json:"new_contacts"`
	EventRegistrations int `json:"event_registrations"`
	SentEmails         int `json:"sent_emails"`
}

// GetDashboardData aggregates and returns dashboard data
func (s *DashboardService) GetDashboardData(ctx context.Context) (*DashboardData, error) {
	data := &DashboardData{}

	// Get user statistics
	if err := s.getUserStats(ctx, data); err != nil {
		return nil, err
	}

	// Get participant statistics
	if err := s.getParticipantStats(ctx, data); err != nil {
		return nil, err
	}

	// Get event statistics
	if err := s.getEventStats(ctx, data); err != nil {
		return nil, err
	}

	// Get contact statistics
	if err := s.getContactStats(ctx, data); err != nil {
		return nil, err
	}

	// Get newsletter statistics
	if err := s.getNewsletterStats(ctx, data); err != nil {
		return nil, err
	}

	// Get system health data
	if err := s.getSystemHealthData(ctx, data); err != nil {
		return nil, err
	}

	// Get recent activity (last 7 days)
	if err := s.getRecentActivity(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *DashboardService) getUserStats(ctx context.Context, data *DashboardData) error {
	// Get total users
	users, err := s.userRepo.List(ctx, 10000, 0) // High limit to get all users
	if err != nil {
		return err
	}
	data.TotalUsers = len(users)

	// For simplicity, consider all users as active (you might want to add last_login tracking)
	data.ActiveUsers = data.TotalUsers

	return nil
}

func (s *DashboardService) getParticipantStats(ctx context.Context, data *DashboardData) error {
	// Get total participants
	participants, err := s.participantRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	data.TotalParticipants = len(participants)

	// Count active participants (those with recent activity)
	activeCount := 0
	for _, p := range participants {
		// Consider participants active if they have been updated in the last 30 days
		if time.Since(p.UpdatedAt) < 30*24*time.Hour {
			activeCount++
		}
	}
	data.ActiveParticipants = activeCount

	return nil
}

func (s *DashboardService) getEventStats(ctx context.Context, data *DashboardData) error {
	// Get total events
	events, err := s.eventRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	data.TotalEvents = len(events)

	// Count upcoming events
	upcomingCount := 0
	now := time.Now()
	for _, event := range events {
		if event.StartTime.After(now) {
			upcomingCount++
		}
	}
	data.UpcomingEvents = upcomingCount

	// Get event registrations
	registrations, err := s.eventRegistrationRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	data.EventRegistrations = len(registrations)

	return nil
}

func (s *DashboardService) getContactStats(ctx context.Context, data *DashboardData) error {
	// Get total contacts
	contacts, err := s.contactRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	data.TotalContacts = len(contacts)

	// Count unanswered contacts
	unansweredCount := 0
	for _, contact := range contacts {
		if contact.AntwoordTekst == "" {
			unansweredCount++
		}
	}
	data.UnansweredContacts = unansweredCount

	return nil
}

func (s *DashboardService) getNewsletterStats(ctx context.Context, data *DashboardData) error {
	// Get total newsletters
	newsletters, err := s.newsletterRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	data.TotalNewsletters = len(newsletters)

	// Count sent newsletters
	sentCount := 0
	for _, newsletter := range newsletters {
		if newsletter.SentAt != nil {
			sentCount++
		}
	}
	data.SentNewsletters = sentCount

	return nil
}

func (s *DashboardService) getSystemHealthData(ctx context.Context, data *DashboardData) error {
	// Get email metrics
	emailMetrics := make(map[string]interface{})
	emailMetrics["total_sent"] = s.emailMetrics.GetTotalEmails()
	emailMetrics["success_rate"] = s.emailMetrics.GetSuccessRate()
	emailMetrics["emails_by_type"] = s.emailMetrics.GetEmailsByType()
	data.SystemHealth.EmailMetrics = emailMetrics

	// Get rate limit status
	data.SystemHealth.RateLimitStatus = s.rateLimiter.GetCurrentValues()

	// Simple uptime (you might want to track actual service start time)
	data.SystemHealth.Uptime = "Service running"

	return nil
}

func (s *DashboardService) getRecentActivity(ctx context.Context, data *DashboardData) error {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// New users in last 7 days
	users, err := s.userRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	newUserCount := 0
	for _, user := range users {
		if user.CreatedAt.After(sevenDaysAgo) {
			newUserCount++
		}
	}
	data.RecentActivity.NewUsers = newUserCount

	// New participants in last 7 days
	participants, err := s.participantRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	newParticipantCount := 0
	for _, participant := range participants {
		if participant.CreatedAt.After(sevenDaysAgo) {
			newParticipantCount++
		}
	}
	data.RecentActivity.NewParticipants = newParticipantCount

	// New contacts in last 7 days
	contacts, err := s.contactRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	newContactCount := 0
	for _, contact := range contacts {
		if contact.CreatedAt.After(sevenDaysAgo) {
			newContactCount++
		}
	}
	data.RecentActivity.NewContacts = newContactCount

	// Event registrations in last 7 days
	registrations, err := s.eventRegistrationRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	recentRegistrationCount := 0
	for _, reg := range registrations {
		if reg.RegisteredAt.After(sevenDaysAgo) {
			recentRegistrationCount++
		}
	}
	data.RecentActivity.EventRegistrations = recentRegistrationCount

	// For sent emails, we can use the email metrics
	data.RecentActivity.SentEmails = int(s.emailMetrics.GetTotalEmails())

	return nil
}
