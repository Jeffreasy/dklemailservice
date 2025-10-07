package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PermissionServiceImpl implementeert de PermissionService interface
type PermissionServiceImpl struct {
	rbacRoleRepo       repository.RBACRoleRepository
	permissionRepo     repository.PermissionRepository
	rolePermissionRepo repository.RolePermissionRepository
	userRoleRepo       repository.UserRoleRepository
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
	impl := &PermissionServiceImpl{
		rbacRoleRepo:       rbacRoleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
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
func (s *PermissionServiceImpl) HasPermission(ctx context.Context, userID, resource, action string) bool {
	// Probeer eerst uit cache te halen
	if s.cacheEnabled {
		if cached := s.getCachedPermission(userID, resource, action); cached != nil {
			return *cached
		}
	}

	// Haal permissies op uit database
	permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen user permissions", "user_id", userID, "error", err)
		return false
	}

	// Debug logging voor admin user
	if userID == "7157f3f6-da85-4058-9d38-19133ec93b03" { // admin user ID from logs
		logger.Info("Admin user permissions check", "user_id", userID, "resource", resource, "action", action, "permissions_count", len(permissions))
		for _, perm := range permissions {
			logger.Info("Admin user permission", "resource", perm.Resource, "action", perm.Action, "role", perm.RoleName)
		}
	}

	hasPermission := s.checkPermissionInList(permissions, resource, action)

	// Debug logging voor alle permissie checks
	logger.Info("Permission check result",
		"user_id", userID,
		"resource", resource,
		"action", action,
		"has_permission", hasPermission,
		"permissions_count", len(permissions))

	// Gedetailleerde logging voor admin user of bij permission denied
	if userID == "7157f3f6-da85-4058-9d38-19133ec93b03" || !hasPermission {
		logger.Info("Detailed permission analysis",
			"user_id", userID,
			"resource", resource,
			"action", action,
			"permissions_found", len(permissions))

		for _, perm := range permissions {
			logger.Info("User permission",
				"user_id", userID,
				"resource", perm.Resource,
				"action", perm.Action,
				"role_name", perm.RoleName)
		}
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

// GetUserPermissions haalt alle permissies op voor een gebruiker
func (s *PermissionServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	return s.userRoleRepo.GetUserPermissions(ctx, userID)
}

// GetUserRoles haalt alle actieve rollen op voor een gebruiker
func (s *PermissionServiceImpl) GetUserRoles(ctx context.Context, userID string) ([]*models.UserRole, error) {
	return s.userRoleRepo.ListActiveByUser(ctx, userID)
}

// AssignRole kent een rol toe aan een gebruiker
func (s *PermissionServiceImpl) AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error {
	// Controleer of de rol bestaat
	_, err := s.rbacRoleRepo.GetByID(ctx, roleID)
	if err != nil {
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
		return fmt.Errorf("fout bij toekennen rol: %w", err)
	}

	// Invalideer cache
	s.InvalidateUserCache(userID)

	logger.Info("Rol toegekend aan gebruiker", "user_id", userID, "role_id", roleID, "assigned_by", assignedBy)
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

// GetRoles haalt alle rollen op
func (s *PermissionServiceImpl) GetRoles(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	return s.rbacRoleRepo.List(ctx, limit, offset)
}

// GetPermissions haalt alle permissies op
func (s *PermissionServiceImpl) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	return s.permissionRepo.List(ctx, limit, offset)
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

	// Cache voor 10 minuten
	err = s.redisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()
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
