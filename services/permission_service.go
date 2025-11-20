package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"encoding/json"
	"errors" // OPLOSSING: Importeer errors package
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm" // OPLOSSING: Importeer gorm
)

// PermissionServiceImpl implementeert de PermissionService interface
// V30+RBAC: Uitgebreid met participant repository voor app access checks
type PermissionServiceImpl struct {
	rbacRoleRepo       repository.RBACRoleRepository
	permissionRepo     repository.PermissionRepository
	rolePermissionRepo repository.RolePermissionRepository
	userRoleRepo       repository.UserRoleRepository
	participantRepo    repository.ParticipantRepository // V30+RBAC: Voor participant app access checks
	redisClient        *redis.Client
	cacheEnabled       bool
}

// NewPermissionService maakt een nieuwe PermissionService
func NewPermissionService(
	rbacRoleRepo repository.RBACRoleRepository,
	permissionRepo repository.PermissionRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	userRoleRepo repository.UserRoleRepository,
) PermissionService {
	return NewPermissionServiceWithRedis(rbacRoleRepo, permissionRepo, rolePermissionRepo, userRoleRepo, nil)
}

// NewPermissionServiceWithRedis maakt een nieuwe PermissionService met Redis ondersteuning
func NewPermissionServiceWithRedis(
	rbacRoleRepo repository.RBACRoleRepository,
	permissionRepo repository.PermissionRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	redisClient *redis.Client,
) PermissionService {
	return NewPermissionServiceWithParticipantSupport(rbacRoleRepo, permissionRepo, rolePermissionRepo, userRoleRepo, nil, redisClient)
}

// NewPermissionServiceWithParticipantSupport maakt een nieuwe PermissionService met volledige participant integratie
// V30+RBAC: Voegt participant repository toe voor app access validatie
func NewPermissionServiceWithParticipantSupport(
	rbacRoleRepo repository.RBACRoleRepository,
	permissionRepo repository.PermissionRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	participantRepo repository.ParticipantRepository,
	redisClient *redis.Client,
) PermissionService {
	impl := &PermissionServiceImpl{
		rbacRoleRepo:       rbacRoleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
		participantRepo:    participantRepo, // V30+RBAC
		redisClient:        redisClient,
		cacheEnabled:       redisClient != nil,
	}

	// Test Redis verbinding indien beschikbaar
	if impl.cacheEnabled {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			logger.Warn("Redis connection failed for PermissionService, disabling cache", "error", err)
			impl.cacheEnabled = false
		} else {
			logger.Info("Redis caching enabled for PermissionService")
		}
	}

	return impl
}

