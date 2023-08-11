package api

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
func GetMethod(w http.ResponseWriter, r *http.Request, mess string, OutputCh OutputAPIChan) {
	if r.Method != http.MethodGet {
		ErrorWriter(w, "Request error", http.StatusBadRequest)
		return
	}
	newChan := make(InputAPIChan)
	msg := APImessage{
		Message: mess,
		ApiChan: newChan,
	}
	OutputCh <- msg
	answ, _ := <-newChan
	JsonWriter(w, answ, http.StatusOK)
}

// абстрактный метод для POST запросов
func PostMethod(w http.ResponseWriter, r *http.Request, mess string, OutputCh OutputAPIChan, body bool) {
	if r.Method != http.MethodPost {
		ErrorWriter(w, "Request error", http.StatusBadRequest)
		return
	}
	var InpData InputDataApi

	if body {
		err := json.NewDecoder(r.Body).Decode(&InpData)
		if err != nil {
			JsonWriter(w, StateAnswer{Err: err.Error()}, http.StatusBadRequest)
		}
	}

	newChan := make(InputAPIChan, 10)
	msg := APImessage{
		Message: mess,
		ApiChan: newChan,
		Data:    InpData,
	}
	OutputCh <- msg

	switch mess {
	case StartAll:
		allAnswer := StateAnswer{
			Info: "start status of all sync",
		}
		newList := make([]StateAnswer, 0)
		for vall := range newChan {
			newList = append(newList, vall)
		}
		allAnswer.Data = newList
		JsonWriter(w, allAnswer, http.StatusOK)
	case StopAll:
		allAnswer := StateAnswer{
			Info: "stop status of all sync",
		}
		newList := make([]StateAnswer, 0)
		for vall := range newChan {
			newList = append(newList, vall)
		}
		allAnswer.Data = newList
		JsonWriter(w, allAnswer, http.StatusOK)
	default:
		answ, _ := <-newChan
		JsonWriter(w, answ, http.StatusOK)
	}

}
