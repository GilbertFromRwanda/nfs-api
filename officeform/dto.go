package officeform

type CreateRequest struct {
	Title    string `json:"title" binding:"required,max=200"`
	Content  string `json:"content" binding:"required"`
	OfficeID uint   `json:"office_id" binding:"required"`
}

type UpdateRequest struct {
	Title    *string `json:"title" binding:"omitempty,max=200"`
	Content  *string `json:"content"`
	OfficeID *uint   `json:"office_id"`
}
