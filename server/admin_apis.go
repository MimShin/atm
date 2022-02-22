package server

import (
	"atm/atm"
	"encoding/json"
	"net/http"
)

func (as *AtmServer) listAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := as.atm.ListAllUsers()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	usersBytes, err := json.Marshal(&users)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(usersBytes)
}

func (as *AtmServer) createUser(w http.ResponseWriter, r *http.Request) {
	user := atm.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	operator := r.Header.Get(hdrXUserID)
	user.CreatedBy = operator

	user, err = as.atm.CreateUser(user)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	userBytes, err := json.Marshal(&user)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(userBytes)
}

func (as *AtmServer) listAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := as.atm.ListAllAccounts()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	accountsBytes, err := json.Marshal(&accounts)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accountsBytes)
}

func (as *AtmServer) createAccount(w http.ResponseWriter, r *http.Request) {
	account := atm.Account{}
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	operator := r.Header.Get(hdrXUserID)
	account.CreatedBy = operator

	account, err = as.atm.CreateAccount(account)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	accountBytes, err := json.Marshal(&account)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(accountBytes)
}
