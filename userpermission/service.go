package userpermission

import (
	"errors"

	"nfs-api/database"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Assign(req AssignRequest) (*UserPermission, error) {
	var existing UserPermission
	if err := database.DB.Where("user_id = ? AND permission_id = ?", req.UserID, req.PermissionID).First(&existing).Error; err == nil {
		return nil, errors.New("permission already assigned to user")
	}

	up := UserPermission{
		UserID:       req.UserID,
		PermissionID: req.PermissionID,
	}

	if err := database.DB.Create(&up).Error; err != nil {
		return nil, errors.New("failed to assign permission")
	}
	return &up, nil
}

func (s *Service) GetByUserID(userID uint) ([]UserPermissionResponse, error) {
	var results []UserPermissionResponse
	err := database.DB.Table("user_permissions").
		Select("user_permissions.id, user_permissions.user_id, user_permissions.permission_id, permissions.name as permission_name").
		Joins("LEFT JOIN permissions ON permissions.id = user_permissions.permission_id").
		Where("user_permissions.user_id = ? AND permissions.deleted = 0", userID).
		Scan(&results).Error
	if err != nil {
		return nil, errors.New("failed to fetch user permissions")
	}
	return results, nil
}

func (s *Service) Revoke(id uint) error {
	result := database.DB.Delete(&UserPermission{}, id)
	if result.Error != nil {
		return errors.New("failed to revoke permission")
	}
	if result.RowsAffected == 0 {
		return errors.New("user permission not found")
	}
	return nil
}

func (s *Service) RevokeByUserAndPermission(userID, permissionID uint) error {
	result := database.DB.Where("user_id = ? AND permission_id = ?", userID, permissionID).Delete(&UserPermission{})
	if result.Error != nil {
		return errors.New("failed to revoke permission")
	}
	if result.RowsAffected == 0 {
		return errors.New("user permission not found")
	}
	return nil
}
func (s *Service) RevokeMany(userID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return errors.New("permission IDs list cannot be empty")
	}

	result := database.DB.Where("user_id = ? AND permission_id IN ?", userID, permissionIDs).Delete(&UserPermission{})
	if result.Error != nil {
		return errors.New("failed to revoke permissions")
	}
	if result.RowsAffected == 0 {
		return errors.New("no matching user permissions found")
	}
	return nil
}

func (s *Service) GetByID(id uint) (*UserPermission, error) {
	var up UserPermission
	if err := database.DB.First(&up, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user permission not found")
		}
		return nil, errors.New("failed to fetch user permission")
	}
	return &up, nil
}
func (s *Service) AssignMany(userID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return errors.New("permission IDs list cannot be empty")
	}

	// Check for existing permissions
	var existingPerms []UserPermission
	database.DB.Where("user_id = ? AND permission_id IN ?", userID, permissionIDs).Find(&existingPerms)

	// Create a map of existing permission IDs for quick lookup
	existingMap := make(map[uint]bool)
	for _, perm := range existingPerms {
		existingMap[perm.PermissionID] = true
	}

	// Prepare new permissions to insert (only non-existing ones)
	var newPermissions []UserPermission
	for _, permID := range permissionIDs {
		if !existingMap[permID] {
			newPermissions = append(newPermissions, UserPermission{
				UserID:       userID,
				PermissionID: permID,
			})
		}
	}

	// If all permissions already exist
	if len(newPermissions) == 0 {
		return errors.New("all permissions are already assigned to user")
	}

	// Bulk insert new permissions
	if err := database.DB.Create(&newPermissions).Error; err != nil {
		return errors.New("failed to assign permissions")
	}

	return nil
}
