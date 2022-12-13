package util

import (
	"net/http"

	"github.com/gorilla/mux"
)

func CreateHealthCheck(r *mux.Router) {
	r.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("smartiko-test is healthy"))
	})
}
