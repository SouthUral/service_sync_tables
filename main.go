package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// запускает читает переменные окружения и запускает горутины
func main() {
	logrus.SetLevel(logrus.DebugLevel)

	chan_mongo_input := make(chan MessCommand, 100)
	chan_mongo_output := make(chan MessCommand, 100)

	go MongoMain(chan_mongo_input, chan_mongo_output)
	InitState(chan_mongo_input, chan_mongo_output)

	fmt.Println("Server is starting")
}
