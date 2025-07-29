package model

type DepartmentRequest struct {
	Name string `json:"name" validate:"required"`
}

type DepartmentResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateDepartmentRequest struct {
	ID   uint `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}