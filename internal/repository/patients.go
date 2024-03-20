package repository

import (
	"strings"

	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// NewPatientRepository --> Interface to PatientRepository
type PatientRepository interface {
	AddPatient(models.Patient) (models.Patient, error)
	DeletePatient(uint) (bool, error)
	UpdatePatient(models.Patient) (models.Patient, error)
	GetPatient(uint) (models.Patient, error)

	GetPatientMedCard(int) (models.Patient, error)
	GetPatientByIdAndCompanyID(uint, uint) (models.Patient, error)
	GetPatientKeyword(string) ([]models.Patient, error)
	GetPatientAmka(int) ([]models.Patient, error)
	GetAllPatients([]uint) ([]models.Patient, error)

	GetAllPatientsByCompanyId(uint) ([]models.Patient, error)
	GetAllPatientsDetails([]uint) ([]models.Patient, error)

	CheckPatient(uint, []uint) (bool, string)
}

type patientService struct {
	dbConnection *gorm.DB
}

// NewPatientRepository --> returns new patient repository
func NewPatientService() *patientService {
	dbConnection := config.ConnectDB()

	return &patientService{dbConnection: dbConnection}
}

func (db *patientService) AddPatient(patient models.Patient) (models.Patient, error) {
	return patient, db.dbConnection.Create(&patient).Error
}

func (db *patientService) UpdatePatient(patient models.Patient) (models.Patient, error) {
	var oldPatient models.Patient
	if err := db.dbConnection.First(&oldPatient, patient.ID).Error; err != nil {
		return oldPatient, err
	}
	return patient, db.dbConnection.Model(&patient).Updates(&patient).Error
}

func (db *patientService) DeletePatient(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.Patient{}, id).Error
}

func (db *patientService) GetPatient(id uint) (patient models.Patient, err error) {
	return patient, db.dbConnection.Preload("GenericNotes").First(&patient, id).Error
}

func (db *patientService) GetPatientByIdAndCompanyID(patientId uint, companyId uint) (models.Patient, error) {
	var patient models.Patient
	return patient, db.dbConnection.Raw("select * from patient_details where id = ? and company_id = ?", patientId, companyId).Scan(&patient).Error
}

func (db *patientService) GetPatientMedCard(id int) (patient models.Patient, err error) {
	return patient, db.dbConnection.Preload("PatientDetails").First(&patient, id).Error
}

func (db *patientService) GetPatientKeyword(keyword string) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Where("firstname LIKE ?", keyword).Or("lastname LIKE ? ", keyword).Find(&patients).Error
}

func (db *patientService) GetPatientAmka(amka int) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Where("amka = ?", amka).Find(&patients).Error
}

func (db *patientService) GetAllPatients(companyID []uint) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Preload("Company").Omit("addedBy").Where("company_id IN ?", companyID).Find(&patients).Error
}

func (db *patientService) GetAllPatientsDetails(companyID []uint) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Preload("Company").Preload("PatientDetails").Find(&patients).Error
}

func (db *patientService) GetAllPatientsByCompanyId(id uint) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Preload("Company").Where("company_id = ?", id).Find(&patients).Error
}

// Check if patient editing belongs to company of caller
func (db *patientService) CheckPatient(id uint, compIDs []uint) (bool, string) {

	var patient models.Patient
	var err = db.dbConnection.First(&patient, id).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Patient does not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	var isOwner = false
	for _, id := range compIDs {

		if patient.CompanyID == id {
			isOwner = true
		}
	}

	if !isOwner {
		return false, "Patient does not belong to your company"
	}

	return true, ""
}
