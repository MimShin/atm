package server

import (
	"atm/atm"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type AtmServer struct {
	listenAddr string
	atm        *atm.ATM
}

func NewAtmServer(listenAddr string, theATM *atm.ATM) *AtmServer {
	return &AtmServer{listenAddr: listenAddr, atm: theATM}
}

func (as *AtmServer) InitRouter() {

	r := mux.NewRouter()
	r.Use(as.logRequest)

	// Session APIs - don't need sessionID
	rSession := r.PathPrefix("/session/v1").Subrouter()
	rSession.HandleFunc("/login", as.login).Methods("POST")

	// User APIs - users operations
	rUser := r.PathPrefix("/api/v1").Subrouter()
	rUser.Use(as.checkAuth)
	rUser.HandleFunc("/accounts", as.listAccounts).Methods("GET")
	rUser.HandleFunc("/accounts/{accountID}", as.getAccount).Methods("GET")
	rUser.HandleFunc("/accounts/{accountID}/transactions", as.listTransactions).Methods("GET")
	rUser.HandleFunc("/accounts/{accountID}/transactions", as.postTransaction).Methods("POST")
	rUser.HandleFunc("/logout", as.logout).Methods("POST")
	// rUser.HandleFunc("/me", me).Methods("GET")

	// Admin APIs - for managing users and accounts
	rAdmin := r.PathPrefix("/admin/v1").Subrouter()
	rAdmin.Use(as.checkAuth)
	rAdmin.HandleFunc("/users", as.createUser).Methods("POST")
	rAdmin.HandleFunc("/users", as.listAllUsers).Methods("GET")
	rAdmin.HandleFunc("/accounts", as.createAccount).Methods("POST")
	rAdmin.HandleFunc("/accounts", as.listAllAccounts).Methods("GET")

	// Listen for REST requests
	log.Infof("starting server at %s", as.listenAddr)
	log.Fatal(http.ListenAndServe(as.listenAddr, r))
}
