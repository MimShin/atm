package server

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	hdrXUserID    = "x-user-id"
	hdrXPin       = "x-pin"
	hdrXSessionID = "x-session-id"
)

func (as *AtmServer) checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("checkAuth")

		userID := r.Header.Get(hdrXUserID)
		if userID == "" {
			err := fmt.Errorf("mandatory header '%s' is missing", hdrXUserID)
			log.Error(err)
			writeError(w, err, http.StatusUnauthorized)
			return
		}

		log.Println(userID, r.URL.Path)
		if userID != "admin" && strings.HasPrefix(r.URL.Path, "/admin/") {
			err := fmt.Errorf("access to admin API denied")
			log.Error(err)
			writeError(w, err, http.StatusForbidden)
			return
		}

		sessionID := r.Header.Get(hdrXSessionID)
		if sessionID == "" {
			err := fmt.Errorf("mandatory header '%s' is missing", hdrXSessionID)
			log.Error(err)
			writeError(w, err, http.StatusUnauthorized)
			return
		}

		if !as.atm.CheckAuth(userID, sessionID) {
			err := fmt.Errorf("invalid session ID")
			log.Error(err)
			writeError(w, err, http.StatusUnauthorized)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
