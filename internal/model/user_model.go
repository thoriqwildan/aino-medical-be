package model

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Name 	 *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

type UserResponse struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Name     *string `json:"name,omitempty"`
	CreatedAt string  `json:"created_at"`
}