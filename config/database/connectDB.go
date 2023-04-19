package config

import (
	"fmt"
	"log"

	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {

	var err error

	// Read configuration
	cfg, err := Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&models.Account{}, &models.Patient{}, &models.Therapy{}, &models.MedicalTherapy{}, &models.DrugTreatment{}, &models.Company{}, &models.Relation{},
		&models.DrugTreatment{}, &models.Drug{}, &models.PersonalAllergy{}, &models.Allergy{}, &models.PersonalDisorder{}, &models.Disorder{}, &models.Visit{},
		&models.Injury{}, &models.Dysfunction{}, &models.Protocol{}, &models.PhTherapy{}, &models.PhTherapyNote{}, &models.PhTherapyKey{}, &models.PatientExercise{},
		&models.ClinicalTestCategory{}, &models.ClinicalTests{}, &models.ClinicalTestDysfunction{}, &models.DysfunctionHistory{}
	)

	return db
}
