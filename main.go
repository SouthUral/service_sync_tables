package main

import (
	Api "github.com/SouthUral/service_sync_tables/api"
	Mongo "github.com/SouthUral/service_sync_tables/database/mongodb"
	Postgres "github.com/SouthUral/service_sync_tables/database/postgres"
	URLStorage "github.com/SouthUral/service_sync_tables/database/urlstorage"
	_ "github.com/SouthUral/service_sync_tables/docs"
	State "github.com/SouthUral/service_sync_tables/state"
	tools "github.com/SouthUral/service_sync_tables/tools"
	log "github.com/sirupsen/logrus"
)

//	@title			sync_service
//	@version		1.0
//	@description	API для взаимодействия с сервисом синхронизации таблиц баз данных.

// @host		localhost:3000
// @BasePath	/
func main() {

	// Загрузка переменных окружения
	envs, err := tools.LoadingEnvVars()
	if err != nil {
		log.Error(err)
		return
	}

	log.SetLevel(envs.LogLevel)

	outgoingSyncCh, ingoingPgCh := Postgres.GenerateChan()
	inputMDBchan := make(Mongo.MongoInputChan, 100)
	inputUrlChan := make(URLStorage.InputUrlStorageAPIch, 100)
	outputMDBchan := make(Mongo.MongoOutputChan, 100)
	outputApiChan := make(Api.OutputAPIChan, 100)

	Mongo.MDBInit(inputMDBchan, outputMDBchan, envs.GetMongoEnvs())
	Api.InitServer(outputApiChan, inputUrlChan, envs.ApiServerPort)
	State.InitState(inputMDBchan, outputMDBchan, outputApiChan, outgoingSyncCh, ingoingPgCh)
	Postgres.InitPostgres(inputUrlChan, ingoingPgCh, outgoingSyncCh)

	// инициализация модуля работы с URL
	URLStorage.InitUrlStorage(inputUrlChan, envs.UrlStoragePass)

	log.Info("Server is starting")
	ch := make(chan struct{})

	<-ch
	// time.Sleep(100000 * time.Second)
}
