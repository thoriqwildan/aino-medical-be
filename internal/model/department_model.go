package model

type DepartmentRequest struct {
	Name string `json:"name" validate:"required"`
}

type DepartmentResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
}

type UpdateDepartmentRequest struct {
	ID   uint `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}