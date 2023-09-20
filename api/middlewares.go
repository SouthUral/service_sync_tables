package api

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// мидлвар с дебаг логом
func middlwareGET(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Error(fmt.Sprintf("Вызван метод %s, ожидается метод %s", r.Method, http.MethodGet))
			ErrorWriter(w, "Request error", http.StatusBadRequest)
			return
		}
		log.Debug(r.Method)
		handler(w, r)
	}
}

// мидлвар с дебаг логом
func middlwarePOST(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Error(fmt.Sprintf("Вызван метод %s, ожидается метод %s", r.Method, http.MethodPost))
			ErrorWriter(w, "Request error", http.StatusBadRequest)
			return
		}
		log.Debug(r.Method)
		handler(w, r)
	}
}

// мидлвар с дебаг логом
func middlwarePUT(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			log.Error(fmt.Sprintf("Вызван метод %s, ожидается метод %s", r.Method, http.MethodPut))
			ErrorWriter(w, "Request error", http.StatusBadRequest)
			return
		}
		log.Debug(r.Method)
		handler(w, r)
	}
}