// HasPermission controleert of een gebruiker een specifieke permissie heeft
// OPLOSSING: Logica aangepast om correct onderscheid te maken tussen Participant en Gebruiker
func (s *PermissionServiceImpl) HasPermission(ctx context.Context, userID, resource, action string) bool {
	// Probeer eerst uit cache te halen
	if s.cacheEnabled {
		if cached := s.getCachedPermission(userID, resource, action); cached != nil {
			logger.Debug("Permission check from cache",
				"user_id", userID,
				"resource", resource,
				"action", action,
				"has_permission", *cached)
			return *cached
		}
	}

	var hasPermission bool

	// OPLOSSING: Eerst checken of het een participant is
	if s.participantRepo != nil {
		logger.Debug("Attempting participant lookup", "user_id", userID)
		participant, err := s.participantRepo.GetByID(ctx, userID)

		if err == nil && participant != nil {
			// ===================================
			// PAD A: Gebruiker IS een Participant
			// ===================================
			logger.Info("User is a participant, checking participant permissions",
				"user_id", userID,
				"resource", resource,
				"action", action,
				"account_type", participant.AccountType,
				"has_app_access", participant.HasAppAccess,
				"gebruiker_id", participant.GebruikerID)

			hasPermission = s.checkParticipantPermission(participant, resource, action)

			if !hasPermission {
				logger.Warn("Permission denied (participant)",
					"user_id", userID,
					"resource", resource,
					"action", action,
					"account_type", participant.AccountType,
					"has_app_access", participant.HasAppAccess)
			} else {
				logger.Debug("Permission granted (participant)",
					"user_id", userID,
					"resource", resource,
					"action", action)
			}

			// Cache het resultaat en return
			if s.cacheEnabled {
				s.cachePermission(userID, resource, action, hasPermission)
			}
			return hasPermission

		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Er is een ECHTE databasefout opgetreden bij het zoeken naar participant
			logger.Error("Fout bij ophalen participant in HasPermission", "user_id", userID, "error", err)
			return false // Fail closed
		}

		// Als we hier zijn, betekent het: err == gorm.ErrRecordNotFound
		// Dit is dus GEEN participant, en we gaan door naar de RBAC check hieronder.
		logger.Debug("Participant not found, proceeding with RBAC check", "user_id", userID)
	}

	// ===================================
	// PAD B: Gebruiker is een Gebruiker (Admin/Staff)
	// ===================================
	permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen user permissions (RBAC)", "user_id", userID, "error", err)
		return false
	}

	hasPermission = s.checkPermissionInList(permissions, resource, action)

	// Log altijd voor debugging van admin toegang problemen
	if !hasPermission {
		// Extra logging voor admin access problemen
		if resource == "admin" && action == "access" {
			roles, _ := s.GetUserRoles(ctx, userID)
			roleNames := make([]string, len(roles))
			for i, role := range roles {
				roleNames[i] = role.Role.Name
			}
			logger.Warn("ADMIN ACCESS DENIED - Detailed analysis",
				"user_id", userID,
				"resource", resource,
				"action", action,
				"permissions_count", len(permissions),
				"roles_count", len(roles),
				"role_names", roleNames,
				"available_permissions", s.formatPermissionsList(permissions))
		} else {
			logger.Warn("Permission denied (RBAC)",
				"user_id", userID,
				"resource", resource,
				"action", action,
				"permissions_count", len(permissions),
				"available_permissions", s.formatPermissionsList(permissions))
		}
	} else {
		logger.Debug("Permission granted (RBAC)",
			"user_id", userID,
			"resource", resource,
			"action", action)
	}

	// Cache het resultaat
	if s.cacheEnabled {
		s.cachePermission(userID, resource, action, hasPermission)
	}

	return hasPermission
}

// checkPermissionInList controleert of een permissie in de lijst staat
func (s *PermissionServiceImpl) checkPermissionInList(permissions []*models.UserPermission, resource, action string) bool {
	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}

// formatPermissionsList formatteert een lijst van permissies voor logging
func (s *PermissionServiceImpl) formatPermissionsList(permissions []*models.UserPermission) string {
	if len(permissions) == 0 {
		return "none"
	}

	var permStrings []string
	for _, perm := range permissions {
		permStrings = append(permStrings, fmt.Sprintf("%s:%s", perm.Resource, perm.Action))
	}

	return strings.Join(permStrings, ", ")
}

// checkParticipantPermission controleert participant permissions (V34)
// OPLOSSING 2: Gebruikt de nieuwe centrale permissie-map
func (s *PermissionServiceImpl) checkParticipantPermission(participant *models.Participant, resource, action string) bool {
	// Check 1: Must have app access
	if !participant.HasAppAccess {
		return false
	}

	// Check 2: Must be full account
	if participant.AccountType != "full" {
		return false
	}

	// OPLOSSING 2: Check permissies tegen de centrale map
	// Speciale mapping: 'profile:read' en 'profile:update' zijn aliassen
	// voor 'participant:read' en 'participant:write'
	effectiveResource := resource
	effectiveAction := action
	if resource == "profile" {
		effectiveResource = "participant"
		if action == "update" {
			effectiveAction = "update_own" // Gecorrigeerd naar de juiste 'action'
		}
		if action == "read" {
			effectiveAction = "view_own" // Gecorrigeerd naar de juiste 'action'
		}
	}

	// Check if resource exists
	actions, exists := participantPermissionsMap[effectiveResource]
	if !exists {
		return false
	}

	// Check if action is allowed
	for _, allowedAction := range actions {
		if allowedAction == effectiveAction {
			return true
		}
	}

	return false
}

