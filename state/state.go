package state

import (
	"fmt"
	"time"

	api "github.com/SouthUral/service_sync_tables/api"
	mongo "github.com/SouthUral/service_sync_tables/database/mongodb"

	log "github.com/sirupsen/logrus"
)

type State struct {
	table         string
	status        int
	mongoError    interface{}
	StorageChanI  StorageChanInput
	mdbInput      mongo.MongoInputChan
	mdbOutput     mongo.MongoOutputChan
	syncOutput    chan syncMessChan
	outputApiChan api.OutputAPIChan
	// если нужно отправить StateStorage в целиком в другую горутину
	// то нужно отправлять глубокую копию StateStorage иначе будет состояние гонки
	stateStorage StateStorage
}

// создает структуру State и запускает горутину StateWorker
func InitState(mdbInput mongo.MongoInputChan, mdbOutput mongo.MongoOutputChan, outputApiChan api.OutputAPIChan) {

	w_state := State{
		StorageChanI:  make(StorageChanInput),
		mdbInput:      mdbInput,
		mdbOutput:     mdbOutput,
		outputApiChan: outputApiChan,
		stateStorage:  make(StateStorage),
	}
	go w_state.StateWorker()
}

// основная функция для работы с состоянием, при запуске она стартует и ждет ответа от монго
// состояние должно принимать парметрами 6 каналов, которые создаются в main
// для обработки сообщений с каждым из модулей можно создать функции (3 функции)
func (state *State) StateWorker() {
	// метод для старта, запрашивает из монго все документы
	state.syncOutput = make(chan syncMessChan, 100)
	state.mdbInput <- mongo.MessCommand{Info: mongo.GetAll}

	// запускается бесконечный цикл обработки сообщений
	for {
		select {
		// ловит сообщение из mdb
		case mess := <-state.mdbOutput:
			state.MongoWorker(mess)
		case mess := <-state.syncOutput:
			state.HandlerSyncThreads(mess)
		case mess := <-state.outputApiChan:
			state.ApiHandler(mess)
		}
	}
}

// обработчик для сообщений которые приходят из горутин
func (state *State) HandlerSyncThreads(mess syncMessChan) {
	itemSync := state.stateStorage[mess.id]

	// отправляет сообщение в API либо об успешном страте либо об ошибке
	switch mess.Info {
	case StartSync:
		if mess.Error == nil {
			state.ResponseAPIRequest(mess.id, nil, StartSync)
		} else {
			state.StopSyncState(mess.id, mess.Error, false)
			state.ResponseAPIRequest(mess.id, mess.Error, StartSync)
		}
		return
	case RegularSync:
		itemSync.Offset = mess.Offset
		itemSync.IsSave = false
		state.stateStorage[mess.id] = itemSync
		state.updateDataMongo(mess.id)
		log.Debug("Данные из горутины отправлены на сохранение в MongoDB")
	case StopSync:
		state.ResponseAPIRequest(mess.id, nil, StopSync)
	}

}

// Метод для остановки синхронизации
// param: activeChan
func (state *State) StopSyncState(key string, err interface{}, activeChan bool) {
	itemSync := state.stateStorage[key]
	if activeChan {
		itemSync.syncChan <- Stop
	}
	itemSync.IsActive = false
	itemSync.syncChan = nil
	itemSync.DateEnd = time.Now()

	if err != nil {
		itemSync.Err = err
	} else {
		itemSync.IsSave = true
	}

	state.stateStorage[key] = itemSync
}

// создает новую запись в словаре stateStorage
func (state *State) AddInfoToStorage(data mongo.StateMess) {
	key := fmt.Sprintf("%s_%s", data.DataBase, data.Table)
	state.stateStorage[key] = StateSyncStorage{
		Id:        data.Oid,
		Table:     data.Table,
		DataBase:  data.DataBase,
		Offset:    data.Offset,
		Err:       nil,
		IsSave:    true,
		IsActive:  data.IsActive,
		DateStart: time.Now(),
		DateEnd:   nil,
	}
	log.Debug("Данные записаны в локальное хранение")
	log.Debug(fmt.Sprintf("%+v\n", state.stateStorage[data.Oid]))
}

// функция для запуска горутин синхронизации
func (state *State) InitSyncT(data mongo.StateMess) {
	// создает канал для связи с горутиной, запускает горутину и записывает канал в структуру по id
	StorageKey := fmt.Sprintf("%s_%s", data.DataBase, data.Table)
	syncInput := make(chan string)
	go SyncTables(data, syncInput, state.syncOutput)
	SyncData := state.stateStorage[StorageKey]
	SyncData.syncChan = syncInput
	state.stateStorage[StorageKey] = SyncData
}
