package scanner

type EnrollRequest struct {
	Name                string `json:"name" binding:"required,max=200"`
	Nid                 string `json:"nid" binding:"required,max=50"`
	Phone               string `json:"phone" binding:"omitempty,max=20"`
	ServiceID           uint   `json:"service_id" binding:"required"`
	Service             string `json:"service" binding:"required,max=200"`
	FingerprintTemplate string `json:"fingerprint_template" binding:"required"`
	DeviceID            string `json:"device_id" binding:"omitempty,max=100"`
	// CompanyID and UserID are injected from JWT — not accepted from request body
	CompanyID uint `json:"-"`
	UserID    uint `json:"-"`
}

type UpdateFingerprintRequest struct {
	FingerprintTemplate string `json:"fingerprint_template" binding:"required"`
	DeviceID            string `json:"device_id" binding:"omitempty,max=100"`
}

type ListRequest struct {
	Q       string `form:"q"`
	Service string `form:"service"`
	From    string `form:"from"`
	To      string `form:"to"`
	Page    int    `form:"page,default=1" binding:"min=1"`
	PerPage int    `form:"per_page,default=50" binding:"min=1,max=500"`
	// CompanyID is injected from JWT — not accepted from query params
	CompanyID uint `form:"-"`
}

// TemplateEntry is returned by GET /scanner/templates — used for duplicate fingerprint detection during enroll.
type TemplateEntry struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Nid      string `json:"nid"`
	Template string `json:"template"`
}
