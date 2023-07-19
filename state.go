package main

import (
	"log"
)

const (
	StatusActive   = 0
	StatusInactive = 1
	InputData      = "input_data"
	GetAll         = "get_all"
	UpdateData     = "update_data"
	DropCollection = "drop_collection"
	DropData       = "drop_data"
)

type State struct {
	table      string
	status     int
	mongoError interface{}
	mdbInput   chan MessCommand
	mdbOutput  chan MessCommand
}

// создает структуру State и запускает горутину StateWorker
func InitState(mongoChInput chan MessCommand, mongoChOutput chan MessCommand) {
	w_state := State{
		mdbInput:  mongoChInput,
		mdbOutput: mongoChOutput,
	}
	go w_state.StateWorker()
}

// основная функция для работы с состоянием, при запуске она стартует и ждет ответа от монго
// состояние должно принимать парметрами 6 каналов, которые создаются в main
// для обработки сообщений с каждым из модулей можно создать функции (3 функции)
func (state *State) StateWorker() {
	// метод для старта, запрашивает из монго все документы
	state.mdbInput <- MessCommand{Info: "get_all"}

	// запускается бесконечный цикл обработки сообщений
	for {
		select {
		// ловит сообщение из mdb
		case mess := <-state.mdbOutput:
			state.MongoWorker(mess)
		}
	}

}

// Обработчик сообщений приходящих от модуля MongoDB
func (state *State) MongoWorker(mess MessCommand) {
	switch mess.Info {
	case GetAll:
		state.mdbGetAll(mess)
	case InputData:
		state.mdbInputData(mess)
	}

}

// обработчик сообщений из монго, работает с сообщниями InputData
// запуск горутины произойдет только после записи о синхронизации в mongo
func (state *State) mdbInputData(mess MessCommand) {
	if mess.Error != nil {
		log.Println("Данные не добавлены в Mongo: ", mess.Error)
		// отправка сообщения в канал REST о неудачном запуске
		return
	}
	if mess.Data.IsActive {
		state.InitSyncT(mess.Data)
	}
}

// обработчик сообщений из монго, работает с сообщниями GetAll
func (state *State) mdbGetAll(mess MessCommand) {
	if mess.Error != nil {
		log.Println("Старт синхронизации не состоялся по причине: ", mess.Error)
		state.mongoError = mess.Error
		return
	}
	// если статус синхронизации True то данные передаются функции для запуска горутины
	if mess.Data.IsActive {
		state.InitSyncT(mess.Data)
	}
}

// функция для запуска горутин синхронизации
func (state *State) InitSyncT(data StateMess) {

}
