package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//TherapyRepository --> Interface to TherapyRepository
type TherapyRepository interface {
	GetTherapy(int) (models.Therapy, error)
	DeleteTherapy(int) (bool, error)
	AddTherapy(models.Therapy) (models.Therapy, error)
	UpdateTherapy(models.Therapy) (models.Therapy, error)
}

type therapyService struct {
	dbConnection *gorm.DB
}

//NewTherapyService --> returns new therapy repository
func NewTherapyService() *therapyService {
	dbConnection := config.ConnectDB()

	return &therapyService{dbConnection: dbConnection}
}

func (db *therapyService) GetTherapy(id int) (therapy models.Therapy, err error) {
	return therapy, db.dbConnection.Preload("Patient").First(&therapy, id).Error
}

func (db *therapyService) AddTherapy(therapy models.Therapy) (models.Therapy, error) {
	return therapy, db.dbConnection.Create(&therapy).Error
}

func (db *therapyService) UpdateTherapy(therapy models.Therapy) (models.Therapy, error) {
	if err := db.dbConnection.Preload("User").First(&therapy, therapy.ID).Error; err != nil {
		return therapy, err
	}
	return therapy, db.dbConnection.Preload("User").Model(&therapy).Updates(&therapy).Error
}

func (db *therapyService) DeleteTherapy(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.Therapy{}, id).Error
}
