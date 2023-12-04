package repository

import (
	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// HighlightRepository --> Interface to HighlightRepository
type HighlightRepository interface {
	AddHighlightDetails(models.Highlight) error
}

type highlightService struct {
	dbConnection *gorm.DB
}

// NewMedHistoryRepository --> returns new medical history repository
func NewHighlightService() *highlightService {
	dbConnection := config.ConnectDB()

	return &highlightService{dbConnection: dbConnection}
}

func (db *highlightService) AddHighlightDetails(details models.Highlight) (err error) {
	return db.dbConnection.Create(&details).Error
}
