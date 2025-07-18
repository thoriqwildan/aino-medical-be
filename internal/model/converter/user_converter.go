package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID: 			user.ID,
		Username: 	user.Username,
		Name: 		user.Name,
		CreatedAt: 	user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}