// GetUserPermissions haalt alle permissies op voor een gebruiker
// RBAC: Only uses user_roles system (legacy role_id removed)
func (s *PermissionServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	// Get permissions from user_roles table only
	permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("Error getting user permissions from user_roles", "user_id", userID, "error", err)
		return nil, err
	}

	logger.Debug("Retrieved permissions from user_roles", "user_id", userID, "permissions_count", len(permissions))
	return permissions, nil
}

// GetUserRoles haalt alle actieve rollen op voor een gebruiker
func (s *PermissionServiceImpl) GetUserRoles(ctx context.Context, userID string) ([]*models.UserRole, error) {
	return s.userRoleRepo.ListActiveByUser(ctx, userID)
}

// AssignRole kent een rol toe aan een gebruiker
func (s *PermissionServiceImpl) AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error {
	// Controleer of de rol bestaat
	role, err := s.rbacRoleRepo.GetByID(ctx, roleID)
	if err != nil {
		// Audit: Failed role assignment
		actorID := ""
		if assignedBy != nil {
			actorID = *assignedBy
		}
		logger.Audit(ctx, logger.AuditEvent{
			EventType:  logger.AuditRoleAssigned,
			ActorID:    actorID,
			TargetID:   userID,
			TargetType: "user",
			ResourceID: roleID,
			Resource:   "role",
			Result:     logger.ResultFailed,
			Reason:     "rol niet gevonden",
		})
		return fmt.Errorf("rol niet gevonden: %w", err)
	}

	// Controleer of de gebruiker de rol al heeft
	existing, err := s.userRoleRepo.GetByUserAndRole(ctx, userID, roleID)
	if err == nil && existing != nil && existing.IsActive {
		return fmt.Errorf("gebruiker heeft deze rol al")
	}

	// Maak nieuwe user-role relatie aan
	userRole := &models.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedAt: time.Now(),
		AssignedBy: assignedBy,
		IsActive:   true,
	}

	if err := s.userRoleRepo.Create(ctx, userRole); err != nil {
		// Audit: Failed role assignment
		actorID := ""
		if assignedBy != nil {
			actorID = *assignedBy
		}
		logger.Audit(ctx, logger.AuditEvent{
			EventType:  logger.AuditRoleAssigned,
			ActorID:    actorID,
			TargetID:   userID,
			TargetType: "user",
			ResourceID: roleID,
			Resource:   "role",
			Result:     logger.ResultFailed,
			Reason:     err.Error(),
		})
		return fmt.Errorf("fout bij toekennen rol: %w", err)
	}

	// Invalideer cache
	s.InvalidateUserCache(userID)

	// Audit: Successful role assignment
	actorID := ""
	if assignedBy != nil {
		actorID = *assignedBy
	}
	logger.Audit(ctx, logger.AuditEvent{
		EventType:  logger.AuditRoleAssigned,
		ActorID:    actorID,
		TargetID:   userID,
		TargetType: "user",
		ResourceID: roleID,
		Resource:   "role",
		Result:     logger.ResultSuccess,
		Metadata: map[string]interface{}{
			"role_name": role.Name,
		},
	})

	logger.Info("Rol toegekend aan gebruiker", "user_id", userID, "role_id", roleID, "role_name", role.Name, "assigned_by", assignedBy)
	return nil
}

// RevokeRole verwijdert een rol van een gebruiker
func (s *PermissionServiceImpl) RevokeRole(ctx context.Context, userID, roleID string) error {
	// Controleer of de relatie bestaat
	existing, err := s.userRoleRepo.GetByUserAndRole(ctx, userID, roleID)
	if err != nil || existing == nil {
		return fmt.Errorf("gebruiker heeft deze rol niet")
	}

	// Deactiveer de relatie
	if err := s.userRoleRepo.Deactivate(ctx, existing.ID); err != nil {
		return fmt.Errorf("fout bij verwijderen rol: %w", err)
	}

	// Invalideer cache
	s.InvalidateUserCache(userID)

	// Audit: Role revoked
	logger.Audit(ctx, logger.AuditEvent{
		EventType:  logger.AuditRoleRevoked,
		TargetID:   userID,
		TargetType: "user",
		ResourceID: roleID,
		Resource:   "role",
		Result:     logger.ResultSuccess,
		Metadata: map[string]interface{}{
			"role_name": existing.Role.Name,
		},
	})

	logger.Info("Rol verwijderd van gebruiker", "user_id", userID, "role_id", roleID)
	return nil
}

