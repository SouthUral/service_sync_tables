package main

import (
	"time"

	Api "github.com/SouthUral/service_sync_tables/api"
	Config "github.com/SouthUral/service_sync_tables/config"
	Mongo "github.com/SouthUral/service_sync_tables/database/mongodb"
	Postgres "github.com/SouthUral/service_sync_tables/database/postgres"
	URLStorage "github.com/SouthUral/service_sync_tables/database/urlstorage"
	_ "github.com/SouthUral/service_sync_tables/docs"
	State "github.com/SouthUral/service_sync_tables/state"
	log "github.com/sirupsen/logrus"
)

//	@title			sync_service
//	@version		1.0
//	@description	This is a sample server Petstore server.

// @host		localhost:3000
// @BasePath	/
func main() {
	log.SetLevel(log.InfoLevel)

	// Загрузка конфигурации
	confStruct, err := Config.GetConf()
	if err != nil {
		log.Error("Конфигурация не загружена, программа завершена")
		return
	}

	outgoingSyncCh, ingoingPgCh := Postgres.GenerateChan()
	inputMDBchan := make(Mongo.MongoInputChan, 100)
	inputUrlChan := make(URLStorage.InputUrlStorageAPIch, 100)
	outputMDBchan := make(Mongo.MongoOutputChan, 100)
	outputApiChan := make(Api.OutputAPIChan, 100)

	Mongo.MDBInit(inputMDBchan, outputMDBchan)
	Api.InitServer(outputApiChan, inputUrlChan)
	State.InitState(inputMDBchan, outputMDBchan, outputApiChan, outgoingSyncCh, ingoingPgCh)
	Postgres.InitPostgres(inputUrlChan, ingoingPgCh, outgoingSyncCh)

	// инициализация модуля работы с URL
	URLStorage.InitUrlStorage(inputUrlChan, confStruct)

	log.Info("Server is starting")
	time.Sleep(1000 * time.Second)
}
