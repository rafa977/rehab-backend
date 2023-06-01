package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//ProductRepository --> Interface to ProductRepository
type CompanyRepository interface {
	GetCompanyByID(int) (models.Company, error)
	UpdateCompany(models.Company) (models.Company, error)
	RegisterCompany(models.Company) (models.Company, error)
	AddRelation(models.Relation) (models.Relation, error)
	GetRelationsByAccountID(uint) ([]models.Relation, error)
	GetRelationIDsByAccountID(int) ([]models.Relation, error)
	GetCompaniesByAccountID(uint) []uint
	GetCompaniesDetailsByAccountID(uint) []models.Company
}

type companyService struct {
	dbConnection *gorm.DB
}

//NewProductRepository --> returns new product repository
func NewCompanyService() *companyService {
	dbConnection := config.ConnectDB()

	return &companyService{dbConnection: dbConnection}
}

func (db *companyService) GetCompanyByID(id int) (company models.Company, err error) {
	return company, db.dbConnection.First(&company, id).Error
}

// func (db *companyService) GetCompanyByID(id int) (relations []models.Relation, err error) {
// 	return relations, db.dbConnection.Preload("Account").Omit("password").Preload("Companies").Joins("JOIN relation_companies ON relation_companies.relation_id = relations.id").Where("relation_companies.company_id = ?", id).Find(&relations).Error
// }

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

func (db *companyService) GetRelationIDsByAccountID(id int) (relation []models.Relation, err error) {
	// var relat models.Relation
	return relation, db.dbConnection.Select("id").Where("account_id = ?", id).Preload("Companies").Find(&relation).Error
}

func (db *companyService) GetCompaniesByAccountID(id uint) (ids []uint) {
	db.dbConnection.Raw("select company_id as id from relation_companies rc where rc.relation_id in (select id as id from relations r where r.account_id = ?)", id).Scan(&ids)
	return ids
}

func (db *companyService) GetCompaniesDetailsByAccountID(id uint) (companies []models.Company) {
	db.dbConnection.Raw("select * from companies where id in (select company_id as id from relation_companies rc where rc.relation_id in (select id as id from relations r where r.account_id = ?))", id).Scan(&companies)
	return companies
}

// func (db *companyService) GetRelationsByCompanyID(id int) (relation []models.Relation, err error) {
// 	return relation, db.dbConnection.Where("account_id = ?", id).Preload("Companies").Find(&relation).Error
// }
