package services

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/repository"
)

// ChatServiceImpl implements the ChatService interface
type ChatServiceImpl struct {
	channelRepo     repository.ChatChannelRepository
	participantRepo repository.ChatChannelParticipantRepository
	messageRepo     repository.ChatMessageRepository
	reactionRepo    repository.ChatMessageReactionRepository
	presenceRepo    repository.ChatUserPresenceRepository
}

// NewChatService creates a new ChatService instance
func NewChatService(
	channelRepo repository.ChatChannelRepository,
	participantRepo repository.ChatChannelParticipantRepository,
	messageRepo repository.ChatMessageRepository,
	reactionRepo repository.ChatMessageReactionRepository,
	presenceRepo repository.ChatUserPresenceRepository,
) *ChatServiceImpl {
	return &ChatServiceImpl{
		channelRepo:     channelRepo,
		participantRepo: participantRepo,
		messageRepo:     messageRepo,
		reactionRepo:    reactionRepo,
		presenceRepo:    presenceRepo,
	}
}

// CreateChannel creates a new chat channel
func (s *ChatServiceImpl) CreateChannel(ctx context.Context, channel *models.ChatChannel) error {
	return s.channelRepo.Create(ctx, channel)
}

// GetChannel retrieves a chat channel by ID
func (s *ChatServiceImpl) GetChannel(ctx context.Context, id string) (*models.ChatChannel, error) {
	return s.channelRepo.GetByID(ctx, id)
}

// ListChannels lists chat channels with pagination
func (s *ChatServiceImpl) ListChannels(ctx context.Context, limit, offset int) ([]*models.ChatChannel, error) {
	return s.channelRepo.List(ctx, limit, offset)
}

// ListChannelsForUser lists channels for a specific user
func (s *ChatServiceImpl) ListChannelsForUser(ctx context.Context, userID string, limit, offset int) ([]*models.ChatChannel, error) {
	return s.channelRepo.ListByUserID(ctx, userID, limit, offset)
}

// UpdateChannel updates a chat channel
func (s *ChatServiceImpl) UpdateChannel(ctx context.Context, channel *models.ChatChannel) error {
	return s.channelRepo.Update(ctx, channel)
}

// DeleteChannel deletes a chat channel
func (s *ChatServiceImpl) DeleteChannel(ctx context.Context, id string) error {
	return s.channelRepo.Delete(ctx, id)
}

// AddParticipant adds a participant to a channel
func (s *ChatServiceImpl) AddParticipant(ctx context.Context, participant *models.ChatChannelParticipant) error {
	return s.participantRepo.Create(ctx, participant)
}

// GetParticipant retrieves a participant by ID
func (s *ChatServiceImpl) GetParticipant(ctx context.Context, id string) (*models.ChatChannelParticipant, error) {
	return s.participantRepo.GetByID(ctx, id)
}

// ListParticipants lists participants with pagination
func (s *ChatServiceImpl) ListParticipants(ctx context.Context, limit, offset int) ([]*models.ChatChannelParticipant, error) {
	return s.participantRepo.List(ctx, limit, offset)
}

// UpdateParticipant updates a participant
func (s *ChatServiceImpl) UpdateParticipant(ctx context.Context, participant *models.ChatChannelParticipant) error {
	return s.participantRepo.Update(ctx, participant)
}

// DeleteParticipant deletes a participant
func (s *ChatServiceImpl) DeleteParticipant(ctx context.Context, id string) error {
	return s.participantRepo.Delete(ctx, id)
}

// ListParticipantsByChannel lists participants by channel ID
func (s *ChatServiceImpl) ListParticipantsByChannel(ctx context.Context, channelID string) ([]*models.ChatChannelParticipant, error) {
	return s.participantRepo.ListByChannelID(ctx, channelID)
}

// CreateMessage creates a new message
func (s *ChatServiceImpl) CreateMessage(ctx context.Context, message *models.ChatMessage) error {
	return s.messageRepo.Create(ctx, message)
}

// GetMessage retrieves a message by ID
func (s *ChatServiceImpl) GetMessage(ctx context.Context, id string) (*models.ChatMessage, error) {
	return s.messageRepo.GetByID(ctx, id)
}

// ListMessages lists messages with pagination
func (s *ChatServiceImpl) ListMessages(ctx context.Context, limit, offset int) ([]*models.ChatMessage, error) {
	return s.messageRepo.List(ctx, limit, offset)
}

// UpdateMessage updates a message
func (s *ChatServiceImpl) UpdateMessage(ctx context.Context, message *models.ChatMessage) error {
	return s.messageRepo.Update(ctx, message)
}

// DeleteMessage deletes a message
func (s *ChatServiceImpl) DeleteMessage(ctx context.Context, id string) error {
	return s.messageRepo.Delete(ctx, id)
}

// ListMessagesByChannel lists messages by channel ID with pagination
func (s *ChatServiceImpl) ListMessagesByChannel(ctx context.Context, channelID string, limit, offset int) ([]*models.ChatMessage, error) {
	return s.messageRepo.ListByChannelID(ctx, channelID, limit, offset)
}

// AddReaction adds a reaction to a message
func (s *ChatServiceImpl) AddReaction(ctx context.Context, reaction *models.ChatMessageReaction) error {
	return s.reactionRepo.Create(ctx, reaction)
}

// GetReaction retrieves a reaction by ID
func (s *ChatServiceImpl) GetReaction(ctx context.Context, id string) (*models.ChatMessageReaction, error) {
	return s.reactionRepo.GetByID(ctx, id)
}

// ListReactions lists reactions with pagination
func (s *ChatServiceImpl) ListReactions(ctx context.Context, limit, offset int) ([]*models.ChatMessageReaction, error) {
	return s.reactionRepo.List(ctx, limit, offset)
}

// DeleteReaction deletes a reaction
func (s *ChatServiceImpl) DeleteReaction(ctx context.Context, id string) error {
	return s.reactionRepo.Delete(ctx, id)
}

// ListReactionsByMessage lists reactions by message ID
func (s *ChatServiceImpl) ListReactionsByMessage(ctx context.Context, messageID string) ([]*models.ChatMessageReaction, error) {
	return s.reactionRepo.ListByMessageID(ctx, messageID)
}

// UpdatePresence updates user presence
func (s *ChatServiceImpl) UpdatePresence(ctx context.Context, presence *models.ChatUserPresence) error {
	return s.presenceRepo.Upsert(ctx, presence)
}

// GetPresence retrieves user presence by user ID
func (s *ChatServiceImpl) GetPresence(ctx context.Context, userID string) (*models.ChatUserPresence, error) {
	return s.presenceRepo.GetByUserID(ctx, userID)
}

// DeletePresence deletes user presence
func (s *ChatServiceImpl) DeletePresence(ctx context.Context, userID string) error {
	return s.presenceRepo.Delete(ctx, userID)
}

// ListOnlineUsers lists online user IDs
func (s *ChatServiceImpl) ListOnlineUsers(ctx context.Context) ([]string, error) {
	return s.presenceRepo.ListOnlineUserIDs(ctx)
}
