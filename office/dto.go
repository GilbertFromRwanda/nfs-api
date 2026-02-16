package office

type CreateRequest struct {
	Name                string  `json:"name" binding:"required,max=100"`
	Phone               string  `json:"phone" binding:"max=20"`
	ScannerID           string  `json:"scanner_id" binding:"max=100"`
	SubscriptionEndDate *string `json:"subscription_end_date"`
	PaymentType         string  `json:"payment_type" binding:"required,oneof=subscription full_paid"`
	Pricing             float64 `json:"pricing"`
	Location            string  `json:"location" binding:"max=255"`
	Status              string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateRequest struct {
	Name                *string  `json:"name" binding:"omitempty,max=100"`
	Phone               *string  `json:"phone" binding:"omitempty,max=20"`
	ScannerID           *string  `json:"scanner_id" binding:"omitempty,max=100"`
	SubscriptionEndDate *string  `json:"subscription_end_date"`
	PaymentType         *string  `json:"payment_type" binding:"omitempty,oneof=subscription full_paid"`
	Pricing             *float64 `json:"pricing"`
	Location            *string  `json:"location" binding:"omitempty,max=255"`
	Status              *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}
