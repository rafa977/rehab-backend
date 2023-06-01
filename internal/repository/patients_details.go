package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//ProductRepository --> Interface to ProductRepository
type PatientDetailRepository interface {
	GetPatientDetails(int) (models.PatientDetails, error)
	GetPatientDetailsFull(int) (models.PatientDetails, error)
	GetPatientDetailsByIdAndCompanyID(int, int) models.PatientDetails
	AddPatientDetails(models.PatientDetails) (models.PatientDetails, error)
	UpdatePatientDetails(models.PatientDetails) (models.PatientDetails, error)

	GetPatientDetailsByCompanyID(int) ([]models.PatientDetails, error)
}

type patientDetailsService struct {
	dbConnection *gorm.DB
}

//NewProductRepository --> returns new product repository
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

func (db *patientService) GetPatientDetailsFull(id int) (patientDetails models.PatientDetails, err error) {
	return patientDetails, db.dbConnection.Preload("Injuries").Preload("PersonalAllergies").Preload("DrugTreatments").Preload("DrugTreatments.Drug").Preload("Therapies").Preload("MedicalTherapies").First(&patientDetails, id).Error
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

// func (db *findStorageRepository) DeleteProduct(product model.Products) (model.Products, error) {
// 	if err := db.connection.Preload("User").First(&product, product.ID).Error; err != nil {
// 		return product, err
// 	}
// 	return product, db.connection.Preload("User").Delete(&product).Error
// }
