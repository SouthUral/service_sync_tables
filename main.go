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

	MDBInit(chan_mongo_input, chan_mongo_output)
	InitState(chan_mongo_input, chan_mongo_output)

	log.Info("Server is starting")
	time.Sleep(150 * time.Second)
}
