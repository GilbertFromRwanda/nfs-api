package auth

type RegisterRequest struct {
	FirstName string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string  `json:"last_name" binding:"required,min=2,max=100"`
	Username  string  `json:"username" binding:"required,min=3,max=100"`
	Password  string  `json:"password" binding:"required,min=6"`
	OfficeID  *uint   `json:"office_id"`
	Role      *string `json:"role"`
	Status    string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" binding:"omitempty,min=2,max=100"`
	Username  *string `json:"username" binding:"omitempty,min=3,max=100"`
	Password  *string `json:"password" binding:"omitempty,min=6"`
	OfficeID  *uint   `json:"office_id"`
	Role      *string `json:"role"`
	Status    *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
	OfficeID  *uint   `json:"office_id"`
	Role      *string `json:"role"`
	Status    string  `json:"status"`
}
