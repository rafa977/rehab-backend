package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//NewPatientRepository --> Interface to PatientRepository
type PatientRepository interface {
	GetPatient(int) (models.Patient, error)
	GetPatientKeyword(string) ([]models.Patient, error)
	GetPatientAmka(int) ([]models.Patient, error)
	GetAllPatients() ([]models.Patient, error)
	AddPatient(models.Patient) (models.Patient, error)
	UpdatePatient(models.Patient) (models.Patient, error)
}

type patientService struct {
	dbConnection *gorm.DB
}

//NewPatientRepository --> returns new patient repository
func NewPatientService() *patientService {
	dbConnection := config.ConnectDB()

	return &patientService{dbConnection: dbConnection}
}

func (db *patientService) GetPatient(id int) (patient models.Patient, err error) {
	return patient, db.dbConnection.First(&patient, id).Error
}

func (db *patientService) GetPatientKeyword(keyword string) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Where("firstname LIKE ?", keyword).Or("lastname LIKE ? ", keyword).Find(&patients).Error
}

func (db *patientService) GetPatientAmka(amka int) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Where("amka = ?", amka).Find(&patients).Error
}

func (db *patientService) GetAllPatients() (patients []models.Patient, err error) {
	return patients, db.dbConnection.Find(&patients).Error
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
