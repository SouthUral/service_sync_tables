package api

import (
	"fmt"
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
	urlSwag := fmt.Sprintf("http://localhost%s/swagger/doc.json", srv.Port)

	http.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL(urlSwag)))
	http.HandleFunc("/all_sync", middlwareGET(srv.allSync))
	http.HandleFunc("/add_sync", middlwarePOST(srv.addNewSync))
	http.HandleFunc("/stop_sync", middlwarePOST(srv.stopSync))
	http.HandleFunc("/start_sync", middlwarePOST(srv.startSync))
	http.HandleFunc("/start-allSync", middlwarePOST(srv.startAllSync))
	http.HandleFunc("/stop-allSync", middlwarePOST(srv.stopAllSync))
	http.HandleFunc("/all-connbd", middlwareGET(srv.getAllDBConn))
	http.HandleFunc("/one-conn-bd", middlwarePOST(srv.getOneDBConn))
	http.HandleFunc("/change-one-conn-bd", middlwarePUT(srv.changeOneDBConn))
	http.HandleFunc("/add-one-conn-bd", middlwarePOST(srv.addOneDBConn))
	http.ListenAndServe(srv.Port, nil)

}

// allSync godoc
//
//	@Summary 		allSync
//	@Description 	some description
//	@Tags 			GET
//	@Accept       	json
//	@Produce      	json
//	@Success		200		{object}	StateAnswer
//	@Router			/all_sync	[get]
func (srv *Server) allSync(w http.ResponseWriter, r *http.Request) {
	GetMethod(w, r, GetAll, srv.OutputCh)
	log.Info("all_sync request processed")
}

// @Summary 	stopSync
// @Tags 		POST
// @Description метод для остановки синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/stop_sync	[post]
func (srv *Server) stopSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StopSync, srv.OutputCh, true)
	log.Info("StopSync request processed")
}

// @Summary addNewSync
// @Tags POST
// @Description метод для добавления новой синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/add_sync	[post]
func (srv *Server) addNewSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, InputData, srv.OutputCh, true)
	log.Info("AddNewSync request processed")
}

// @Summary startSync
// @Tags POST
// @Description метод для старта приостановленной синхронизации
// @Param 		request 	body 	InputDataApi 	false 	"body example"
// @Router		/start_sync	[post]
func (srv *Server) startSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StartSync, srv.OutputCh, true)
	log.Info("StartSync request processed")
}

// @Summary startAllSync
// @Tags POST
// @Description метод для старта всех синхронизаций
// @Router		/start-allSync	[post]
func (srv *Server) startAllSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StartAll, srv.OutputCh, false)
	log.Info("StartAllSync request processed")
}

// @Summary stopAllSync
// @Tags POST
// @Description метод остановки всех синхронизаций
// @Router		/stop-allSync	[post]
func (srv *Server) stopAllSync(w http.ResponseWriter, r *http.Request) {
	PostMethod(w, r, StopAll, srv.OutputCh, false)
	log.Info("stopAllSync request processed")
}

// @Summary getAllDBConn
// @Tags GET
// @Accept       	json
// @Produce      	json
// @Success		200		{object}	url.StorageConnDB
// @Description метод для получения параметров подключения к базам данных и их элиасу
// @Router		/all-connbd	[get]
func (srv *Server) getAllDBConn(w http.ResponseWriter, r *http.Request) {
	GetURLmethod(w, r, srv.URLInputCh)
	log.Info("GetAllDBConn request processed")
}

// @Summary getOneDBConn
// @Tags GET
// @Accept       	json
// @Produce      	json
// @Success		200		{object}	url.ConnDBData
// @Failure     500     {object}    ErrorResponse
// @Param 		request 	body 	RequestDBConn 	false 	"body example"
// @Description метод для получения параметров подключения к БД по элиасу
// @Router		/one-conn-bd	[post]
func (srv *Server) getOneDBConn(w http.ResponseWriter, r *http.Request) {
	getOneURLMethod(w, r, srv.URLInputCh)
	log.Info("GetOneDBConn request processed")
}

// @Summary ChangeOneDBConn
// @Tags PUT
// @Accept       	json
// @Produce      	json
// @Success		200		{object}	StateAnswer
// @Param 		request 	body 	url.JsonFormat 	false 	"body example"
// @Description метод для изменения параметров подключения к БД по элиасу
// @Router		/change-one-conn-bd	[put]
func (srv *Server) changeOneDBConn(w http.ResponseWriter, r *http.Request) {
	changeOneURLMethod(w, r, srv.URLInputCh)
	log.Info("ChangeOneDBConn request processed")
}

// @Summary AddOneDBConn
// @Tags POST
// @Accept       	json
// @Produce      	json
// @Success		200		{object}	url.StorageConnDB
// @Param 		request 	body 	url.JsonFormat 	false 	"body example"
// @Description метод для добавления параметров подключения к БД
// @Router		/add-one-conn-bd	[post]
func (srv *Server) addOneDBConn(w http.ResponseWriter, r *http.Request) {
	addOneURLMethod(w, r, srv.URLInputCh)
	log.Info("AddOneDBConn request processed")
}
