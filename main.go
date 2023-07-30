package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// запускает читает переменные окружения и запускает горутины
func main() {
	log.SetLevel(log.DebugLevel)

	chan_mongo_input := make(chan MessCommand, 100)
	chan_mongo_output := make(chan MessCommand, 100)
	chan_api_state := make(StateAPIChan, 100)

	MDBInit(chan_mongo_input, chan_mongo_output)
	InitServer(chan_api_state)
	InitState(chan_mongo_input, chan_mongo_output, chan_api_state)

	log.Info("Server is starting")
	time.Sleep(150 * time.Second)
}
