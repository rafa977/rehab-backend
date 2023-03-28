package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//ProductRepository --> Interface to ProductRepository
type CompanyRepository interface {
	GetCompanyByID(int) ([]models.Relation, error)
	UpdateCompany(models.Company) (models.Company, error)
	RegisterCompany(models.Company) (models.Company, error)
	AddRelation(models.Relation) (models.Relation, error)
	GetRelationsByAccountID(uint) ([]models.Relation, error)
	// GetRelationsByCompanyID(int) ([]models.Relation, error)
}

type companyService struct {
	dbConnection *gorm.DB
}

//NewProductRepository --> returns new product repository
func NewCompanyService() *companyService {
	dbConnection := config.ConnectDB()

	return &companyService{dbConnection: dbConnection}
}

// func (db *companyService) GetCompanyByID(id int) (company models.Company, err error) {
// 	return company, db.dbConnection.Preload("Relations").First(&company, id).Error
// }

func (db *companyService) GetCompanyByID(id int) (relations []models.Relation, err error) {
	return relations, db.dbConnection.Preload("Account").Omit("password").Preload("Companies").Joins("JOIN relation_companies ON relation_companies.relation_id = relations.id").Where("relation_companies.company_id = ?", id).Find(&relations).Error
}

func (db *companyService) UpdateCompany(company models.Company) (models.Company, error) {
	var data models.Company
	if err := db.dbConnection.First(&data, company.ID).Error; err != nil {
		return data, err
	}
	return company, db.dbConnection.Model(&data).Updates(&company).Error
}

func (db *companyService) RegisterCompany(company models.Company) (models.Company, error) {
	return company, db.dbConnection.Create(&company).Error
}

func (db *companyService) AddRelation(relation models.Relation) (models.Relation, error) {
	return relation, db.dbConnection.Create(&relation).Error
}

func (db *companyService) GetRelationsByAccountID(id uint) (relation []models.Relation, err error) {
	return relation, db.dbConnection.Where("account_id = ?", id).Preload("Companies").Find(&relation).Error
}

// func (db *companyService) GetRelationsByCompanyID(id int) (relation []models.Relation, err error) {
// 	return relation, db.dbConnection.Where("account_id = ?", id).Preload("Companies").Find(&relation).Error
// }
