package repository

import (
	"strings"

	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

// ProductRepository --> Interface to ProductRepository
type PatientDetailRepository interface {
	GetPatientDetails(int) (models.PatientDetails, error)
	GetPatientDetailsFull(int) (models.PatientDetails, error)
	GetPatientDetailsByIdAndCompanyID(int, int) models.PatientDetails
	AddPatientDetails(models.PatientDetails) (models.PatientDetails, error)
	UpdatePatientDetails(models.PatientDetails) (models.PatientDetails, error)
	DeletePatientDetails(int) (bool, error)

	GetPatientDetailsByCompanyID(int) ([]models.PatientDetails, error)
	GetPatientDetailsByPatientID(int) ([]models.PatientDetails, error)

	CheckPatientByPatientDetailsID(uint, []uint) (bool, string)
	CheckPatientByDiseaseID(uint, []uint) (bool, string)

	AddPatientDetailsPermission(models.PatientDetailsPermission) (models.PatientDetailsPermission, error)
	GetPatientDetailsPermission(int, uint) (models.PatientDetailsPermission, error)
}

type patientDetailsService struct {
	dbConnection *gorm.DB
}

// NewProductRepository --> returns new product repository
func NewPatientDetailsService() *patientService {
	dbConnection := config.ConnectDB()

	return &patientService{dbConnection: dbConnection}
}

func (db *patientService) GetPatientDetails(id int) (patient models.PatientDetails, err error) {
	return patient, db.dbConnection.First(&patient, id).Error
}

func (db *patientService) GetPatientDetailsByIdAndCompanyID(patientDetailsId int, companyId int) (patient models.PatientDetails) {
	db.dbConnection.Raw("select * from patient_details where id = ? and company_id = ?", patientDetailsId, companyId).Scan(&patient)
	return patient
}

func (db *patientService) GetPatientDetailsByCompanyID(companyId int) (patientDetails []models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Preload("Patient").Where("company_id = ? ", companyId).Find(&patientDetails).Error
}

func (db *patientService) GetPatientDetailsByPatientID(patientId int) (patientDetails []models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Where("patient_id = ? ", patientId).Find(&patientDetails).Error
}

func (db *patientService) GetPatientDetailsFull(id int) (patientDetails models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Preload("Diseases").First(&patientDetails, id).Error
}

// func (db *patientService) GetPatientFull(id int) (patient models.Patient, err error) {
// 	return patient, db.dbConnection.Preload("PersnoalAllergies").Preload("DrugTreatments").Preload("Therapies").Preload("MedicalTherapies").First(&patient, id).Error
// }

// func (db *patientService) GetPatientWithTherapies(id int) (patient models.Patient, err error) {
// 	return patient, db.dbConnection.Preload("Therapies").First(&patient, id).Error
// }

// func (db *patientService) GetAllPatients() (patients []models.Patient, err error) {
// 	return patients, db.dbConnection.Find(&patients).Error
// }

// func (db *findStorageRepository) GetCurrentusersProducts(id int) (products []model.Products, err error) {
// 	return products, db.connection.Preload("User").Where("user_id = ?", id).Find(&products).Error

// }

func (db *patientService) DeletePatientDetails(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.PatientDetails{}, id).Error
}

func (db *patientService) AddPatientDetails(patient models.PatientDetails) (models.PatientDetails, error) {
	return patient, db.dbConnection.Create(&patient).Error
}

func (db *patientService) UpdatePatientDetails(patient models.PatientDetails) (models.PatientDetails, error) {
	var oldPatient models.PatientDetails
	if err := db.dbConnection.First(&oldPatient, patient.ID).Error; err != nil {
		return oldPatient, err
	}
	return patient, db.dbConnection.Model(&patient).Updates(&patient).Error
}

func (db *patientService) CheckPatientByPatientDetailsID(id uint, compIDs []uint) (bool, string) {

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

func (db *patientService) CheckPatientByDiseaseID(id uint, compIDs []uint) (bool, string) {

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

func (db *patientService) AddPatientDetailsPermission(patientDetailsPermission models.PatientDetailsPermission) (models.PatientDetailsPermission, error) {
	return patientDetailsPermission, db.dbConnection.Create(&patientDetailsPermission).Error
}

func (db *patientService) GetPatientDetailsPermission(patient_details_id int, account_id uint) (patientDetailsPermission models.PatientDetailsPermission, err error) {
	return patientDetailsPermission, db.dbConnection.Where("patient_details_id = ? AND account_id = ?", patient_details_id, account_id).Find(&patientDetailsPermission).Error
}
