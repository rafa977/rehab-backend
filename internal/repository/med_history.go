package repository

import (
	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// ProductRepository --> Interface to ProductRepository
type MedHistoryRepository interface {
	AddMedHistory(models.MedHistory) (models.MedHistory, error)
	GetMedicalHistoryFull(uint) (models.MedHistory, error)
	GetMedicalHistorySpecific(uint, string) (models.MedHistory, error)
	UpdateMedicalHistory(models.MedHistory) (models.MedHistory, error)

	AddMedHistoryPermission(models.MedHistoryPermission) (models.MedHistoryPermission, error)
	GetMedHistoryPermission(uint, uint) (models.MedHistoryPermission, error)

	// Surgeries
	DeleteSurgeryById(uint) (bool, error)

	// Injuries
	DeleteInjuryById(uint) (bool, error)
}

type medHistoryService struct {
	dbConnection *gorm.DB
}

// NewMedHistoryRepository --> returns new medical history repository
func NewMedHistoryService() *medHistoryService {
	dbConnection := config.ConnectDB()

	return &medHistoryService{dbConnection: dbConnection}
}

func (db *medHistoryService) AddMedHistory(history models.MedHistory) (models.MedHistory, error) {
	return history, db.dbConnection.Create(&history).Error
}

func (db *medHistoryService) GetMedicalHistoryFull(id uint) (history models.MedHistory, err error) {
	return history, db.dbConnection.Preload("Company").Preload("AddedBy").Preload("Therapies").Preload("MedicalTherapies").
		Preload("DrugTreatments").Preload("Injuries").Preload("PersonalAllergies").Preload("PersonalDisorders").
		Preload("Surgeries").Where("patient_id = ?", id).First(&history).Error
}

func (db *medHistoryService) GetMedicalHistorySpecific(id uint, historyType string) (history models.MedHistory, err error) {
	return history, db.dbConnection.Preload(historyType).First(&history, id).Error
}

func (db *medHistoryService) UpdateMedicalHistory(history models.MedHistory) (models.MedHistory, error) {
	var oldHistory models.MedHistory
	if err := db.dbConnection.First(&oldHistory, history.ID).Error; err != nil {
		return oldHistory, err
	}
	return history, db.dbConnection.Session(&gorm.Session{FullSaveAssociations: true}).Model(&history).Updates(&history).Error
}

// Surgeries
func (db *medHistoryService) DeleteSurgeryById(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.Surgery{}, id).Error
}

// Injuries
func (db *medHistoryService) DeleteInjuryById(id uint) (bool, error) {
	return true, db.dbConnection.Delete(&models.Injury{}, id).Error
}

///////////////////////// Medical History Permissions ////////////////////////////////////////////////////////

func (db *medHistoryService) AddMedHistoryPermission(historyPermission models.MedHistoryPermission) (models.MedHistoryPermission, error) {
	return historyPermission, db.dbConnection.Create(&historyPermission).Error
}

func (db *medHistoryService) GetMedHistoryPermission(med_history_id uint, account_id uint) (medHistoryPermission models.MedHistoryPermission, err error) {
	return medHistoryPermission, db.dbConnection.Where("med_history_id = ? AND account_id = ?", med_history_id, account_id).Find(&medHistoryPermission).Error
}
