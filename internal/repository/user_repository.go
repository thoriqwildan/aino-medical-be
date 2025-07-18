package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (ur *UserRepository) GetByUsername(db *gorm.DB, username string, user *entity.User) error	{
	return db.Where("username = ?", username).First(user).Error
}