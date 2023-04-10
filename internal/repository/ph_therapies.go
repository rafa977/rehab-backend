package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//TherapyRepository --> Interface to TherapyRepository
type PhTherapyRepository interface {
	AddDysfunction(models.Dysfunction) (models.Dysfunction, error)
	GetDysfunction(int) (models.Dysfunction, error)
	UpdateDysfunction(models.Dysfunction) (models.Dysfunction, error)
	DeleteDysfunction(int) (bool, error)

	AddPhTherapy(models.PhTherapy) (models.PhTherapy, error)
	GetPhTherapy(int) (models.PhTherapy, error)
	GetPhTherapiesByCompanyID(int) ([]models.PhTherapy, error)
	// DeletePhTherapy(int) (bool, error)
	// UpdatePhTherapy(models.Therapy) (models.Therapy, error)
}

type phTherapyService struct {
	dbConnection *gorm.DB
}

//NewTherapyService --> returns new therapy repository
func NewPhTherapyService() *phTherapyService {
	dbConnection := config.ConnectDB()

	return &phTherapyService{dbConnection: dbConnection}
}

///////////////////////////////////// Dysfunctions /////////////////////////////////////////////////////
func (db *phTherapyService) AddDysfunction(dysfnuction models.Dysfunction) (models.Dysfunction, error) {
	return dysfnuction, db.dbConnection.Create(&dysfnuction).Error
}

func (db *phTherapyService) GetDysfunction(id int) (dysfnuction models.Dysfunction, err error) {
	return dysfnuction, db.dbConnection.First(&dysfnuction, id).Error
}

func (db *phTherapyService) UpdateDysfunction(dysfnuction models.Dysfunction) (models.Dysfunction, error) {
	if err := db.dbConnection.First(&dysfnuction, dysfnuction.ID).Error; err != nil {
		return dysfnuction, err
	}
	return dysfnuction, db.dbConnection.Model(&dysfnuction).Updates(&dysfnuction).Error
}

func (db *phTherapyService) DeleteDysfunction(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.Dysfunction{}, id).Error
}

///////////////////////////////////// Dysfunctions /////////////////////////////////////////////////////

///////////////////////////////////// Physio Therapy /////////////////////////////////////////////////////
func (db *phTherapyService) AddPhTherapy(phTherapy models.PhTherapy) (models.PhTherapy, error) {
	return phTherapy, db.dbConnection.Create(&phTherapy).Error
}

func (db *phTherapyService) GetPhTherapy(id int) (therapy models.PhTherapy, err error) {
	return therapy, db.dbConnection.Preload("Patient").First(&therapy, id).Error
}

func (db *phTherapyService) GetPhTherapiesByCompanyID(companyId int) (therapies []models.PhTherapy, err error) {
	return therapies, db.dbConnection.Where("company_id = ?", companyId).Find(&therapies).Error
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
