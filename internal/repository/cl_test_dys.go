package repository

import (
	"strings"

	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//ClTestDysRepository --> Interface to ClTestDysRepository
type ClTestDysRepository interface {
	GetClTestDys(int) (models.ClinicalTestDysfunction, error)
	GetClTestDysByDysID(int) ([]models.ClinicalTestDysfunction, error)
	DeleteClTestDys(int) (bool, error)
	AddClTestDys(models.ClinicalTestDysfunction) (models.ClinicalTestDysfunction, error)
	UpdateClTestDys(models.ClinicalTestDysfunction) (models.ClinicalTestDysfunction, error)
	CheckDysfunctionCompanyClinical([]uint, int) (bool, string)
	CheckCompany([]uint, models.ClinicalTestDysfunction) (bool, string)
}

type clTestDysService struct {
	dbConnection *gorm.DB
}

//NewClTestDysService --> returns new clinical test dysfunction repository
func NewClTestDysService() *clTestDysService {
	dbConnection := config.ConnectDB()

	return &clTestDysService{dbConnection: dbConnection}
}

func (db *clTestDysService) GetClTestDys(id int) (clinical models.ClinicalTestDysfunction, err error) {
	return clinical, db.dbConnection.Preload("Patient").First(&clinical, id).Error
}

func (db *clTestDysService) GetClTestDysByDysID(id int) (tests []models.ClinicalTestDysfunction, err error) {
	return tests, db.dbConnection.Where("dysfunction_id = ?", id).Find(&tests).Error
}

func (db *clTestDysService) AddClTestDys(therapy models.ClinicalTestDysfunction) (models.ClinicalTestDysfunction, error) {
	return therapy, db.dbConnection.Create(&therapy).Error
}

func (db *clTestDysService) UpdateClTestDys(clinical models.ClinicalTestDysfunction) (models.ClinicalTestDysfunction, error) {
	if err := db.dbConnection.Preload("User").First(&clinical, clinical.ID).Error; err != nil {
		return clinical, err
	}
	return clinical, db.dbConnection.Preload("User").Model(&clinical).Updates(&clinical).Error
}

func (db *clTestDysService) DeleteClTestDys(id int) (bool, error) {
	return true, db.dbConnection.Delete(&models.ClinicalTestDysfunction{}, id).Error
}

// function to check if the user is under the same company where the dysfunction/category is registered
func (db *clTestDysService) CheckDysfunctionCompanyClinical(compIDs []uint, dysfunctionID int) (bool, string) {

	// get dysfunction company ID
	var dysfunction models.Dysfunction
	err := db.dbConnection.First(&dysfunction, dysfunctionID).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Dysfunction does not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	var isOwnerTest = false
	for _, id := range compIDs {
		if dysfunction.CompanyID == id {
			isOwnerTest = true
		}
	}

	if !isOwnerTest {
		return false, "Account does not belong to the same company"
	}

	return true, ""
}

// function to check if the user is under the same company where the dysfunction/category is registered
func (db *clTestDysService) CheckCompany(compIDs []uint, clinical models.ClinicalTestDysfunction) (bool, string) {

	// get dysfunction company ID
	var dysfunction models.Dysfunction
	err := db.dbConnection.First(&dysfunction, clinical.DysfunctionID).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Dysfunction does not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	// get clinical category company ID
	var clinicalTest models.ClinicalTests
	err = db.dbConnection.First(&clinicalTest, clinical.ClinicalTestsID).Error
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "Clinical Test does not exist"
		} else {
			msg = "Bad Request"
		}
		return false, msg
	}

	// if from both models company ID is different
	if clinicalTest.CompanyID != dysfunction.CompanyID {
		var isOwnerTest = false
		for _, id := range compIDs {
			if clinicalTest.CompanyID == id {
				isOwnerTest = true
			}
		}

		if !isOwnerTest {
			return false, "Account is not under category company"
		}

		var isOwnerDysfunction = false
		for _, id := range compIDs {
			if dysfunction.CompanyID == id {
				isOwnerDysfunction = true
			}
		}

		if !isOwnerDysfunction {
			return false, "Account is not under dysfunction company"
		}

	} else {
		var isOwner = false
		for _, id := range compIDs {
			if dysfunction.CompanyID == id {
				isOwner = true
			}
		}
		if !isOwner {
			return false, "Account is not under same company"
		}
	}

	return true, ""
}
