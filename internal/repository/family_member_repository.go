package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
)

type FamilyMemberRepository struct {
	Repository[entity.FamilyMember]
	Log *logrus.Logger
}

func NewFamilyMemberRepository(log *logrus.Logger) *FamilyMemberRepository {
	return &FamilyMemberRepository{
		Log: log,
	}
}