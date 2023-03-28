package repository

import (
	"fmt"

	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//GeneralRepository --> Interface to GeneralRepository
type GeneralRepository interface {
	// Drug
	GetDrug(int) (models.Drug, error)
	GetAllDrugs() ([]models.Drug, error)
	DeleteDrug(uint) (bool, error)
	AddDrug(models.Drug) (models.Drug, error)
	UpdateDrug(models.Drug) (models.Drug, error)

	// Allergy
	GetAllergy(int) (models.Allergy, error)
	GetAllAllergies() ([]models.Allergy, error)
	DeleteAllergy(uint) (bool, error)
	AddAllergy(models.Allergy) (models.Allergy, error)
	UpdateAllergy(models.Allergy) (models.Allergy, error)
}

type generalService struct {
	dbConnection *gorm.DB
}

//NewGeneralRepoService --> returns new general repository
func NewGeneralRepoService() *generalService {
	dbConnection := config.ConnectDB()

	return &generalService{dbConnection: dbConnection}
}

// ############# Drug CRUD ############################# //
func (db *generalService) GetDrug(id int) (drug models.Drug, err error) {
	return drug, db.dbConnection.First(&drug, id).Error
}

func (db *generalService) GetAllDrugs() (drugs []models.Drug, err error) {
	return drugs, db.dbConnection.Find(&drugs).Error
}

func (db *generalService) AddDrug(drug models.Drug) (models.Drug, error) {
	return drug, db.dbConnection.Create(&drug).Error
}

func (db *generalService) UpdateDrug(drug models.Drug) (models.Drug, error) {
	fmt.Println(drug.DrugTitle)

	var oldDrug models.Drug
	if err := db.dbConnection.First(&oldDrug, drug.ID).Error; err != nil {
		return oldDrug, err
	}
	fmt.Println(drug.DrugTitle)
	return drug, db.dbConnection.Model(&drug).Updates(&drug).Error
}

func (db *generalService) DeleteDrug(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.Drug{}, id).Error
}

// ############# Drug CRUD END ############################# //

// ############# Allergy CRUD ############################# //
func (db *generalService) GetAllergy(id int) (allergy models.Allergy, err error) {
	return allergy, db.dbConnection.First(&allergy, id).Error
}

func (db *generalService) GetAllAllergies() (allergies []models.Allergy, err error) {
	return allergies, db.dbConnection.Find(&allergies).Error
}

func (db *generalService) AddAllergy(allergy models.Allergy) (models.Allergy, error) {
	return allergy, db.dbConnection.Create(&allergy).Error
}

func (db *generalService) UpdateAllergy(allergy models.Allergy) (models.Allergy, error) {

	var oldAllergy models.Allergy

	if err := db.dbConnection.First(&oldAllergy, allergy.ID).Error; err != nil {
		return oldAllergy, err
	}
	return allergy, db.dbConnection.Model(&allergy).Updates(&allergy).Error
}

func (db *generalService) DeleteAllergy(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.Allergy{}, id).Error
}

// ############# Allergy CRUD END ############################# //
