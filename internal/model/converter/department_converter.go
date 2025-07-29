package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func DepartmentToResponse(department *entity.Department) *model.DepartmentResponse {
	return &model.DepartmentResponse{
		ID:  department.ID,
		Name: department.Name,
		CreatedAt: department.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: department.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}