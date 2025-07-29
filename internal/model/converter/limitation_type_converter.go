package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func LimitationTypeToResponse(lt *entity.LimitationType) *model.LimitationTypeResponse {
	return &model.LimitationTypeResponse{
		ID:   lt.ID,
		Name: lt.Name,
	}
}