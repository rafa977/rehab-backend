package repository

import (
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/gorm"
)

//ProductRepository --> Interface to ProductRepository
type AccountsRepository interface {
	GetAccountByID(int) (models.Account, error)
	GetAccountByUsernameForLogin(string) (models.Account, error)
	GetAccountByUsername(string) (models.Account, error)
	UpdateAccount(models.Account) (models.Account, error)
	AddUser(models.Account) (models.Account, error)
}

type accountService struct {
	dbConnection *gorm.DB
}

//NewProductRepository --> returns new product repository
func NewAccountService() *accountService {
	dbConnection := config.ConnectDB()

	return &accountService{dbConnection: dbConnection}
}

func (db *accountService) GetAccountByID(id int) (account models.Account, err error) {
	return account, db.dbConnection.Omit("password").First(&account, id).Error
}

func (db *accountService) GetCompanyByAccountByID(id int) (account models.Account, err error) {
	return account, db.dbConnection.Omit("password").First(&account, id).Error
}

func (db *accountService) GetAccountByUsernameForLogin(username string) (account models.Account, err error) {
	return account, db.dbConnection.Where("username = ?", username).First(&account).Error
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
