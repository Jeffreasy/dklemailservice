package tests

import (
	"context"
	"dklautomationgo/models"

	"github.com/stretchr/testify/mock"
)

// ============================================================================
// COMPLETE MOCK FOR PERMISSION SERVICE
// ============================================================================

type AuthMockPermissionService struct {
	mock.Mock
}

func (m *AuthMockPermissionService) HasPermission(ctx context.Context, userID, resource, action string) bool {
	args := m.Called(ctx, userID, resource, action)
	return args.Bool(0)
}

func (m *AuthMockPermissionService) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserPermission), args.Error(1)
}

func (m *AuthMockPermissionService) GetUserRoles(ctx context.Context, userID string) ([]*models.UserRole, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserRole), args.Error(1)
}

func (m *AuthMockPermissionService) AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error {
	args := m.Called(ctx, userID, roleID, assignedBy)
	return args.Error(0)
}

func (m *AuthMockPermissionService) RevokeRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *AuthMockPermissionService) CreateRole(ctx context.Context, role *models.RBACRole, createdBy *string) error {
	args := m.Called(ctx, role, createdBy)
	return args.Error(0)
}

func (m *AuthMockPermissionService) UpdateRole(ctx context.Context, role *models.RBACRole) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *AuthMockPermissionService) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *AuthMockPermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string, assignedBy *string) error {
	args := m.Called(ctx, roleID, permissionID, assignedBy)
	return args.Error(0)
}

func (m *AuthMockPermissionService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *AuthMockPermissionService) GetRoles(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.RBACRole), args.Error(1)
}

func (m *AuthMockPermissionService) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Permission), args.Error(1)
}

func (m *AuthMockPermissionService) InvalidateUserCache(userID string) {
	m.Called(userID)
}

func (m *AuthMockPermissionService) RefreshCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
