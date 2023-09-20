package api

import (
	"net/http"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"
	_ "github.com/SouthUral/service_sync_tables/docs"

	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	Port       string
	OutputCh   OutputAPIChan
	URLInputCh url.InputUrlStorageAPIch
}

// инициализатор go сервера
func InitServer(OutPutChan OutputAPIChan, URLChan url.InputUrlStorageAPIch, serverPort string) {
	srv := Server{
		Port:       ":" + serverPort,
		OutputCh:   OutPutChan,
		URLInputCh: URLChan,
	}
	go srv.StartServer()
}

func (srv *Server) StartServer() {
	http.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:3000/swagger/doc.json")))
	http.HandleFunc("/all_sync", middlwareGET(srv.allSync))
	http.HandleFunc("/add_sync", middlwarePOST(srv.addNewSync))
	http.HandleFunc("/stop_sync", middlwarePOST(srv.stopSync))
	http.HandleFunc("/start_sync", middlwarePOST(srv.startSync))
	http.HandleFunc("/start-allSync", middlwarePOST(srv.startAllSync))
	http.HandleFunc("/stop-allSync", middlwarePOST(srv.stopAllSync))
	http.HandleFunc("/all-conn-bd", middlwareGET(srv.GetAllDBConn))
	http.HandleFunc("/one-conn-bd", middlwareGET(srv.GetOneDBConn))
	http.HandleFunc("/change-one-conn-bd", middlwarePUT(srv.ChangeOneDBConn))
	http.HandleFunc("/add-one-conn-bd", middlwarePOST(srv.AddOneDBConn))
	http.ListenAndServe(srv.Port, nil)

}

// allSync godoc
//
//	@Summary 		allSync
//	@Description 	some description
//	@Tags 			Get
//	@Accept       	json
//	@Produce      	json
//	@Success		200		{object}	StateAnswer
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
	PostMethod(w, r, StopSync, srv.OutputCh, true)
	log.Info("StopSync request processed")
}

// @Summary addNewSync
// @Tags Post
// @Description метод для добавления новой синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/add_sync	[post]
func (srv *Server) addNewSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, InputData, srv.OutputCh, true)
	log.Info("AddNewSync request processed")
}

// @Summary startSync
// @Tags Post
// @Description метод для старта приостановленной синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/start_sync	[post]
func (srv *Server) startSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StartSync, srv.OutputCh, true)
	log.Info("StartSync request processed")
}

// @Summary startAllSync
// @Tags Post
// @Description метод для старта всех синхронизаций
// @Router		/start-allSync	[post]
func (srv *Server) startAllSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StartAll, srv.OutputCh, false)
	log.Info("StartAllSync request processed")
}

// @Summary stopAllSync
// @Tags Post
// @Description метод остановки всех синхронизаций
// @Router		/stop-allSync	[post]
func (srv *Server) stopAllSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StopAll, srv.OutputCh, false)
	log.Info("stopAllSync request processed")
}

func (srv *Server) GetAllDBConn(w http.ResponseWriter, r *http.Request) {
	GetURLmethod(w, r, srv.URLInputCh)
	log.Info("GetAllDBConn request processed")
}

func (srv *Server) GetOneDBConn(w http.ResponseWriter, r *http.Request) {
	getOneURLMethod(w, r, srv.URLInputCh)
	log.Info("GetOneDBConn request processed")
}

func (srv *Server) ChangeOneDBConn(w http.ResponseWriter, r *http.Request) {
	changeOneURLMethod(w, r, srv.URLInputCh)
	log.Info("ChangeOneDBConn request processed")
}

func (srv *Server) AddOneDBConn(w http.ResponseWriter, r *http.Request) {
	addOneURLMethod(w, r, srv.URLInputCh)
	log.Info("AddOneDBConn request processed")
}
