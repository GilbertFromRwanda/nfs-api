package userpermission

type AssignRequest struct {
	UserID       uint `json:"user_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
}

type UserPermissionResponse struct {
	ID             uint   `json:"id"`
	UserID         uint   `json:"user_id"`
	PermissionID   uint   `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}
