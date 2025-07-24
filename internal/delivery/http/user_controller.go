package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type UserController struct {
	UseCase *usecase.UserUseCase
	Log *logrus.Logger
	Config *viper.Viper
}

func NewUserController(useCase *usecase.UserUseCase, log *logrus.Logger, config *viper.Viper) *UserController {
	return &UserController{
		UseCase: useCase,
		Log: log,
		Config: config,
	}
}

func (uc *UserController) GetTest(c *fiber.Ctx) error {
	// Example handler for testing purposes
	uc.Log.Info("Test endpoint hit")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Test endpoint is working",
	})
}

func (uc *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterRequest)
	ctx.BodyParser(request)

	response, err := uc.UseCase.CreateUser(ctx.Context(), request) 
	if err != nil {
		uc.Log.WithError(err).Error("Error creating user")
		return err
	}

	token, err := uc.GenerateToken(response.Username)
	if err != nil {
		uc.Log.WithError(err).Error("Error generating token")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.UserResponse]{
		Code: fiber.StatusCreated,
		Message: "User created successfully",
		AccessToken: token,
		Data: response,
	})
}

func (uc *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)
	ctx.BodyParser(request)

	response, err := uc.UseCase.LoginUser(ctx.Context(), request)
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[model.UserResponse]{
				Code: fiber.StatusUnauthorized,
				Message: "Invalid username or password",
			})
		} else if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[model.UserResponse]{
				Code: fiber.StatusNotFound,
				Message: "User not found",
			})
		}
		uc.Log.WithError(err).Error("Error logging in user")
		return err
	}

	token, err := uc.GenerateToken(response.Username)
	if err != nil {
		uc.Log.WithError(err).Error("Error generating token")
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.UserResponse]{
		Code: fiber.StatusOK,
		Message: "Login successful",
		AccessToken: token,
		Data: response,
	})
}

func (uc *UserController) GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(30)).Unix()

	t, err := token.SignedString([]byte(uc.Config.GetString("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}