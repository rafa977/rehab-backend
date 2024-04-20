package repository

import (
	config "rehab/config/database"
	"rehab/internal/pkg/models"

	"gorm.io/gorm"
)

// ProductRepository --> Interface to ProductRepository
type AccountsRepository interface {
	GetAccountByID(uint) (models.Account, error)
	GetAccountByIDWithPassword(uint) (models.Account, error)
	GetAccountByUsernameForLogin(string) (models.Account, error)
	GetAccountByUsername(string) (models.Account, error)
	GetAccountsByCompanyId(uint) ([]models.Account, error)

	UpdateAccount(models.Account) (models.Account, error)
	AddUser(models.Account) (models.Account, error)
	GetCompanyByAccountID(uint) models.Company
	GetCompaniesByAccountID(uint) []uint
}

type accountService struct {
	dbConnection *gorm.DB
}

// NewProductRepository --> returns new product repository
func NewAccountService() *accountService {
	dbConnection := config.ConnectDB()

	return &accountService{dbConnection: dbConnection}
}

func (db *accountService) GetAccountByID(id uint) (account models.Account, err error) {
	return account, db.dbConnection.Omit("password").First(&account, id).Error
}

func (db *accountService) GetCompanyByAccountID(id uint) (company models.Company) {
	db.dbConnection.Raw("select * from companies where id = (select company_id as id from relation_companies rc where rc.relation_id = (select id as id from relations r where r.account_id = ?))", id).Scan(&company)
	return company
}

func (db *accountService) GetCompaniesByAccountID(id uint) (ids []uint) {
	db.dbConnection.Raw("select company_id as id from relation_companies rc where rc.relation_id in (select id as id from relations r where r.account_id = ?)", id).Scan(&ids)
	return ids
}

func (db *accountService) GetCompanyByAccountByID(id int) (account models.Account, err error) {
	return account, db.dbConnection.Omit("password").First(&account, id).Error
}

func (db *accountService) GetAccountByUsernameForLogin(username string) (account models.Account, err error) {
	return account, db.dbConnection.Where("username = ?", username).First(&account).Error
}

func (db *accountService) GetAccountsByCompanyId(id uint) (accounts []models.Account, err error) {
	return accounts, db.dbConnection.
		Joins("JOIN relations ON accounts.id = relations.account_id").
		Joins("JOIN relation_companies ON relations.id = relation_companies.relation_id").
		Where("relation_companies.company_id = ?", id).
		Select("accounts.*"). // Selects only fields from the accounts table
		Find(&accounts).Error
}

func (db *accountService) GetAccountByIDWithPassword(id uint) (account models.Account, err error) {
	return account, db.dbConnection.First(&account, id).Error
}

func (db *accountService) GetAccountByUsername(username string) (account models.Account, err error) {
	return account, db.dbConnection.Omit("password").Where("username = ?", username).First(&account).Error
}

func (db *accountService) UpdateAccount(account models.Account) (models.Account, error) {
	var user models.Account
	if err := db.dbConnection.First(&user, account.ID).Error; err != nil {
		return user, err
	}
	return account, db.dbConnection.Model(&user).Updates(&account).Error
}

func (db *accountService) AddUser(account models.Account) (models.Account, error) {
	return account, db.dbConnection.Create(&account).Error
}
