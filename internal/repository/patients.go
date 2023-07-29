package repository

import (
	"fmt"
	"strings"

	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//NewPatientRepository --> Interface to PatientRepository
type PatientRepository interface {
	GetPatient(int) (models.Patient, error)
	GetPatientMedCard(int) (models.Patient, error)
	GetPatientByIdAndCompanyID(uint, uint) (models.Patient, error)
	GetPatientKeyword(string) ([]models.Patient, error)
	GetPatientAmka(int) ([]models.Patient, error)
	GetAllPatients() ([]models.Patient, error)
	AddPatient(models.Patient) (models.Patient, error)
	UpdatePatient(models.Patient) (models.Patient, error)
	CheckPatient(uint, []uint) (bool, string)
	GetAllPatientsByCompanyId(int) ([]models.Patient, error)
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

func (db *patientService) GetPatientMedCard(id int) (patient models.Patient, err error) {
	return patient, db.dbConnection.Preload("PatientDetails").First(&patient, id).Error
}

func (db *patientService) GetPatientKeyword(keyword string) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Where("firstname LIKE ?", keyword).Or("lastname LIKE ? ", keyword).Find(&patients).Error
}

func (db *patientService) GetPatientAmka(amka int) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Where("amka = ?", amka).Find(&patients).Error
}

func (db *patientService) GetAllPatients() (patients []models.Patient, err error) {
	return patients, db.dbConnection.Preload("Company").Find(&patients).Error
}

func (db *patientService) GetAllPatientsByCompanyId(id int) (patients []models.Patient, err error) {
	return patients, db.dbConnection.Preload("Company").Where("company_id = ?", id).Find(&patients).Error
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

func (db *patientService) GetPatientByIdAndCompanyID(patientId uint, companyId uint) (models.Patient, error) {
	var patient models.Patient
	return patient, db.dbConnection.Raw("select * from patient_details where id = ? and company_id = ?", patientId, companyId).Scan(&patient).Error
}

func (db *patientService) CheckPatient(id uint, compIDs []uint) (bool, string) {

	var patient models.Patient
	fmt.Println(id)
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
	fmt.Println(patient.ID)

	var isOwner = false
	for _, id := range compIDs {
		fmt.Println("we are here: ", id)
		fmt.Println("we are here patient: ", patient.CompanyID)
		fmt.Println("we are here patient: ", patient.Firstname)

		if patient.CompanyID == id {
			isOwner = true
		}
	}

	if !isOwner {
		return false, "Patient does not belong to your company"
	}

	return true, ""
}
