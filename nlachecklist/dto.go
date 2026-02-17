package nlachecklist

import "nfs-api/utils"

type DossierItem struct {
	Name   string `json:"name" binding:"required,max=200"`
	Status string `json:"status" binding:"required,oneof=yes no no_needed"`
}

type CreateRequest struct {
	UPI          string        `json:"upi" binding:"required,max=40"`
	Seller1Names string        `json:"seller_1_names" binding:"required,max=150"`
	Seller2Names string        `json:"seller_2_names" binding:"omitempty,max=150"`
	Buyer1Names  string        `json:"buyer_1_names" binding:"required,max=150"`
	Buyer2Names  string        `json:"buyer_2_names" binding:"omitempty,max=150"`
	Date         string        `json:"date" binding:"required"` // YYYY-MM-DD
	Dossiers     []DossierItem `json:"dossiers" binding:"required,min=1,dive"`
}

type UpdateRequest struct {
	UPI          *string       `json:"upi" binding:"omitempty,max=40"`
	Seller1Names *string       `json:"seller_1_names" binding:"omitempty,max=150"`
	Seller2Names *string       `json:"seller_2_names" binding:"omitempty,max=150"`
	Buyer1Names  *string       `json:"buyer_1_names" binding:"omitempty,max=150"`
	Buyer2Names  *string       `json:"buyer_2_names" binding:"omitempty,max=150"`
	Date         *string       `json:"date"`
	Dossiers     []DossierItem `json:"dossiers" binding:"omitempty,dive"`
}

type FilterRequest struct {
	utils.PaginationRequest
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
	Status   string `form:"status" binding:"omitempty,oneof=completed uncompleted"`
	OfficeID uint   `form:"office_id"`
	UPI      string `form:"upi"`
}
