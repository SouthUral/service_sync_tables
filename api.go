package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// инициализатор go сервера
func startApiServer() {
	go server()
}

// go сервер
func server() {
	http.HandleFunc("/all_sync", allSync)

	log.Info("Go server is starting")
	http.ListenAndServe(":1234", nil)

}

func allSync(w http.ResponseWriter, r *http.Request) {

}
