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

// @Summary allSync
// @Tags Get
// @Description метод для получения всех синхронизаций
// Обработчик для запроса на все синхронизации в сервисе
func (srv *Server) allSync(w http.ResponseWriter, r *http.Request) {
	GetMethod(w, r, GetAll, srv.OutputCh)
	log.Info("all_sync request processed")
}

// @Summary stopSync
// @Tags Post
// @Description метод для остановки синхронизации
// @Param input body InputDataApi
// Обработчик для остановки синхронизации
func (srv *Server) stopSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StopSync, srv.OutputCh)
	log.Info("StopSync request processed")
}

// @Summary addNewSync
// @Tags Post
// @Description метод для добавления новой синхронизации
// @Param input body InputDataApi
// Обработчик для добавления синхронизации
func (srv *Server) addNewSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, InputData, srv.OutputCh)
	log.Info("AddNewSync request processed")
}

// @Summary startSync
// @Tags Post
// @Description метод для старта приостановленной синхронизации
// @Param input body InputDataApi
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
