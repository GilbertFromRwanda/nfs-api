package appfeature

type CreateRequest struct {
	Name   string `json:"name" binding:"required,max=100"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateRequest struct {
	Name   *string `json:"name" binding:"omitempty,max=100"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}
