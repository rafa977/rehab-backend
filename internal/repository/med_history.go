package repository

import (
	"strings"

	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// ProductRepository --> Interface to ProductRepository
type MedHistoryRepository interface {
	AddMedHistory(models.MedHistory) (models.MedHistory, error)
	GetMedicalHistoryFull(uint) (models.MedHistory, error)
	GetMedicalHistorySpecific(uint, string) (models.MedHistory, error)
	UpdateMedicalHistory(models.MedHistory) (models.MedHistory, error)

	GetPatientDetails(int) (models.PatientDetails, error)
	GetPatientDetailsFull(uint) (models.PatientDetails, error)
	GetPatientDetailsByIdAndCompanyID(int, int) models.PatientDetails
	AddPatientDetails(models.PatientDetails) (models.PatientDetails, error)
	UpdatePatientDetails(models.PatientDetails) (models.PatientDetails, error)
	DeletePatientDetails(int) (bool, error)

	GetPatientDetailsByCompanyID(int) ([]models.PatientDetails, error)
	GetPatientDetailsByPatientID(uint) ([]models.PatientDetails, error)
	GetPatientDetailsForEmployeeID(uint, uint) ([]models.PatientDetailsPermission, error)

	CheckPatientByPatientDetailsID(uint, []uint) (bool, string)
	CheckPatientByDiseaseID(uint, []uint) (bool, string)

	AddMedHistoryPermission(models.MedHistoryPermission) (models.MedHistoryPermission, error)
	GetMedHistoryPermission(uint, uint) (models.MedHistoryPermission, error)
}

type medHistoryService struct {
	dbConnection *gorm.DB
}

// NewProductRepository --> returns new product repository
func NewMedHistoryService() *medHistoryService {
	dbConnection := config.ConnectDB()

	return &medHistoryService{dbConnection: dbConnection}
}

func (db *medHistoryService) AddMedHistory(history models.MedHistory) (models.MedHistory, error) {
	return history, db.dbConnection.Create(&history).Error
}

func (db *medHistoryService) GetMedicalHistoryFull(id uint) (history models.MedHistory, err error) {
	return history, db.dbConnection.Preload("PersonalAllergies").First(&history, id).Error
}

func (db *medHistoryService) GetMedicalHistorySpecific(id uint, historyType string) (history models.MedHistory, err error) {
	return history, db.dbConnection.Preload(historyType).First(&history, id).Error
}

func (db *medHistoryService) UpdateMedicalHistory(history models.MedHistory) (models.MedHistory, error) {
	var oldHistory models.MedHistory
	if err := db.dbConnection.First(&oldHistory, history.ID).Error; err != nil {
		return oldHistory, err
	}
	return history, db.dbConnection.Model(&history).Updates(&history).Error
}

func (db *medHistoryService) GetPatientDetails(id int) (patient models.PatientDetails, err error) {
	return patient, db.dbConnection.First(&patient, id).Error
}

func (db *medHistoryService) GetPatientDetailsByIdAndCompanyID(patientDetailsId int, companyId int) (patient models.PatientDetails) {
	db.dbConnection.Raw("select * from patient_details where id = ? and company_id = ?", patientDetailsId, companyId).Scan(&patient)
	return patient
}

func (db *medHistoryService) GetPatientDetailsByCompanyID(companyId int) (patientDetails []models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Preload("Patient").Where("company_id = ? ", companyId).Find(&patientDetails).Error
}

func (db *medHistoryService) GetPatientDetailsByPatientID(patientId uint) (patientDetails []models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Where("patient_id = ? ", patientId).Order("created_at desc").Find(&patientDetails).Error
}

func (db *medHistoryService) GetPatientDetailsFull(id uint) (patientDetails models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Preload("Diseases").First(&patientDetails, id).Error
}

func (db *medHistoryService) GetPatientDetailsForEmployeeID(patientID uint, account_id uint) (patientDetailsPermission []models.PatientDetailsPermission, err error) {
	return patientDetailsPermission, db.dbConnection.Preload("PatientDetails").
		Joins("JOIN patient_details ON patient_details_permissions.patient_details_id = patient_details.id").
		Where("patient_details_permissions.account_id = ? AND patient_details.patient_id = ?", account_id, patientID).
		Find(&patientDetailsPermission).Error
}

func (db *medHistoryService) DeletePatientDetails(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.PatientDetails{}, id).Error
}

func (db *medHistoryService) AddPatientDetails(patient models.PatientDetails) (models.PatientDetails, error) {
	return patient, db.dbConnection.Create(&patient).Error
}

func (db *medHistoryService) UpdatePatientDetails(patient models.PatientDetails) (models.PatientDetails, error) {
	var oldPatient models.PatientDetails
	if err := db.dbConnection.First(&oldPatient, patient.ID).Error; err != nil {
		return oldPatient, err
	}
	return patient, db.dbConnection.Model(&patient).Updates(&patient).Error
}

func (db *medHistoryService) CheckPatientByPatientDetailsID(id uint, compIDs []uint) (bool, string) {

	var patientDetails models.PatientDetails

	var err = db.dbConnection.Preload("Patient").First(&patientDetails, id).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Patient details do not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	var isOwner = false
	for _, id := range compIDs {

		if patientDetails.Patient.CompanyID == id {
			isOwner = true
		}
	}

	if !isOwner {
		return false, "Patient does not belong to your company"
	}

	return true, ""
}

func (db *medHistoryService) CheckPatientByDiseaseID(id uint, compIDs []uint) (bool, string) {

	var disease models.Disease

	var err = db.dbConnection.Preload("PatientDetails").Preload("PatientDetails.Patient").First(&disease, id).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Patient details do not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	var isOwner = false
	for _, id := range compIDs {

		if disease.PatientDetails.Patient.CompanyID == id {
			isOwner = true
		}
	}

	if !isOwner {
		return false, "Patient does not belong to your company"
	}

	return true, ""
}

///////////////////////// Patient Details Permissions ////////////////////////////////////////////////////////

func (db *medHistoryService) AddMedHistoryPermission(historyPermission models.MedHistoryPermission) (models.MedHistoryPermission, error) {
	return historyPermission, db.dbConnection.Create(&historyPermission).Error
}

func (db *medHistoryService) GetMedHistoryPermission(med_history_id uint, account_id uint) (medHistoryPermission models.MedHistoryPermission, err error) {
	return medHistoryPermission, db.dbConnection.Where("med_history_id = ? AND account_id = ?", med_history_id, account_id).Find(&medHistoryPermission).Error
}
