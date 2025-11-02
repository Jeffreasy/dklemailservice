package logger

import (
	"context"
	"time"
)

// AuditEventType definieert het type audit event
type AuditEventType string

const (
	// Role management events
	AuditRoleAssigned AuditEventType = "ROLE_ASSIGNED"
	AuditRoleRevoked  AuditEventType = "ROLE_REVOKED"
	AuditRoleCreated  AuditEventType = "ROLE_CREATED"
	AuditRoleUpdated  AuditEventType = "ROLE_UPDATED"
	AuditRoleDeleted  AuditEventType = "ROLE_DELETED"

	// Permission management events
	AuditPermissionAssigned AuditEventType = "PERMISSION_ASSIGNED"
	AuditPermissionRevoked  AuditEventType = "PERMISSION_REVOKED"
	AuditPermissionCreated  AuditEventType = "PERMISSION_CREATED"
	AuditPermissionDeleted  AuditEventType = "PERMISSION_DELETED"

	// Access events
	AuditAccessGranted AuditEventType = "ACCESS_GRANTED"
	AuditAccessDenied  AuditEventType = "ACCESS_DENIED"

	// Authentication events
	AuditLoginSuccess    AuditEventType = "LOGIN_SUCCESS"
	AuditLoginFailed     AuditEventType = "LOGIN_FAILED"
	AuditLogout          AuditEventType = "LOGOUT"
	AuditTokenRefreshed  AuditEventType = "TOKEN_REFRESHED"
	AuditPasswordChanged AuditEventType = "PASSWORD_CHANGED"

	// User management events
	AuditUserCreated AuditEventType = "USER_CREATED"
	AuditUserUpdated AuditEventType = "USER_UPDATED"
	AuditUserDeleted AuditEventType = "USER_DELETED"

	// Cache events
	AuditCacheInvalidated AuditEventType = "CACHE_INVALIDATED"
)

// AuditEvent representeert een audit log event
type AuditEvent struct {
	Timestamp  time.Time              `json:"timestamp"`
	EventType  AuditEventType         `json:"event_type"`
	ActorID    string                 `json:"actor_id,omitempty"`    // Wie voert de actie uit
	ActorEmail string                 `json:"actor_email,omitempty"` // Email van actor
	TargetID   string                 `json:"target_id,omitempty"`   // Op wie/wat wordt de actie uitgevoerd
	TargetType string                 `json:"target_type,omitempty"` // Type van target (user, role, permission)
	ResourceID string                 `json:"resource_id,omitempty"` // ID van de resource
	Resource   string                 `json:"resource,omitempty"`    // Resource naam
	Action     string                 `json:"action,omitempty"`      // Action naam
	IPAddress  string                 `json:"ip_address,omitempty"`  // IP adres van requester
	UserAgent  string                 `json:"user_agent,omitempty"`  // User agent
	Result     string                 `json:"result"`                // Success/Failed/Denied
	Reason     string                 `json:"reason,omitempty"`      // Reden bij failure
	Metadata   map[string]interface{} `json:"metadata,omitempty"`    // Extra metadata
}

// Audit logt een audit event met structured data
func Audit(ctx context.Context, event AuditEvent) {
	// Zet timestamp als deze niet is ingesteld
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Log als structured event
	Info("AUDIT_EVENT",
		"event_type", event.EventType,
		"timestamp", event.Timestamp.Format(time.RFC3339),
		"actor_id", event.ActorID,
		"actor_email", event.ActorEmail,
		"target_id", event.TargetID,
		"target_type", event.TargetType,
		"resource_id", event.ResourceID,
		"resource", event.Resource,
		"action", event.Action,
		"ip_address", event.IPAddress,
		"user_agent", event.UserAgent,
		"result", event.Result,
		"reason", event.Reason,
		"metadata", event.Metadata,
	)
}

// AuditWithContext is een convenience functie die context gebruikt voor extra info
func AuditWithContext(ctx context.Context, eventType AuditEventType, actorID, targetID, result string, metadata map[string]interface{}) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: eventType,
		ActorID:   actorID,
		TargetID:  targetID,
		Result:    result,
		Metadata:  metadata,
	}

	Audit(ctx, event)
}

// Success result constants
const (
	ResultSuccess = "SUCCESS"
	ResultFailed  = "FAILED"
	ResultDenied  = "DENIED"
)