// CreateRole maakt een nieuwe rol aan
func (s *PermissionServiceImpl) CreateRole(ctx context.Context, role *models.RBACRole, createdBy *string) error {
	role.CreatedBy = createdBy
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	if err := s.rbacRoleRepo.Create(ctx, role); err != nil {
		return fmt.Errorf("fout bij aanmaken rol: %w", err)
	}

	logger.Info("Nieuwe rol aangemaakt", "role_name", role.Name, "created_by", createdBy)
	return nil
}

// UpdateRole werkt een rol bij
func (s *PermissionServiceImpl) UpdateRole(ctx context.Context, role *models.RBACRole) error {
	// Controleer of rol systeemrol is en niet verwijderd mag worden
	existing, err := s.rbacRoleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return fmt.Errorf("rol niet gevonden: %w", err)
	}

	if existing.IsSystemRole && !role.IsSystemRole {
		return fmt.Errorf("kan systeemrol niet wijzigen naar niet-systeemrol")
	}

	role.UpdatedAt = time.Now()

	if err := s.rbacRoleRepo.Update(ctx, role); err != nil {
		return fmt.Errorf("fout bij bijwerken rol: %w", err)
	}

	// Refresh cache voor alle gebruikers met deze rol
	s.refreshUsersWithRole(ctx, role.ID)

	logger.Info("Rol bijgewerkt", "role_id", role.ID, "role_name", role.Name)
	return nil
}

// DeleteRole verwijdert een rol
func (s *PermissionServiceImpl) DeleteRole(ctx context.Context, roleID string) error {
	// Controleer of het een systeemrol is
	role, err := s.rbacRoleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("rol niet gevonden: %w", err)
	}

	if role.IsSystemRole {
		return fmt.Errorf("kan systeemrol niet verwijderen")
	}

	// Verwijder alle user-role relaties
	if err := s.userRoleRepo.DeleteByRole(ctx, roleID); err != nil {
		return fmt.Errorf("fout bij verwijderen user-role relaties: %w", err)
	}

	// Verwijder alle role-permission relaties
	if err := s.rolePermissionRepo.DeleteByRoleID(ctx, roleID); err != nil {
		return fmt.Errorf("fout bij verwijderen role-permission relaties: %w", err)
	}

	// Verwijder de rol
	if err := s.rbacRoleRepo.Delete(ctx, roleID); err != nil {
		return fmt.Errorf("fout bij verwijderen rol: %w", err)
	}

	// Refresh cache voor alle gebruikers
	s.RefreshCache(ctx)

	logger.Info("Rol verwijderd", "role_id", roleID, "role_name", role.Name)
	return nil
}

// AssignPermissionToRole kent een permissie toe aan een rol
func (s *PermissionServiceImpl) AssignPermissionToRole(ctx context.Context, roleID, permissionID string, assignedBy *string) error {
	// Controleer of rol en permissie bestaan
	_, err := s.rbacRoleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("rol niet gevonden: %w", err)
	}

	_, err = s.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("permissie niet gevonden: %w", err)
	}

	// Controleer of de relatie al bestaat
	exists, err := s.rolePermissionRepo.HasPermission(ctx, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("fout bij controleren bestaande relatie: %w", err)
	}
	if exists {
		return fmt.Errorf("rol heeft deze permissie al")
	}

	// Maak nieuwe relatie aan
	rp := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		AssignedAt:   time.Now(),
		AssignedBy:   assignedBy,
	}

	if err := s.rolePermissionRepo.Create(ctx, rp); err != nil {
		return fmt.Errorf("fout bij toekennen permissie: %w", err)
	}

	// Refresh cache voor alle gebruikers met deze rol
	s.refreshUsersWithRole(ctx, roleID)

	logger.Info("Permissie toegekend aan rol", "role_id", roleID, "permission_id", permissionID, "assigned_by", assignedBy)
	return nil
}

