package permission

type CreateRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

type UpdateRequest struct {
	Name *string `json:"name" binding:"omitempty,max=100"`
}
