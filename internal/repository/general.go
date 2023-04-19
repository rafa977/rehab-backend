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

	// Disorder
	GetDisorder(int) (models.Disorder, error)
	GetAllDisorders() ([]models.Disorder, error)
	DeleteDisorder(uint) (bool, error)
	AddDisorder(models.Disorder) (models.Disorder, error)
	UpdateDisorder(models.Disorder) (models.Disorder, error)

	// Clinical Test Category
	GetClinicalTestCategory(int) (models.ClinicalTestCategory, error)
	GetAllClinicalTestCategories() ([]models.ClinicalTestCategory, error)
	DeleteClinicalTestCategory(uint) (bool, error)
	AddClinicalTestCategory(models.ClinicalTestCategory) (models.ClinicalTestCategory, error)
	UpdateClinicalTestCategory(models.ClinicalTestCategory) (models.ClinicalTestCategory, error)
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

// ############# Disorder CRUD ############################# //
func (db *generalService) GetDisorder(id int) (disorder models.Disorder, err error) {
	return disorder, db.dbConnection.First(&disorder, id).Error
}

func (db *generalService) GetAllDisorders() (allergies []models.Disorder, err error) {
	return allergies, db.dbConnection.Find(&allergies).Error
}

func (db *generalService) AddDisorder(disorder models.Disorder) (models.Disorder, error) {
	return disorder, db.dbConnection.Create(&disorder).Error
}

func (db *generalService) UpdateDisorder(disorder models.Disorder) (models.Disorder, error) {

	var oldDisorder models.Disorder

	if err := db.dbConnection.First(&oldDisorder, disorder.ID).Error; err != nil {
		return oldDisorder, err
	}
	return disorder, db.dbConnection.Model(&disorder).Updates(&disorder).Error
}

func (db *generalService) DeleteDisorder(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.Disorder{}, id).Error
}

// ############# Disorder CRUD END ############################# //

// ############# ClinicalTestCategory CRUD ############################# //
func (db *generalService) GetClinicalTestCategory(id int) (category models.ClinicalTestCategory, err error) {
	return category, db.dbConnection.First(&category, id).Error
}

func (db *generalService) GetAllClinicalTestCategories() (category []models.ClinicalTestCategory, err error) {
	return category, db.dbConnection.Find(&category).Error
}

func (db *generalService) AddClinicalTestCategory(category models.ClinicalTestCategory) (models.ClinicalTestCategory, error) {
	return category, db.dbConnection.Create(&category).Error
}

func (db *generalService) UpdateClinicalTestCategory(category models.ClinicalTestCategory) (models.ClinicalTestCategory, error) {

	var oldCategory models.ClinicalTestCategory

	if err := db.dbConnection.First(&oldCategory, category.ID).Error; err != nil {
		return oldCategory, err
	}
	return category, db.dbConnection.Model(&category).Updates(&category).Error
}

func (db *generalService) DeleteClinicalTestCategory(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.ClinicalTestCategory{}, id).Error
}

// ############# ClinicalTestCategory CRUD END ############################# //
