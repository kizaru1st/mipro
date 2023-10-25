package fakers

import (
	"time"

	"gorm.io/gorm"

	"github.com/bxcodec/faker/v4"

	"github.com/google/uuid"

	"github.com/kizaru1st/mipro/app/models"
)

func UserFaker(db *gorm.DB) *models.User {
	return &models.User{
		ID:            uuid.New().String(),
		FirstName:     faker.FirstName(),
		LastName:      faker.LastName(),
		Email:         faker.Email(),
		Password:      "5558967962808cbb969855c134aa0d3263b85003f9dafadb2bf62d1b1e8a12db",
		RememberToken: "",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     time.Time{},
	}
}
