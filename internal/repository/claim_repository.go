package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
)

type ClaimRepository struct {
	Repository[entity.Claim]
	Log *logrus.Logger
}

func NewClaimRepository(log *logrus.Logger) *ClaimRepository {
	return &ClaimRepository{
		Log: log,
	}
}