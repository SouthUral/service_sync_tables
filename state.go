package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	StatusActive   = 0
	StatusInactive = 1
	InputData      = "input_data"
	GetAll         = "get_all"
	UpdateData     = "update_data"
	DropCollection = "drop_collection"
	DropData       = "drop_data"
	Stop           = "stop"
	Continue       = "continue"
)

type State struct {
	table        string
	status       int
	mongoError   interface{}
	mdbInput     chan MessCommand
	mdbOutput    chan MessCommand
	syncOutput   chan syncMessChan
	stateStorage map[string]StateSyncStorage
}

// создает структуру State и запускает горутину StateWorker
func InitState(mongoChInput chan MessCommand, mongoChOutput chan MessCommand) {

	w_state := State{
		mdbInput:     mongoChInput,
		mdbOutput:    mongoChOutput,
		stateStorage: make(map[string]StateSyncStorage),
	}
	go w_state.StateWorker()
}

// основная функция для работы с состоянием, при запуске она стартует и ждет ответа от монго
// состояние должно принимать парметрами 6 каналов, которые создаются в main
// для обработки сообщений с каждым из модулей можно создать функции (3 функции)
func (state *State) StateWorker() {
	// метод для старта, запрашивает из монго все документы
	state.syncOutput = make(chan syncMessChan, 100)
	state.mdbInput <- MessCommand{Info: GetAll}

	// запускается бесконечный цикл обработки сообщений
	for {
		select {
		// ловит сообщение из mdb
		case mess := <-state.mdbOutput:
			state.MongoWorker(mess)
		case mess := <-state.syncOutput:
			state.handlerSyncThreads(mess)
		}
	}
}

// обработчик для сообщений которые приходят из горутин
func (state *State) handlerSyncThreads(mess syncMessChan) {
	itemSync := state.stateStorage[mess.id]
	if mess.Error != nil {
		itemSync.err = mess.Error
		log.Error(mess.Error)
		return
	}
	itemSync.offset = mess.Offset
	itemSync.isSave = false
	state.mdbInput <- MessCommand{
		Info: UpdateData,
		Data: StateMess{
			oid:      itemSync.id,
			DataBase: itemSync.dataBase,
			Table:    itemSync.table,
			Offset:   fmt.Sprintf("%s", itemSync.offset),
			IsActive: itemSync.isActive,
		},
	}
	log.Debug("Данные из горутины отправлены на сохранение в MongoDB")
}

// Обработчик сообщений приходящих от модуля MongoDB
func (state *State) MongoWorker(mess MessCommand) {
	switch mess.Info {
	case GetAll:
		state.mdbGetAll(mess)
	case InputData:
		state.mdbInputData(mess)
	case UpdateData:
		state.mdbUpdateData(mess)
	}
}

// метод обработчик для сообщений UpdateData из модуля MongoDB
func (state *State) mdbUpdateData(mess MessCommand) {
	id := fmt.Sprintf("%s", mess.Data.oid)
	itemSync := state.stateStorage[id]
	if mess.Error != nil {
		log.Error("Данные не обновлены в Mongo: ", mess.Error)
		// добавляет данные об ошибке в хранилище
		itemSync.err = mess.Error
		itemSync.isSave = false
		itemSync.dateEnd = time.Now()
		// отправляет сообщение горутине об остановке
		itemSync.syncChan <- Stop
		return
	}
	// если данные обновлены то в горутину отпрвляется сообщение о продолжении работы
	itemSync.syncChan <- Continue
}

// обработчик сообщений из монго, работает с сообщниями InputData
// запуск горутины произойдет только после записи о синхронизации в mongo
func (state *State) mdbInputData(mess MessCommand) {
	if mess.Error != nil {
		log.Error("Данные не добавлены в Mongo: ", mess.Error)
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
		log.Error("Старт синхронизации не состоялся по причине: ", mess.Error)
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
	// создает канал для связи с горутиной, запускает горутину и записывает канал в структуру по id
	syncInput := make(chan string)
	go SyncTables(data, syncInput, state.syncOutput)
	// создает новую запись в словаре stateStorage
	// id := fmt.Sprintf("%s", data.oid)
	state.stateStorage[data.oid] = StateSyncStorage{
		id:        data.oid,
		table:     data.Table,
		dataBase:  data.DataBase,
		offset:    data.Offset,
		err:       nil,
		isSave:    true,
		isActive:  true,
		syncChan:  syncInput,
		dateStart: time.Now(),
		dateEnd:   nil,
	}
	log.Debug("Данные записаны в локальное хранение")
	log.Debug(fmt.Sprintf("%+v\n", state.stateStorage[data.oid]))
}

// функция запускается в отдельном потоке, ее задача подключиться к БД1 и БД2 и синхронизировать их
// останавливается горутина сообщением из канала inputChan
func SyncTables(data StateMess, inputChan chan string, outputChan chan syncMessChan) {
	//  здесь должно быть подключение к БД1 и БД2
	// в случае неудачного подключения нужно отправить ошибку в канал outputChan и завершить работу горутины
	var answer string

	intOffset, err := strconv.Atoi(data.Offset)
	newOffset := intOffset + 1
	if err != nil {
		log.Error(err)
	}

	for answer != Stop {

		test_answer := syncMessChan{
			Offset: strconv.Itoa(newOffset),
			id:     fmt.Sprintf("%s", data.oid),
			Error:  err,
		}
		outputChan <- test_answer
		log.Debug("Сообщение отправлено, offset: ", newOffset)
		answer = <-inputChan
		time.Sleep(1 * time.Millisecond)
		newOffset++
	}
	// close(outputChan)

	defer log.Debug("Получено сообщение, цикл прекращен: ", answer)
	return
}
