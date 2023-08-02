package main

import (
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
	http.HandleFunc("/all_sync", midlware(srv.allSync))
	http.HandleFunc("/add_sync", midlware(srv.addNewSync))
	http.HandleFunc("/stop_sync", midlware(srv.stopSync))
	http.HandleFunc("/start_sync", midlware(srv.startSync))
	http.ListenAndServe(srv.Port, nil)

}

// Обработчик для запроса на все синхронизации в сервисе
func (srv *Server) allSync(w http.ResponseWriter, r *http.Request) {
	GetMethod(w, r, GetAll, srv.OutputCh)
	log.Info("all_sync request processed")
}

// Обработчик для остановки синхронизации
func (srv *Server) stopSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StopSync, srv.OutputCh)
	log.Info("StopSync request processed")
}

// Обработчик для добавления синхронизации
func (srv *Server) addNewSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, InputData, srv.OutputCh)
	log.Info("AddNewSync request processed")
}

// Обработчик для старта остановленной синхронизации
func (srv *Server) startSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StartSync, srv.OutputCh)
	log.Info("AddNewSync request processed")
}

// мидлвар с дебаг логом
func midlware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(r.Method)
		handler(w, r)
	}
}
