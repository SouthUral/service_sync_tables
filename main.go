package main

import (
	"time"

	Api "github.com/SouthUral/service_sync_tables/api"
	Mongo "github.com/SouthUral/service_sync_tables/database/mongodb"
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

	inputMDBchan := make(Mongo.MongoInputChan, 100)
	outputMDBchan := make(Mongo.MongoOutputChan, 100)
	outputApiChan := make(Api.OutputAPIChan, 100)

	Mongo.MDBInit(inputMDBchan, outputMDBchan)
	Api.InitServer(outputApiChan)
	State.InitState(inputMDBchan, outputMDBchan, outputApiChan)

	log.Info("Server is starting")
	time.Sleep(1000 * time.Second)
}
