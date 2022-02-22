package atm

import (
	"atm/session"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Money - money
type Money = int64

// AuditFields - common fields to keep userID and time for last change
type AuditFields struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedBy UserID    `json:"created_by"`
	UpdatedBy UserID    `json:"updated_by"`
	DeletedBy UserID    `json:"deleted_by"`
}

type UserID = string

// User - user info for account owners or transaction operators
type User struct {
	AuditFields `gorm:"embedded"`
	ID          string `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Pin         string `json:"pin"`
	IsActive    bool   `json:"is_active"`
}
type Users = []User

type AccountID = string
type Currency = string

// Account - a bank account
type Account struct {
	AuditFields `gorm:"embedded"`
	ID          AccountID `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Currency    Currency  `gorm:"not null" json:"currency"`
	Balance     Money     `gorm:"not null" json:"balance"`
	OwnerID     UserID    `gorm:"not null" json:"owner_id"`
	Owner       *User     `gorm:"foreignKey:OwnerID" json:"-"`
	IsActive    bool      `json:"is_active"`
}

type Accounts = []Account

type TxID = uint64
type TxType = string

// Transaction - account transactions
type Transaction struct {
	AuditFields `gorm:"embedded"`
	ID          TxID      `gorm:"primaryKey,autoIncrement" json:"id"`
	GroupID     TxID      `gorm:"not null" json:"-"`
	Type        TxType    `gorm:"not null" json:"type"`
	AccID       AccountID `gorm:"not null" json:"acc_id"`
	Acc         *Account  `gorm:"foreignKey:AccID" json:"-"`
	Value       Money     `gorm:"not null" json:"value"`
}

type Transactions = []Transaction

// Migrate (create or update) database schema
func MigrateModel(db *gorm.DB, adminPin string) error {
	if result := db.Exec("PRAGMA foreign_keys = ON", nil); result.Error != nil {
		log.Errorf("db referential integrity failed: err=%v", result.Error)
		return result.Error
	}

	if err := db.AutoMigrate(&User{}, &Account{}, &Transaction{}); err != nil {
		log.Errorf("db migration failed: err=%v", err)
		return err
	}

	// Add default admin
	admin := User{
		ID:       "admin",
		Name:     "Admin",
		Pin:      session.Hash("admin" + adminPin),
		IsActive: true,
	}
	result := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&admin)
	return result.Error
}
