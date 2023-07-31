package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	Port     string
	OutputCh StateAPIChan
}

// инициализатор go сервера
func InitServer(OutPutChan StateAPIChan) {
	srv := Server{
		Port:     ":" + getEnv("SERVER_PORT"),
		OutputCh: OutPutChan,
	}
	go srv.StartServer()
}

// go сервер
func (srv *Server) StartServer() {
	http.HandleFunc("/all_sync", midlware(srv.AllSync))
	http.ListenAndServe(srv.Port, nil)

}

// Обработчик для запроса на все синхронизации в сервисе
func (srv *Server) AllSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorWriter(w, "Request error", http.StatusBadRequest)
	}
	newChan := make(APImessChan)
	msg := APImessage{
		Message: GetAll,
		ApiChan: newChan,
	}
	srv.OutputCh <- msg
	log.Debug("Отправлено сообщение от API в STATE")
	answ, _ := <-newChan
	log.Debug("Получено сообщение от STATE")
	// Отправка сообщения клиенту
	JsonWriter(w, answ, http.StatusOK)
	log.Info("all_sync request processed")

}

// Обработчик для добавления синхронизации
func (srv *Server) AddNewSync(w http.ResponseWriter, r *http.Request) {
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
		Message: InputData,
		ApiChan: newChan,
		Data:    InpData,
	}
	srv.OutputCh <- msg
	log.Debug("Отправлено сообщение от API в STATE")
	answ, _ := <-newChan
	log.Debug("Получено сообщение от STATE")
	JsonWriter(w, answ, http.StatusOK)
	log.Info("all_sync request processed")

}

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

// мидлвар с дебаг логом
func midlware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(r.Method)
		handler(w, r)
	}
}
