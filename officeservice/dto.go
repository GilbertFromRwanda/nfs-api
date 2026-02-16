package officeservice

type CreateRequest struct {
	Name           string `json:"name" binding:"required,max=100"`
	NumberOfCopies int     `json:"number_of_copies" binding:"required,min=1"`
	Pricing        float64 `json:"pricing"`
	OfficeID       uint    `json:"office_id" binding:"required"`
}

type UpdateRequest struct {
	Name           *string  `json:"name" binding:"omitempty,max=100"`
	NumberOfCopies *int     `json:"number_of_copies" binding:"omitempty,min=1"`
	Pricing        *float64 `json:"pricing"`
	OfficeID       *uint    `json:"office_id"`
}
