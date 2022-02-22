package server

import (
	"atm/atm"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (as *AtmServer) listAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(hdrXUserID)
	accounts, err := as.atm.ListAccounts(userID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	accBytes, err := json.Marshal(accounts)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accBytes)
}

func (as *AtmServer) getAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(hdrXUserID)
	accID := mux.Vars(r)["accountID"]
	account, err := as.atm.GetAccount(userID, accID)

	as.showAccount(w, account, err)
}

func (as *AtmServer) listTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(hdrXUserID)
	accID := mux.Vars(r)["accountID"]
	transactions, err := as.atm.ListTransactions(userID, accID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	accBytes, err := json.Marshal(transactions)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accBytes)
}

func (as *AtmServer) postTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(hdrXUserID)
	accID := mux.Vars(r)["accountID"]

	tx := TxRequest{}
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	atmTx := atm.Transaction{AccID: accID, Type: tx.Type, Value: tx.Value}
	account, err := as.atm.DoTransaction(userID, atmTx)

	as.showAccount(w, account, err)
}

func (as *AtmServer) showAccount(w http.ResponseWriter, account atm.Account, err error) {

	if err == gorm.ErrRecordNotFound {
		writeError(w, err, http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	accBytes, err := json.Marshal(account)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accBytes)
}
