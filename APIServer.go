package main

import (
	"net/http"

	_ "github.com/SouthUral/service_sync_tables/docs"

	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

func (srv *Server) StartServer() {
	http.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:3000/swagger/doc.json")))
	http.HandleFunc("/all_sync", midlware(srv.allSync))
	http.HandleFunc("/add_sync", midlware(srv.addNewSync))
	http.HandleFunc("/stop_sync", midlware(srv.stopSync))
	http.HandleFunc("/start_sync", midlware(srv.startSync))
	http.ListenAndServe(srv.Port, nil)

}

// allSync godoc
//
//	@Summary 		allSync
//	@Description 	some description
//	@Tags 			Get
//	@Accept       	json
//	@Produce      	json
//	@Success		200		{object}	StateStorage
//	@Router			/all_sync	[get]
func (srv *Server) allSync(w http.ResponseWriter, r *http.Request) {
	GetMethod(w, r, GetAll, srv.OutputCh)
	log.Info("all_sync request processed")
}

// @Summary 	stopSync
// @Tags 		Post
// @Description метод для остановки синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/stop_sync	[post]
func (srv *Server) stopSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StopSync, srv.OutputCh)
	log.Info("StopSync request processed")
}

// @Summary addNewSync
// @Tags Post
// @Description метод для добавления новой синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/add_sync	[post]
func (srv *Server) addNewSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, InputData, srv.OutputCh)
	log.Info("AddNewSync request processed")
}

// @Summary startSync
// @Tags Post
// @Description метод для старта приостановленной синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/start_sync	[post]
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
