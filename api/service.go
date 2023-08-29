package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"

	log "github.com/sirupsen/logrus"
)

// Метод для отправки сообщения клиенту в JSON
func JsonWriter[DataType StateAnswer | url.StorageConnDB](w http.ResponseWriter, data DataType, status int, err any) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		errString := fmt.Sprintf("%s", err)
		ErrorWriter(w, errString, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	errEncode := json.NewEncoder(w).Encode(data)
	if err != nil {
		ErrorWriter(w, errEncode.Error(), http.StatusInternalServerError)
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
	JsonWriter[StateAnswer](w, answ, http.StatusOK, answ.Err)
}

// GET метод для запросов к urlstorage
func GetURLmethod(w http.ResponseWriter, r *http.Request, OutputCh url.InputUrlStorageAPIch) {
	if r.Method != http.MethodGet {
		ErrorWriter(w, "Request error", http.StatusBadRequest)
		return
	}
	result, err := url.AllConn(OutputCh)
	if err != nil {
		ErrorWriter(w, "", http.StatusInternalServerError)
	}
	JsonWriter[url.StorageConnDB](w, result, http.StatusOK, err)
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
			JsonWriter(w, StateAnswer{Err: err.Error()}, http.StatusBadRequest, err.Error())
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
		JsonWriter(w, allAnswer, http.StatusOK, allAnswer.Err)
	case StopAll:
		allAnswer := StateAnswer{
			Info: "stop status of all sync",
		}
		newList := make([]StateAnswer, 0)
		for vall := range newChan {
			newList = append(newList, vall)
		}
		allAnswer.Data = newList
		JsonWriter(w, allAnswer, http.StatusOK, allAnswer.Err)
	default:
		answ, _ := <-newChan
		JsonWriter(w, answ, http.StatusOK, answ.Err)
	}

}
