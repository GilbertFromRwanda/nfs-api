package userpermission

import "time"

type UserPermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	PermissionID uint      `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}
