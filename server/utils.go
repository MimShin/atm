package server

import (
	"fmt"
	"net/http"
)

func writeError(w http.ResponseWriter, err error, httpStatus int) {
	msg := fmt.Sprintf("{%q: %q}", "error", err.Error())
	w.Header().Add("error", err.Error())
	w.WriteHeader(httpStatus)
	w.Write([]byte(msg))
}

func writeMessage(w http.ResponseWriter, msg string, httpStatus int) {
	w.WriteHeader(httpStatus)
	msgJson := fmt.Sprintf("{ %q: %q }", "message", msg)
	w.Write([]byte(msgJson))
}
