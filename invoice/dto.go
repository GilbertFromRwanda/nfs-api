package invoice

import "nfs-api/utils"

type CreateServiceItem struct {
	OfficeServiceID uint    `json:"office_service_id" binding:"required"`
	Name            string  `json:"name" binding:"required,max=100"`
	Copy            int     `json:"copy" binding:"required,min=1"`
	Price           float64 `json:"price" binding:"required,min=0"`
	PaymentStatus   string  `json:"payment_status" binding:"omitempty,oneof=paid unpaid"`
}

type CreateRequest struct {
	UPI          string              `json:"upi" binding:"required,max=40"`
	ClientName   string              `json:"client_name" binding:"required,max=150"`
	NotaryID     uint                `json:"notary_id" binding:"required"`
	Date         string              `json:"date" binding:"required"` // YYYY-MM-DD
	Note         string              `json:"note"`
	RefNo        string              `json:"ref_no" binding:"omitempty,max=50"`
	NotaryStatus string              `json:"notary_status" binding:"omitempty,oneof=pending signed_unpaid signed_paid"`
	Services     []CreateServiceItem `json:"services" binding:"required,min=1,dive"`
}

type UpdateRequest struct {
	UPI          *string `json:"upi" binding:"omitempty,max=40"`
	ClientName   *string `json:"client_name" binding:"omitempty,max=150"`
	NotaryID     *uint   `json:"notary_id"`
	Date         *string `json:"date"`
	Note         *string `json:"note"`
	RefNo        *string `json:"ref_no" binding:"omitempty,max=50"`
	NotaryStatus *string `json:"notary_status" binding:"omitempty,oneof=pending unpaid paid"`
}

type UpdateServiceStatusRequest struct {
	PaymentStatus string `json:"payment_status" binding:"required,oneof=paid unpaid"`
}

type FilterRequest struct {
	utils.PaginationRequest
	DateFrom     string `form:"date_from"`
	DateTo       string `form:"date_to"`
	UserID       uint   `form:"user_id"`
	OfficeID     uint   `form:"office_id"`
	ServiceIDs   string `form:"service_ids"` // comma-separated office_service_ids
	NotaryID     uint   `form:"notary_id"`
	UserStatus   string `form:"user_status" binding:"omitempty,oneof=paid unpaid"`
	NotaryStatus string `form:"notary_status" binding:"omitempty,oneof=pending signed_unpaid signed_paid"`
}

type UserServiceSummary struct {
	OfficeServiceID uint    `json:"office_service_id"`
	Name            string  `json:"name"`
	TotalAmount     float64 `json:"total_amount"`
	PaidAmount      float64 `json:"paid_amount"`
	UnpaidAmount    float64 `json:"unpaid_amount"`
}

type UserSummary struct {
	UserID       uint                 `json:"user_id"`
	FirstName    string               `json:"first_name"`
	LastName     string               `json:"last_name"`
	Username     string               `json:"username"`
	TotalAmount  float64              `json:"total_amount"`
	PaidAmount   float64              `json:"paid_amount"`
	UnpaidAmount float64              `json:"unpaid_amount"`
	Services     []UserServiceSummary `json:"services"`
}