// RevokePermissionFromRole verwijdert een permissie van een rol
func (s *PermissionServiceImpl) RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	if err := s.rolePermissionRepo.Delete(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("fout bij verwijderen permissie: %w", err)
	}

	// Refresh cache voor alle gebruikers met deze rol
	s.refreshUsersWithRole(ctx, roleID)

	logger.Info("Permissie verwijderd van rol", "role_id", roleID, "permission_id", permissionID)
	return nil
}

// GetRoles haalt alle rollen op met validatie
func (s *PermissionServiceImpl) GetRoles(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	// Valideer en normaliseer limit
	if limit <= 0 {
		limit = 100 // Default
	}
	if limit > 1000 {
		limit = 1000 // Max voor performance
		logger.Warn("GetRoles limit overschreden, beperkt tot 1000", "requested", limit)
	}

	return s.rbacRoleRepo.List(ctx, limit, offset)
}

// GetPermissions haalt alle permissies op met validatie
func (s *PermissionServiceImpl) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	// Valideer en normaliseer limit
	if limit <= 0 {
		limit = 100 // Default
	}
	if limit > 1000 {
		limit = 1000 // Max voor performance
		logger.Warn("GetPermissions limit overschreden, beperkt tot 1000", "requested", limit)
	}

	return s.permissionRepo.List(ctx, limit, offset)
}

// GetRoleByName haalt een rol op basis van naam
func (s *PermissionServiceImpl) GetRoleByName(ctx context.Context, name string) (*models.RBACRole, error) {
	role, err := s.rbacRoleRepo.GetByName(ctx, name)
	if err != nil {
		logger.Error("Fout bij ophalen rol op naam", "name", name, "error", err)
		return nil, err
	}
	return role, nil
}

// GetPermissionsByRole haalt alle permissies op voor een specifieke rol
func (s *PermissionServiceImpl) GetPermissionsByRole(ctx context.Context, roleID string) ([]*models.Permission, error) {
	permissions, err := s.rolePermissionRepo.GetPermissionsByRole(ctx, roleID)
	if err != nil {
		logger.Error("Fout bij ophalen permissies voor rol", "role_id", roleID, "error", err)
		return nil, err
	}
	return permissions, nil
}

// getCachedPermission haalt een permissie uit de Redis cache
func (s *PermissionServiceImpl) getCachedPermission(userID, resource, action string) *bool {
	if !s.cacheEnabled {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("perm:%s:%s:%s", userID, resource, action)

	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil // Not in cache
	}
	if err != nil {
		logger.Error("Redis cache get error", "error", err, "key", cacheKey)
		return nil
	}

	var hasPermission bool
	if err := json.Unmarshal([]byte(val), &hasPermission); err != nil {
		logger.Error("Redis cache unmarshal error", "error", err, "key", cacheKey)
		return nil
	}

	return &hasPermission
}

// cachePermission slaat een permissie op in de Redis cache
func (s *PermissionServiceImpl) cachePermission(userID, resource, action string, hasPermission bool) {
	if !s.cacheEnabled {
		return
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("perm:%s:%s:%s", userID, resource, action)

	data, err := json.Marshal(hasPermission)
	if err != nil {
		logger.Error("Redis cache marshal error", "error", err, "key", cacheKey)
		return
	}

	// Cache voor 5 minuten voor snellere updates
	err = s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute).Err()
	if err != nil {
		logger.Error("Redis cache set error", "error", err, "key", cacheKey)
	}
}

// InvalidateUserCache wist de cache voor een gebruiker
func (s *PermissionServiceImpl) InvalidateUserCache(userID string) {
	if !s.cacheEnabled {
		logger.Debug("User cache invalidation requested but no caching enabled", "user_id", userID)
		return
	}

	ctx := context.Background()
	pattern := fmt.Sprintf("perm:%s:*", userID)

	// Haal alle keys op die matchen met het patroon
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		logger.Error("Redis keys error during cache invalidation", "error", err, "pattern", pattern)
		return
	}

	if len(keys) > 0 {
		err = s.redisClient.Del(ctx, keys...).Err()
		if err != nil {
			logger.Error("Redis del error during cache invalidation", "error", err, "keys", keys)
			return
		}
		logger.Debug("User cache invalidated", "user_id", userID, "keys_deleted", len(keys))
	}
}

