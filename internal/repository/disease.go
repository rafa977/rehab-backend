package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//DiseaseRepository --> Interface to DiseaseRepository
type DiseaseRepository interface {
	GetDisease(int) (models.Disease, error)
	GetAllDiseasesPatientDetailsID(int) ([]models.Disease, error)
	AddDisease(models.Disease) (models.Disease, error)
	DeleteDisease(int) (bool, error)
	UpdateDisease(models.Disease) (models.Disease, error)
	// CheckPatientDetails(uint, uint) (bool, string)
	// CheckPatientDetailsOwning(int, []uint) (bool, string)
	// CheckDysfunctionCompany([]uint, int) (bool, string)
}

type diseaseService struct {
	dbConnection *gorm.DB
}

//NewDiseaseService --> returns new disease repository
func NewDiseaseService() *diseaseService {
	dbConnection := config.ConnectDB()

	return &diseaseService{dbConnection: dbConnection}
}

func (db *diseaseService) GetDisease(id int) (clinical models.Disease, err error) {
	return clinical, db.dbConnection.Preload("PatientDetails").First(&clinical, id).Error
}

func (db *diseaseService) GetAllDiseasesPatientDetailsID(patientDetailsId int) (dysfunctions []models.Disease, err error) {
	return dysfunctions, db.dbConnection.Where("patient_details_id = ?", patientDetailsId).Find(&dysfunctions).Error
}

func (db *diseaseService) AddDisease(therapy models.Disease) (models.Disease, error) {
	return therapy, db.dbConnection.Create(&therapy).Error
}

func (db *diseaseService) UpdateDisease(clinical models.Disease) (models.Disease, error) {
	if err := db.dbConnection.Preload("User").First(&clinical, clinical.ID).Error; err != nil {
		return clinical, err
	}
	return clinical, db.dbConnection.Preload("User").Model(&clinical).Updates(&clinical).Error
}

func (db *diseaseService) DeleteDisease(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.Disease{}, id).Error
}

// // function to check if the user is under the same company where the dysfunction/category is registered
// func (db *diseaseService) CheckDysfunctionCompany(compIDs []uint, dysfunctionID int) (bool, string) {

// 	// get dysfunction company ID
// 	var dysfunction models.Disease
// 	err := db.dbConnection.First(&dysfunction, dysfunctionID).Error
// 	if err != nil {
// 		var msg string
// 		if strings.Contains(err.Error(), "record not found") {
// 			msg = "Dysfunction does not exist"
// 		} else {
// 			msg = "Bad Request"
// 		}
// 		return false, msg
// 	}

// 	var isOwnerTest = false
// 	for _, id := range compIDs {
// 		if dysfunction.CompanyID == id {
// 			isOwnerTest = true
// 		}
// 	}

// 	if !isOwnerTest {
// 		return false, "Account does not belong to the same company"
// 	}

// 	return true, ""
// }

// func (db *diseaseService) CheckPatientDetails(id uint, compID uint) (bool, string) {

// 	var patient models.PatientDetails

// 	var err = db.dbConnection.Preload("Patient").First(&patient, id).Error
// 	if err != nil {
// 		var msg string
// 		if strings.Contains(err.Error(), "record not found") {
// 			msg = "You are not authorized to access these data."
// 		} else {
// 			msg = "Bad Request"
// 		}
// 		return false, msg
// 	}

// 	if patient.Patient.CompanyID != compID {
// 		return false, "Patient does not belong to your company"
// 	}

// 	return true, ""
// }

// func (db *diseaseService) CheckPatientDetailsOwning(id int, compIDs []uint) (bool, string) {

// 	var patient models.PatientDetails

// 	var err = db.dbConnection.Preload("Patient").First(&patient, id).Error
// 	if err != nil {
// 		var msg string
// 		if strings.Contains(err.Error(), "record not found") {
// 			msg = "You are not authorized to access these data"
// 		} else {
// 			msg = "Bad Request"
// 		}
// 		return false, msg
// 	}

// 	var isOwnerTest = false
// 	for _, id := range compIDs {
// 		if patient.Patient.CompanyID == id {
// 			isOwnerTest = true
// 		}
// 	}

// 	if !isOwnerTest {
// 		return false, "Patient does not belong to your company"
// 	}

// 	return true, ""
// }
