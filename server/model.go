package server

import "atm/atm"

// LoginRequest - body for login request
type LoginRequest struct {
	UserID atm.UserID `json:"user_id"`
	Pin    string     `json:"pin"`
}

// LoginResponse - response for successful login
type LoginResponse struct {
	XUserID    atm.UserID `json:"x-user-id"`
	XSessionID string     `json:"x-session-id"`
	Help       string     `json:"help"`
}

type TxRequest struct {
	Type  string    `json:"type"`
	Value atm.Money `json:"value"`
}
