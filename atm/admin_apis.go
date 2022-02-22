package atm

import (
	"atm/session"
)

func (atm *ATM) CreateUser(user User) (User, error) {
	createdUser := User{}

	user.Pin = session.Hash(user.ID + user.Pin)
	result := atm.db.Create(&user)
	if result.Error != nil {
		return createdUser, result.Error
	}

	result = atm.db.First(&createdUser, "id = ?", user.ID)
	return createdUser, result.Error
}

func (atm *ATM) ListAllUsers() (Users, error) {
	users := Users{}
	result := atm.db.Find(&users)
	return users, result.Error
}

func (atm *ATM) CreateAccount(account Account) (Account, error) {
	createdAcc := Account{}

	result := atm.db.Create(&account)
	if result.Error != nil {
		return createdAcc, result.Error
	}

	result = atm.db.First(&createdAcc, "id = ?", account.ID)
	return createdAcc, result.Error
}

func (atm *ATM) ListAllAccounts() (Accounts, error) {
	accounts := Accounts{}
	result := atm.db.Find(&accounts)
	return accounts, result.Error
}
