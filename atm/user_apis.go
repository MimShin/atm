package atm

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (atm *ATM) ListAccounts(userID UserID) (Accounts, error) {
	accounts := Accounts{}

	result := atm.db.Find(&accounts, "owner_id = ?", userID)
	return accounts, result.Error
}

func (atm *ATM) GetAccount(userID, accID string) (Account, error) {
	account := Account{}

	result := atm.db.First(&account, "owner_id = ? and id = ?", userID, accID)
	return account, result.Error
}

func (atm *ATM) ListTransactions(userID UserID, accID AccountID) (transactions Transactions, err error) {
	transactions = Transactions{}

	_, err = atm.GetAccount(userID, accID)
	if err != nil {
		log.Errorf("fetch account failed: err=%v", err)
		return
	}

	result := atm.db.Find(&transactions, "acc_id = ?", accID)
	return transactions, result.Error
}

func (atm *ATM) DoTransaction(userID UserID, tx Transaction) (account Account, err error) {
	switch tx.Type {
	case TxTypeDeposit:
		return atm.deposit(userID, tx)
	case TxTypeWithdraw:
		return atm.withdraw(userID, tx)
	default:
		err = fmt.Errorf("invalid transaction: type=%s", tx.Type)
		log.Error(err)
		return
	}

}

func (atm *ATM) withdraw(userID UserID, tx Transaction) (account Account, err error) {
	account = Account{}
	result := atm.db.First(&account, "owner_id = ? and id = ?", userID, tx.AccID)
	if result.Error != nil {
		err = result.Error
		log.Errorf("withdraw failed: err=%v", err)
		return
	}
	if account.Balance < tx.Value {
		err = fmt.Errorf("insufficient fund")
		log.Errorf("withdraw failed: err=%v", err)
		return
	}

	err = atm.db.Transaction(func(dbTx *gorm.DB) error {
		updates := map[string]interface{}{"balance": gorm.Expr("balance - ?", tx.Value), "updated_by": userID}
		result = dbTx.Model(&account).
			Where("owner_id = ? and id = ? and balance >= ?", userID, tx.AccID, tx.Value).Updates(updates)

		if result.Error != nil {
			err = result.Error
			log.Errorf("withdraw failed: err=%v", err)
			return err
		}

		tx.CreatedBy = userID
		result = dbTx.Create(&tx)
		if result.Error != nil {
			err = result.Error
			log.Errorf("updating transaction list failed: err=%v", err)
			return err
		}

		account.Balance -= tx.Value
		return nil
	})

	return
}

func (atm *ATM) deposit(userID UserID, tx Transaction) (account Account, err error) {
	account = Account{}
	result := atm.db.First(&account, "owner_id = ? and id = ?", userID, tx.AccID)
	if result.Error != nil {
		err = result.Error
		log.Errorf("deposit failed: err=%v", err)
		return
	}

	err = atm.db.Transaction(func(dbTx *gorm.DB) error {
		updates := map[string]interface{}{"balance": gorm.Expr("balance + ?", tx.Value), "updated_by": userID}
		result = dbTx.Model(&account).
			Where("owner_id = ? and id = ?", userID, tx.AccID).
			Updates(updates)

		if result.Error != nil {
			err = result.Error
			log.Errorf("deposit failed: err=%v", err)
			return err
		}

		tx.CreatedBy = userID
		result = dbTx.Create(&tx)
		if result.Error != nil {
			err = result.Error
			log.Errorf("updating transaction list failed: err=%v", err)
			return err
		}

		account.Balance += tx.Value
		return nil
	})

	return
}
