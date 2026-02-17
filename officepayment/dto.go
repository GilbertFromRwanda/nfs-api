package officepayment

import "nfs-api/utils"

type CreateRequest struct {
	OfficeID    uint    `json:"office_id" binding:"required"`
	PaidAmount  float64 `json:"paid_amount" binding:"required,min=0"`
	PaidMonth   int     `json:"paid_month" binding:"required,min=1"`
	PaymentDate string  `json:"payment_date" binding:"required"` // YYYY-MM-DD
	Comment     string  `json:"comment"`
}

type UpdateRequest struct {
	OfficeID    *uint    `json:"office_id"`
	PaidAmount  *float64 `json:"paid_amount" binding:"omitempty,min=0"`
	PaidMonth   *int     `json:"paid_month" binding:"omitempty,min=1"`
	PaymentDate *string  `json:"payment_date"`
	Comment     *string  `json:"comment"`
}

type FilterRequest struct {
	utils.PaginationRequest
	OfficeID uint   `form:"office_id"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}
