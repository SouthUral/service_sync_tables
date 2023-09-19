package main

import (
	"fmt"
	"time"

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
//	@description	This is a sample server Petstore server.

// @host		localhost:3000
// @BasePath	/
func main() {

	// Загрузка переменных окружения
	envVars := make(map[string]string)

	envVars["urlStoragePass"] = tools.GetEnv("URL_STORAGE_PASS")
	envVars["mongoHost"] = tools.GetEnv("MONGO_HOST")
	envVars["mongoPort"] = tools.GetEnv("MONGO_PORT")
	envVars["mongoCollection"] = tools.GetEnv("MONGO_COLLECTION")
	envVars["mongoDataBase"] = tools.GetEnv("MONGO_DATABASE")
	envVars["apiServerPort"] = tools.GetEnv("API_SERVER_PORT")
	envVars["logLevel"] = tools.GetEnv("LOG_LEVEL")

	for key, value := range envVars {
		if value == "" {
			log.Error(fmt.Sprintf("Переменная окружения %s не обнаружена!", key))
			return
		}
	}

	log.SetLevel(log.InfoLevel)

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
	URLStorage.InitUrlStorage(inputUrlChan, envVars["urlStoragePass"])

	log.Info("Server is starting")
	time.Sleep(100000 * time.Second)
}
