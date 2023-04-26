package repository

import (
	"strings"

	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//DysfunctionRepository --> Interface to DysfunctionRepository
type DysfunctionRepository interface {
	GetDysfunction(int) (models.Dysfunction, error)
	DeleteDysfunction(int) (bool, error)
	AddDysfunction(models.Dysfunction) (models.Dysfunction, error)
	UpdateDysfunction(models.Dysfunction) (models.Dysfunction, error)
	CheckPatientDetails(uint, uint) (bool, string)
}

type dysfunctionService struct {
	dbConnection *gorm.DB
}

//NewDysfunctionService --> returns new clinical test dysfunction repository
func NewDysfunctionService() *dysfunctionService {
	dbConnection := config.ConnectDB()

	return &dysfunctionService{dbConnection: dbConnection}
}

func (db *dysfunctionService) GetDysfunction(id int) (clinical models.Dysfunction, err error) {
	return clinical, db.dbConnection.Preload("PatientDetails").First(&clinical, id).Error
}

func (db *dysfunctionService) AddDysfunction(therapy models.Dysfunction) (models.Dysfunction, error) {
	return therapy, db.dbConnection.Create(&therapy).Error
}

func (db *dysfunctionService) UpdateDysfunction(clinical models.Dysfunction) (models.Dysfunction, error) {
	if err := db.dbConnection.Preload("User").First(&clinical, clinical.ID).Error; err != nil {
		return clinical, err
	}
	return clinical, db.dbConnection.Preload("User").Model(&clinical).Updates(&clinical).Error
}

func (db *dysfunctionService) DeleteDysfunction(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.Dysfunction{}, id).Error
}

func (db *dysfunctionService) CheckPatientDetails(id uint, compID uint) (bool, string) {

	var patient models.PatientDetails

	var err = db.dbConnection.Preload("Patient").First(&patient, id).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Patient Details does not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	if patient.Patient.CompanyID != compID {
		return false, "Patient does not belong to your company"
	}

	return true, ""
}
