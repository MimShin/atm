package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (as *AtmServer) login(w http.ResponseWriter, r *http.Request) {
	loginReq := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil || loginReq.UserID == "" || loginReq.Pin == "" {
		err := fmt.Errorf("login failed")
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	sessionID, ok := as.atm.Login(loginReq.UserID, loginReq.Pin)
	if !ok {
		err := fmt.Errorf("login failed")
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	help := fmt.Sprintf("please add '%s' and '%s' headers to your requests", hdrXUserID, hdrXSessionID)
	loginResp := LoginResponse{XUserID: loginReq.UserID, XSessionID: sessionID, Help: help}

	respBytes, err := json.Marshal(loginResp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func (as *AtmServer) logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(hdrXUserID)
	as.atm.Logout(userID)
	writeMessage(w, "session terminated", http.StatusOK)
}
