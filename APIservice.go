package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Метод для отправки сообщения клиенту в JSON
func JsonWriter(w http.ResponseWriter, data StateAnswer, status int) {
	w.Header().Set("Content-Type", "application/json")
	if data.Err != nil {
		errString := fmt.Sprintf("%s", data.Err)
		ErrorWriter(w, errString, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		ErrorWriter(w, err.Error(), http.StatusInternalServerError)
	}
}

// метод для отправки ошибки клиенту
func ErrorWriter(w http.ResponseWriter, err string, status int) {
	errData := ErrorResponse{
		Status: false,
		Error:  err,
	}
	log.Error("server error: ", err)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errData)
}

// абстрактный метод для GET запросов
func GetMethod(w http.ResponseWriter, r *http.Request, mess string, OutputCh StateAPIChan) {
	if r.Method != http.MethodGet {
		ErrorWriter(w, "Request error", http.StatusBadRequest)
	}
	newChan := make(APImessChan)
	msg := APImessage{
		Message: mess,
		ApiChan: newChan,
	}
	OutputCh <- msg
	answ, _ := <-newChan
	JsonWriter(w, answ, http.StatusOK)
}

// абстрактный метод для POST запросов
func PostMethod(w http.ResponseWriter, r *http.Request, mess string, OutputCh StateAPIChan) {
	if r.Method != http.MethodPost {
		ErrorWriter(w, "Request error", http.StatusBadRequest)
	}
	var InpData InputDataApi
	err := json.NewDecoder(r.Body).Decode(&InpData)
	if err != nil {
		JsonWriter(w, StateAnswer{Err: err.Error()}, http.StatusBadRequest)
	}
	newChan := make(APImessChan)
	msg := APImessage{
		Message: mess,
		ApiChan: newChan,
		Data:    InpData,
	}
	OutputCh <- msg
	answ, _ := <-newChan
	JsonWriter(w, answ, http.StatusOK)
}
