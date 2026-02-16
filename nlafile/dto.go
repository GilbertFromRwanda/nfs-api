package nlafile

import "nfs-api/utils"

type CreateRequest struct {
	OfficeServiceID  uint   `json:"office_service_id" binding:"required"`
	Seller           string `json:"seller" binding:"required,max=200"`
	Buyer            string `json:"buyer" binding:"required,max=200"`
	Date             string `json:"date" binding:"required"` // YYYY-MM-DD
	PendingPeriodDays int   `json:"pending_period_days" binding:"required,min=1"`
	UPI              string `json:"upi" binding:"required,max=100"`
	OfficeID         uint   `json:"office_id" binding:"required"`
}

type UpdateRequest struct {
	OfficeServiceID  *uint   `json:"office_service_id"`
	Seller           *string `json:"seller" binding:"omitempty,max=200"`
	Buyer            *string `json:"buyer" binding:"omitempty,max=200"`
	Date             *string `json:"date"`
	PendingPeriodDays *int   `json:"pending_period_days" binding:"omitempty,min=1"`
	UPI              *string `json:"upi" binding:"omitempty,max=100"`
	OfficeID         *uint   `json:"office_id"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type FilterRequest struct {
	utils.PaginationRequest
	UPI      string `form:"upi"`
	Name     string `form:"name"` // searches seller and buyer
	OfficeID uint   `form:"office_id"`
	Status   string `form:"status"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}
