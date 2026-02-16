package agreement

import "nfs-api/utils"

type CreateRequest struct {
	Name           string   `json:"name" binding:"required,max=200"`
	NumberOfCopies int      `json:"number_of_copies" binding:"required,min=1"`
	Pricing        *float64 `json:"pricing"`
	Date           string   `json:"date" binding:"required"` // YYYY-MM-DD
	Content        string   `json:"content" binding:"required"`
	NotaryID       uint     `json:"notary_id" binding:"required"`
	OfficeID       uint     `json:"office_id" binding:"required"`
}

type UpdateRequest struct {
	Name           *string  `json:"name" binding:"omitempty,max=200"`
	NumberOfCopies *int     `json:"number_of_copies" binding:"omitempty,min=1"`
	Pricing        *float64 `json:"pricing"`
	Date           *string  `json:"date"`
	Content        *string  `json:"content"`
	NotaryID       *uint    `json:"notary_id"`
	OfficeID       *uint    `json:"office_id"`
}

type FilterRequest struct {
	utils.PaginationRequest
	DateFrom string `form:"date_from"` // YYYY-MM-DD
	DateTo   string `form:"date_to"`   // YYYY-MM-DD
	NotaryID uint   `form:"notary_id"`
	Name     string `form:"name"`
	OfficeID uint   `form:"office_id"`
}
