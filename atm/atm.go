package atm

import (
	"atm/session"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ATM struct {
	db         *gorm.DB
	sessionMgr *session.SessionManager
}

func NewATM(db *gorm.DB, adminPin string) *ATM {
	if err := MigrateModel(db, adminPin); err != nil {
		log.Errorf("db model migration failed: err=%v", err)
		return nil
	}

	sessionMgr := session.NewSessionManager()
	return &ATM{db: db, sessionMgr: sessionMgr}
}
