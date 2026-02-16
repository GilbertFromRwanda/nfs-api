package notary

type CreateRequest struct {
	Names    string `json:"names" binding:"required,max=200"`
	Phone    string `json:"phone" binding:"max=20"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
	OfficeID uint   `json:"office_id" binding:"required"`
}

type UpdateRequest struct {
	Names    *string `json:"names" binding:"omitempty,max=200"`
	Phone    *string `json:"phone" binding:"omitempty,max=20"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive"`
	OfficeID *uint   `json:"office_id"`
}
