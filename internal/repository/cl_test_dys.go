package repository

import (
	"strings"

	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// ClTestDysRepository --> Interface to ClTestDysRepository
type ClTestDisRepository interface {
	GetClTestDis(int) (models.ClinicalTestDisease, error)
	GetClTestDisByDisID(uint) ([]models.ClinicalTestDisease, error)
	DeleteClTestDis(int) (bool, error)
	AddClTestDis(models.ClinicalTestDisease) (models.ClinicalTestDisease, error)
	UpdateClTestDis(models.ClinicalTestDisease) (models.ClinicalTestDisease, error)
	CheckDiseaseCompanyClinical([]uint, uint) (bool, string)
	CheckCompany([]uint, models.ClinicalTestDisease) (bool, string)
}

type clTestDisService struct {
	dbConnection *gorm.DB
}

// NewClTestDysService --> returns new clinical test dysfunction repository
func NewClTestDisService() *clTestDisService {
	dbConnection := config.ConnectDB()

	return &clTestDisService{dbConnection: dbConnection}
}

func (db *clTestDisService) GetClTestDis(id int) (clinical models.ClinicalTestDisease, err error) {
	return clinical, db.dbConnection.Preload("Patient").First(&clinical, id).Error
}

func (db *clTestDisService) GetClTestDisByDisID(id uint) (tests []models.ClinicalTestDisease, err error) {
	return tests, db.dbConnection.Preload("ClinicalTests").Where("disease_id = ?", id).Find(&tests).Error
}

func (db *clTestDisService) AddClTestDis(test models.ClinicalTestDisease) (models.ClinicalTestDisease, error) {
	return test, db.dbConnection.Create(&test).Error
}

func (db *clTestDisService) UpdateClTestDis(clinical models.ClinicalTestDisease) (models.ClinicalTestDisease, error) {
	if err := db.dbConnection.Preload("User").First(&clinical, clinical.ID).Error; err != nil {
		return clinical, err
	}
	return clinical, db.dbConnection.Preload("User").Model(&clinical).Updates(&clinical).Error
}

func (db *clTestDisService) DeleteClTestDis(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.ClinicalTestDisease{}, id).Error
}

// function to check if the user is under the same company where the disease/category is registered
func (db *clTestDisService) CheckDiseaseCompanyClinical(compIDs []uint, diseaseID uint) (bool, string) {

	var diseaseCompanyID uint

	// Execute the raw SQL query
	query := "select company_id from patients p where p.id = (select patient_id from patient_details pd where pd.id = (select patient_details_id from diseases d where d.id = ?))"
	result := db.dbConnection.Raw(query, diseaseID).Scan(&diseaseCompanyID)
	if result.Error != nil {
		var msg string
		if strings.Contains(result.Error.Error(), "record not found") {
			msg = "Disease does not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	var isOwnerTest = false
	for _, id := range compIDs {
		if diseaseCompanyID == id {
			isOwnerTest = true
		}
	}

	if !isOwnerTest {
		return false, "Account does not belong to the same company"
	}

	return true, ""
}

// function to check if the user is under the same company where the dysfunction/category is registered
func (db *clTestDisService) CheckCompany(compIDs []uint, clinical models.ClinicalTestDisease) (bool, string) {

	// // get dysfunction company ID
	// var disease models.Disease
	// err := db.dbConnection.First(&disease, clinical.DiseaseID).Error
	// if err != nil {
	// 	var msg string
	// 	if strings.Contains(err.Error(), "record not found") {
	// 		msg = "Dysfunction does not exist"
	// 	} else {
	// 		msg = "Bad Request"
	// 	}
	// 	return false, msg
	// }

	// // get clinical category company ID
	// var clinicalTest models.ClinicalTests
	// err = db.dbConnection.First(&clinicalTest, clinical.ClinicalTestsID).Error
	// if err != nil {
	// 	var msg string
	// 	if strings.Contains(err.Error(), "record not found") {
	// 		msg = "Clinical Test does not exist"
	// 	} else {
	// 		msg = "Bad Request"
	// 	}
	// 	return false, msg
	// }

	// // if from both models company ID is different
	// if clinicalTest.CompanyID != dysfunction.CompanyID {
	// 	var isOwnerTest = false
	// 	for _, id := range compIDs {
	// 		if clinicalTest.CompanyID == id {
	// 			isOwnerTest = true
	// 		}
	// 	}

	// 	if !isOwnerTest {
	// 		return false, "Account is not under category company"
	// 	}

	// 	var isOwnerDysfunction = false
	// 	for _, id := range compIDs {
	// 		if dysfunction.CompanyID == id {
	// 			isOwnerDysfunction = true
	// 		}
	// 	}

	// 	if !isOwnerDysfunction {
	// 		return false, "Account is not under dysfunction company"
	// 	}

	// } else {
	// 	var isOwner = false
	// 	for _, id := range compIDs {
	// 		if dysfunction.CompanyID == id {
	// 			isOwner = true
	// 		}
	// 	}
	// 	if !isOwner {
	// 		return false, "Account is not under same company"
	// 	}
	// }

	return true, ""
}
