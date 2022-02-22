package atm

import (
	"atm/session"

	log "github.com/sirupsen/logrus"
)

func (atm *ATM) Login(userID, pin string) (string, bool) {
	user := User{}
	result := atm.db.First(&user, "id = ?", userID)
	if result.Error != nil || user.Pin != session.Hash(userID+pin) {
		log.Errorf("login failed: err=%v", result.Error)
		return "", false
	}
	sessionID := atm.sessionMgr.LoadOrStore(userID, pin)
	return sessionID, sessionID != ""
}

func (atm *ATM) Logout(userID string) {
	atm.sessionMgr.Delete(userID)
}

func (atm *ATM) CheckAuth(userID, sessionID string) bool {
	return atm.sessionMgr.Check(userID, sessionID)
}