// RefreshCache vernieuwt alle caches
func (s *PermissionServiceImpl) RefreshCache(ctx context.Context) error {
	if !s.cacheEnabled {
		logger.Debug("Cache refresh requested but no caching enabled")
		return nil
	}

	// Verwijder alle permissie caches
	pattern := "perm:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		logger.Error("Redis keys error during full cache refresh", "error", err, "pattern", pattern)
		return err
	}

	if len(keys) > 0 {
		err = s.redisClient.Del(ctx, keys...).Err()
		if err != nil {
			logger.Error("Redis del error during full cache refresh", "error", err, "keys", keys)
			return err
		}
		logger.Info("Full cache refresh completed", "keys_deleted", len(keys))
	}

	return nil
}

// refreshUsersWithRole vernieuwt de cache voor alle gebruikers met een specifieke rol
func (s *PermissionServiceImpl) refreshUsersWithRole(ctx context.Context, roleID string) {
	if !s.cacheEnabled {
		logger.Debug("User role cache refresh requested but no caching enabled", "role_id", roleID)
		return
	}

	// Haal alle gebruikers op die deze rol hebben
	userRoles, err := s.userRoleRepo.ListByRole(ctx, roleID)
	if err != nil {
		logger.Error("Error getting users with role for cache refresh", "error", err, "role_id", roleID)
		return
	}

	// Invalideer cache voor elke gebruiker
	for _, userRole := range userRoles {
		if userRole.IsActive { // Alleen actieve rollen
			s.InvalidateUserCache(userRole.UserID)
		}
	}

	logger.Debug("User role cache refresh completed", "role_id", roleID, "users_affected", len(userRoles))
}

// ==============================================================================
// V30+RBAC: PARTICIPANT-SPECIFIEKE PERMISSION CHECKS
// ==============================================================================

// HasParticipantAppAccess controleert of een gebruiker app toegang heeft via participant account
// V30+RBAC: Combineert RBAC permission check met participant account type check
func (s *PermissionServiceImpl) HasParticipantAppAccess(ctx context.Context, userID string) bool {
	// Check 1: RBAC permission check - heeft user de 'app:access' permission?
	if !s.HasPermission(ctx, userID, "app", "access") {
		logger.Debug("User heeft geen app:access permission", "user_id", userID)
		return false
	}

	// Check 2: Participant account type check (alleen als participant repo beschikbaar)
	if s.participantRepo == nil {
		logger.Debug("Participant repository niet beschikbaar, allow access op basis van permission alleen", "user_id", userID)
		return true
	}

	// Haal gebruiker email op
	// Dit vereist een lookup - we kunnen dit optimaliseren door email in JWT claims te hebben
	// Voor nu doen we het via participant lookup

	// TODO: Implementeer GetByGebruikerID in participant repository voor betere performance
	// Voor nu: return true als permission check geslaagd is (participant check gebeurt al in auth service login)
	return true
}

// CanParticipantRegisterForEvent controleert of een participant zich kan registreren voor een event
// V30+RBAC: Zowel full als temporary accounts kunnen zich registreren, maar met verschillende rechten
func (s *PermissionServiceImpl) CanParticipantRegisterForEvent(ctx context.Context, userID string) bool {
	// Full account users kunnen altijd registreren (via participant:register_event permission)
	if s.HasPermission(ctx, userID, "participant", "register_event") {
		return true
	}

	// Voor temporary accounts: public endpoint heeft geen user ID, dus dit is niet van toepassing
	// Temporary registratie gebeurt via public endpoint zonder authenticatie
	return false
}

// GetParticipantPermissionLevel bepaalt het permission level van een participant
// V30+RBAC: Returns 'full', 'temporary', of 'none'
func (s *PermissionServiceImpl) GetParticipantPermissionLevel(ctx context.Context, userID string) string {
	if s.participantRepo == nil {
		// Geen participant repo = legacy user
		if s.HasPermission(ctx, userID, "app", "access") {
			return "full"
		}
		return "none"
	}

	// Check app access permission
	if !s.HasPermission(ctx, userID, "app", "access") {
		return "none"
	}

	return "full"
}
