package middleware

import (
	"errors"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type MiddlewareConfig struct {
	Viper *viper.Viper
	App   *fiber.App
	DB    *gorm.DB
}

func NewMiddlewareConfig(v *viper.Viper, app *fiber.App, db *gorm.DB) *MiddlewareConfig {
	return &MiddlewareConfig{
		Viper: v,
		App:   app,
		DB:    db,
	}
}

func (mc *MiddlewareConfig) JWTProtected() func(*fiber.Ctx) error {
	jwtwareConfig := jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(mc.Viper.GetString("JWT_SECRET"))},
		ContextKey:     "user",
		ErrorHandler:   mc.jwtError,
		SuccessHandler: mc.verifyTokenExpiration,
	}

	return jwtware.New(jwtwareConfig)
}

func (mc *MiddlewareConfig) verifyTokenExpiration(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	expires := int64(claims["exp"].(float64))
	var count int64
	if err := mc.DB.WithContext(c.Context()).Model(&entity.User{}).Where("username = ?", claims["username"]).Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error when find user "+err.Error())
	}
	if count == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if time.Now().Unix() > expires {
		return mc.jwtError(c, errors.New("Token has expired"))
	}
	return c.Next()
}

func (mc *MiddlewareConfig) jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Code:    fiber.StatusBadRequest,
			Message: "JWT not Found",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[any]{
		Code:    fiber.StatusUnauthorized,
		Message: "Unauthorized",
		Errors:  err.Error(),
	})
}
