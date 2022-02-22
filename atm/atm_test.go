package atm

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	adminID  = "admin"
	adminPin = "0000"
	hashLen  = 64
)

var user1 User = User{
	ID:       "user-1",
	Name:     "ATM User",
	Pin:      "1111",
	IsActive: true,
}

var account1 Account = Account{
	ID:       "acc-1",
	Name:     "ATM Account",
	Currency: "CAD",
	Balance:  0,
	OwnerID:  user1.ID,
	IsActive: true,
}

func TestUserCreation(t *testing.T) {

	dbPath := fmt.Sprintf("/tmp/%d.db", os.Getpid())
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.Nil(t, err)
	require.NotNil(t, db)
	require.IsType(t, &gorm.DB{}, db)
	defer os.Remove(dbPath)

	newAtm := NewATM(db, adminPin)
	require.IsType(t, &ATM{}, newAtm)

	// Login with wrong pin -
	sessID, ok := newAtm.Login(adminID, "wrong pin")
	assert.False(t, ok)
	assert.Empty(t, sessID)

	// Admin login -
	sessID, ok = newAtm.Login(adminID, adminPin)
	assert.True(t, ok)
	assert.Len(t, sessID, hashLen)

	// Create user -
	createdUser, err := newAtm.CreateUser(user1)
	require.Nil(t, err)
	assert.True(t, createdUser.equals(user1))

	// Create Account for User -
	createdAcc, err := newAtm.CreateAccount(account1)
	require.Nil(t, err)
	assert.True(t, createdAcc.equals(account1))

	// User login -
	sessID, ok = newAtm.Login(user1.ID, user1.Pin)
	assert.True(t, ok)
	assert.Len(t, sessID, hashLen)

	// Deposit -
	tx := Transaction{
		Type:  "deposit",
		AccID: account1.ID,
		Value: 5000}

	updatedAcc, err := newAtm.deposit(user1.ID, tx)
	require.Nil(t, err)
	require.Equal(t, tx.Value+account1.Balance, updatedAcc.Balance)
	account1.Balance = updatedAcc.Balance

	// Withdraw with insufficient fund -
	tx = Transaction{
		Type:  "withdraw",
		AccID: account1.ID,
		Value: 5001}

	updatedAcc, err = newAtm.withdraw(user1.ID, tx)
	require.NotNil(t, err)
	require.Equal(t, account1.Balance, updatedAcc.Balance)

	// Withdraw -
	tx = Transaction{
		Type:  "withdraw",
		AccID: account1.ID,
		Value: 4000}

	updatedAcc, err = newAtm.withdraw(user1.ID, tx)
	require.Nil(t, err)
	require.Equal(t, account1.Balance-tx.Value, updatedAcc.Balance)

	// Logout -
	newAtm.Logout(user1.ID)

}

// helper functions
func (u1 User) equals(u2 User) bool {
	return u1.ID == u2.ID &&
		u1.Name == u2.Name &&
		u1.IsActive == u2.IsActive
}

func (a1 Account) equals(a2 Account) bool {
	return a1.ID == a2.ID &&
		a1.Name == a2.Name &&
		a1.Balance == a2.Balance &&
		a1.Currency == a2.Currency &&
		a1.OwnerID == a2.OwnerID &&
		a1.IsActive == a2.IsActive
}
