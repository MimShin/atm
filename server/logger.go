package server

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var skipHeaders = map[string]bool{
	"Accept": true, "Accept-Encoding": true, "Connection": true,
}

func (as *AtmServer) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Request: %s %s", r.Method, r.URL.String())

		if log.GetLevel() == log.TraceLevel {

			for name, value := range r.Header {
				if !skipHeaders[name] {
					log.Tracef("Head: %s: %s", name, value)
				}
			}
			body, _ := ioutil.ReadAll(r.Body)                // get the body
			r.Body = ioutil.NopCloser(bytes.NewReader(body)) // put the body back
			log.Tracef("Body: %s", body)
		}

		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
