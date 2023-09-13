package state

import (
	"fmt"
	"time"

	api "github.com/SouthUral/service_sync_tables/api"
	mongo "github.com/SouthUral/service_sync_tables/database/mongodb"
	pg "github.com/SouthUral/service_sync_tables/database/postgres"

	log "github.com/sirupsen/logrus"
)

type State struct {
	table         string
	status        int
	mongoError    interface{}
	StorageChanI  StorageChanInput
	mdbInput      mongo.MongoInputChan
	mdbOutput     mongo.MongoOutputChan
	syncOutput    pg.OutgoingChanSync
	outputApiChan api.OutputAPIChan
	chanStartSync pg.IncomCh
	// если нужно отправить StateStorage в целиком в другую горутину
	// то нужно отправлять глубокую копию StateStorage иначе будет состояние гонки
	stateStorage StateStorage
}

// создает структуру State и запускает горутину StateWorker
func InitState(mdbInput mongo.MongoInputChan, mdbOutput mongo.MongoOutputChan, outputApiChan api.OutputAPIChan, syncOutput pg.OutgoingChanSync, chanStartSync pg.IncomCh) {

	w_state := State{
		StorageChanI:  make(StorageChanInput),
		mdbInput:      mdbInput,
		mdbOutput:     mdbOutput,
		outputApiChan: outputApiChan,
		stateStorage:  make(StateStorage),
		syncOutput:    syncOutput,
		chanStartSync: chanStartSync,
	}
	go w_state.StateWorker()
}

// основная функция для работы с состоянием, при запуске она стартует и ждет ответа от монго
// состояние должно принимать парметрами 6 каналов, которые создаются в main
// для обработки сообщений с каждым из модулей можно создать функции (3 функции)
func (state *State) StateWorker() {
	// метод для старта, запрашивает из монго все документы
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
func (state *State) HandlerSyncThreads(mess pg.OutgoingMessSync) {
	itemSync := state.stateStorage[mess.GetID()]

	// отправляет сообщение в API либо об успешном страте либо об ошибке
	switch mess.Info {
	case pg.StartSync:
		if mess.Error == nil {
			state.ResponseAPIRequest(mess.GetID(), nil, pg.StartSync)
		} else {
			state.StopSyncState(mess.GetID(), mess.Error, false)
			state.ResponseAPIRequest(mess.GetID(), mess.Error, pg.StartSync)
		}
		return
	case pg.RegularSync:
		if mess.Error != nil {
			state.StopSyncState(mess.GetID(), mess.Error, false)
			state.ResponseAPIRequest(mess.GetID(), mess.Error, pg.StartSync)
			return
		}
		itemSync.Offset = mess.Offset
		itemSync.IsSave = false
		state.stateStorage[mess.GetID()] = itemSync
		state.updateDataMongo(mess.GetID())
		log.Debug("Данные из горутины отправлены на сохранение в MongoDB")
	case pg.StopSync:
		state.ResponseAPIRequest(mess.GetID(), nil, pg.StopSync)
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
		Schema:    data.Schema,
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

	chanSync := pg.StartSyncPg(state.chanStartSync, data.DataBase, data.Table, data.Schema, data.Offset)

	SyncData := state.stateStorage[StorageKey]
	SyncData.syncChan = chanSync
	state.stateStorage[StorageKey] = SyncData
}
