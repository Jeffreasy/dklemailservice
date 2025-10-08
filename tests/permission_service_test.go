package tests

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for testing
type MockRBACRoleRepository struct {
	mock.Mock
}

func (m *MockRBACRoleRepository) Create(ctx context.Context, role *models.RBACRole) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRBACRoleRepository) GetByID(ctx context.Context, id string) (*models.RBACRole, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RBACRole), args.Error(1)
}

func (m *MockRBACRoleRepository) GetByName(ctx context.Context, name string) (*models.RBACRole, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RBACRole), args.Error(1)
}

func (m *MockRBACRoleRepository) List(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.RBACRole), args.Error(1)
}

func (m *MockRBACRoleRepository) Update(ctx context.Context, role *models.RBACRole) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRBACRoleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRBACRoleRepository) GetSystemRoles(ctx context.Context) ([]*models.RBACRole, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.RBACRole), args.Error(1)
}

func (m *MockRBACRoleRepository) ListWithPermissions(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.RBACRole), args.Error(1)
}

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(ctx context.Context, permission *models.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByResourceAction(ctx context.Context, resource, action string) (*models.Permission, error) {
	args := m.Called(ctx, resource, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) List(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListByResource(ctx context.Context, resource string) ([]*models.Permission, error) {
	args := m.Called(ctx, resource)
	return args.Get(0).([]*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Update(ctx context.Context, permission *models.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetSystemPermissions(ctx context.Context) ([]*models.Permission, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Permission), args.Error(1)
}

type MockRolePermissionRepository struct {
	mock.Mock
}

func (m *MockRolePermissionRepository) Create(ctx context.Context, rp *models.RolePermission) error {
	args := m.Called(ctx, rp)
	return args.Error(0)
}

func (m *MockRolePermissionRepository) Delete(ctx context.Context, roleID, permissionID string) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockRolePermissionRepository) DeleteByRoleID(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockRolePermissionRepository) DeleteByPermissionID(ctx context.Context, permissionID string) error {
	args := m.Called(ctx, permissionID)
	return args.Error(0)
}

func (m *MockRolePermissionRepository) GetPermissionsByRole(ctx context.Context, roleID string) ([]*models.Permission, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]*models.Permission), args.Error(1)
}

func (m *MockRolePermissionRepository) GetRolesByPermission(ctx context.Context, permissionID string) ([]*models.RBACRole, error) {
	args := m.Called(ctx, permissionID)
	return args.Get(0).([]*models.RBACRole), args.Error(1)
}

func (m *MockRolePermissionRepository) HasPermission(ctx context.Context, roleID, permissionID string) (bool, error) {
	args := m.Called(ctx, roleID, permissionID)
	return args.Bool(0), args.Error(1)
}

type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) Create(ctx context.Context, ur *models.UserRole) error {
	args := m.Called(ctx, ur)
	return args.Error(0)
}

func (m *MockUserRoleRepository) GetByID(ctx context.Context, id string) (*models.UserRole, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) GetByUserAndRole(ctx context.Context, userID, roleID string) (*models.UserRole, error) {
	args := m.Called(ctx, userID, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) ListByUser(ctx context.Context, userID string) ([]*models.UserRole, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) ListByRole(ctx context.Context, roleID string) ([]*models.UserRole, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]*models.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) ListActiveByUser(ctx context.Context, userID string) ([]*models.UserRole, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) Update(ctx context.Context, ur *models.UserRole) error {
	args := m.Called(ctx, ur)
	return args.Error(0)
}

func (m *MockUserRoleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRoleRepository) DeleteByUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRoleRepository) DeleteByRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockUserRoleRepository) Deactivate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRoleRepository) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.UserPermission), args.Error(1)
}

func TestPermissionService_HasPermission(t *testing.T) {
	// Setup mocks
	mockRoleRepo := new(MockRBACRoleRepository)
	mockPermRepo := new(MockPermissionRepository)
	mockRolePermRepo := new(MockRolePermissionRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	// Create service
	service := services.NewPermissionService(
		mockRoleRepo,
		mockPermRepo,
		mockRolePermRepo,
		mockUserRoleRepo,
	)

	// Test data
	userID := "user-123"
	permissions := []*models.UserPermission{
		{Resource: "contact", Action: "read"},
		{Resource: "contact", Action: "write"},
		{Resource: "admin", Action: "access"},
	}

	// Setup expectations
	mockUserRoleRepo.On("GetUserPermissions", mock.Anything, userID).Return(permissions, nil)

	// Test cases
	tests := []struct {
		name     string
		resource string
		action   string
		expected bool
	}{
		{"Has contact read permission", "contact", "read", true},
		{"Has contact write permission", "contact", "write", true},
		{"Has admin access permission", "admin", "access", true},
		{"Does not have newsletter delete permission", "newsletter", "delete", false},
		{"Does not have user manage permission", "user", "manage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.HasPermission(context.Background(), userID, tt.resource, tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Verify expectations
	mockUserRoleRepo.AssertExpectations(t)
}

func TestPermissionService_AssignRole(t *testing.T) {
	// Setup mocks
	mockRoleRepo := new(MockRBACRoleRepository)
	mockPermRepo := new(MockPermissionRepository)
	mockRolePermRepo := new(MockRolePermissionRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)

	// Create service
	service := services.NewPermissionService(
		mockRoleRepo,
		mockPermRepo,
		mockRolePermRepo,
		mockUserRoleRepo,
	)

	// Test data
	userID := "user-123"
	roleID := "role-admin"
	assignedBy := "admin-user"

	// Mock role exists
	mockRoleRepo.On("GetByID", mock.Anything, roleID).Return(&models.RBACRole{ID: roleID, Name: "admin"}, nil)

	// Mock user doesn't have role yet
	mockUserRoleRepo.On("GetByUserAndRole", mock.Anything, userID, roleID).Return((*models.UserRole)(nil), nil)

	// Mock create user role
	mockUserRoleRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.UserRole")).Return(nil)

	// Test
	err := service.AssignRole(context.Background(), userID, roleID, &assignedBy)

	// Assert
	assert.NoError(t, err)

	// Verify expectations
	mockRoleRepo.AssertExpectations(t)
	mockUserRoleRepo.AssertExpectations(t)
}
