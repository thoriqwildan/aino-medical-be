package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/model/converter"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB *gorm.DB
	Log *logrus.Logger
	UserRepository *repository.UserRepository
	Validate *validator.Validate
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, userRepository *repository.UserRepository, validate *validator.Validate) *UserUseCase {
	return &UserUseCase{
		DB: db,
		Log: log,
		UserRepository: userRepository,
		Validate: validate,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Validation error in CreateUser")
		return nil, err
	}

	
	if err := u.UserRepository.GetByUsername(tx, request.Username, &entity.User{}); err == nil {
		u.Log.WithField("name", request.Username).Error("Username already exists")
		return nil, fiber.ErrConflict
	}

	hash, err := helper.HashPassword(request.Password)
	if err != nil {
		u.Log.WithError(err).Error("Error hashing password in CreateUser")
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Username: request.Username,
		Password: hash,
		Name: request.Name,
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.WithError(err).Error("Error creating user in CreateUser")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Error committing transaction in CreateUser")
		return nil, err
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) LoginUser(ctx context.Context, request *model.LoginRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("Validation error in LoginUser")
		return nil, err
	}

	user := &entity.User{}
	if err := u.UserRepository.GetByUsername(tx, request.Username, user); err != nil {
		u.Log.WithError(err).Error("User not found in LoginUser")
		return nil, err
	}

	if err := helper.CheckPasswordHash(request.Password, user.Password); err != nil {
		u.Log.WithError(err).Error("Password mismatch in LoginUser")
		return nil, err
	}

	return converter.UserToResponse(user), nil
}

