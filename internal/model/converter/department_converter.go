package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func DepartmentToResponse(department *entity.Department) *model.DepartmentResponse {
	return &model.DepartmentResponse{
		ID:  department.ID,
		Name: department.Name,
	}
}