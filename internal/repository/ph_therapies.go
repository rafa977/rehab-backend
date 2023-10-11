package repository

import (
	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// TherapyRepository --> Interface to TherapyRepository
type PhTherapyRepository interface {
	AddPhTherapy(models.PhTherapy) (models.PhTherapy, error)
	GetPhTherapy(int) (models.PhTherapy, error)
	GetPhTherapiesByCompanyID(int) ([]models.PhTherapy, error)
	GetAllTherapiesByDiseaseID(int) ([]models.PhTherapy, error)
	GetNumberOfTherapiesByDiseaseID(int) (int64, error)

	// DeletePhTherapy(int) (bool, error)
	// UpdatePhTherapy(models.Therapy) (models.Therapy, error)
}

type phTherapyService struct {
	dbConnection *gorm.DB
}

// NewTherapyService --> returns new therapy repository
func NewPhTherapyService() *phTherapyService {
	dbConnection := config.ConnectDB()

	return &phTherapyService{dbConnection: dbConnection}
}

func (db *phTherapyService) AddPhTherapy(phTherapy models.PhTherapy) (models.PhTherapy, error) {
	return phTherapy, db.dbConnection.Create(&phTherapy).Error
}

func (db *phTherapyService) GetPhTherapy(id int) (therapy models.PhTherapy, err error) {
	return therapy, db.dbConnection.Preload("AccountSuperVisor", func(tx *gorm.DB) *gorm.DB { return tx.Omit("Password") }).
		Preload("AccountEmployee", func(tx *gorm.DB) *gorm.DB { return tx.Omit("Password") }).
		Preload("Protocols").
		Preload("TherapyKeys").
		Preload("Exercises").First(&therapy, id).Error
}

func (db *phTherapyService) GetPhTherapiesByCompanyID(companyId int) (therapies []models.PhTherapy, err error) {
	return therapies, db.dbConnection.Where("company_id = ?", companyId).Find(&therapies).Error
}

func (db *phTherapyService) GetAllTherapiesByDiseaseID(diseaseID int) (therapies []models.PhTherapy, err error) {
	return therapies, db.dbConnection.Order("created_at desc").
		Preload("AccountSuperVisor", func(tx *gorm.DB) *gorm.DB { return tx.Omit("Password") }).
		Preload("AccountEmployee", func(tx *gorm.DB) *gorm.DB { return tx.Omit("Password") }).
		Where("disease_id = ?", diseaseID).Find(&therapies).Error
}

func (db *phTherapyService) GetNumberOfTherapiesByDiseaseID(diseaseID int) (count int64, err error) {
	return count, db.dbConnection.Model(&models.PhTherapy{}).Where("disease_id = ?", diseaseID).Count(&count).Error
}

// func (db *phTherapyService) UpdatePhTherapy(therapy models.Therapy) (models.Therapy, error) {
// 	if err := db.dbConnection.Preload("User").First(&therapy, therapy.ID).Error; err != nil {
// 		return therapy, err
// 	}
// 	return therapy, db.dbConnection.Preload("User").Model(&therapy).Updates(&therapy).Error
// }

// func (db *phTherapyService) DeletePhTherapy(id int) (bool, error) {
// 	return true, db.dbConnection.Delete(&models.Therapy{}, id).Error
// }
