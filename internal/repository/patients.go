package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//ProductRepository --> Interface to ProductRepository
type PatientRepository interface {
	GetPatient(int) (models.Patient, error)
	GetAllPatients() ([]models.Patient, error)
	AddPatient(models.Patient) (models.Patient, error)
}

type patientService struct {
	dbConnection *gorm.DB
}

//NewProductRepository --> returns new product repository
func NewPatientService() *patientService {
	dbConnection := config.ConnectDB()

	return &patientService{dbConnection: dbConnection}
}

func (db *patientService) GetPatient(id int) (patient models.Patient, err error) {
	return patient, db.dbConnection.Preload("Patient").First(&patient, id).Error
}

func (db *patientService) GetAllPatients() (patient []models.Patient, err error) {
	return patient, db.dbConnection.Find(&patient).Error
}

// func (db *findStorageRepository) GetCurrentusersProducts(id int) (products []model.Products, err error) {
// 	return products, db.connection.Preload("User").Where("user_id = ?", id).Find(&products).Error

// }

func (db *patientService) AddPatient(patient models.Patient) (models.Patient, error) {
	return patient, db.dbConnection.Create(&patient).Error
}

// func (db *findStorageRepository) UpdateProduct(product model.Products) (model.Products, error) {
// 	if err := db.connection.Preload("User").First(&product, product.ID).Error; err != nil {
// 		return product, err
// 	}
// 	return product, db.connection.Preload("User").Model(&product).Updates(&product).Error
// }

// func (db *findStorageRepository) DeleteProduct(product model.Products) (model.Products, error) {
// 	if err := db.connection.Preload("User").First(&product, product.ID).Error; err != nil {
// 		return product, err
// 	}
// 	return product, db.connection.Preload("User").Delete(&product).Error
// }
