package checklistnote

import "nfs-api/utils"

type CreateRequest struct {
	Name     string `json:"name" binding:"required"`
	OfficeID uint   `json:"office_id" binding:"required"`
}

type UpdateRequest struct {
	Name     *string `json:"name"`
	OfficeID *uint   `json:"office_id"`
}

type FilterRequest struct {
	utils.PaginationRequest
	OfficeID uint `form:"office_id"`
}
