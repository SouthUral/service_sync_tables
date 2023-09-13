package state

import (
	"fmt"

	api "github.com/SouthUral/service_sync_tables/api"
	mongo "github.com/SouthUral/service_sync_tables/database/mongodb"

	log "github.com/sirupsen/logrus"
)

// Обработчик сообщений приходящих от модуля MongoDB
func (state *State) MongoWorker(mess mongo.MessCommand) {
	switch mess.Info {
	case mongo.GetAll:
		state.mdbGetAll(mess)
	case mongo.InputData:
		state.mdbInputData(mess)
	case mongo.UpdateData:
		state.mdbUpdateData(mess)
	}
}

// обработчик сообщений из монго, работает с сообщниями GetAll
func (state *State) mdbGetAll(mess mongo.MessCommand) {
	if mess.Error != nil {
		log.Error("Старт синхронизации не состоялся по причине: ", mess.Error)
		state.mongoError = mess.Error
		return
	}
	state.AddInfoToStorage(mess.Data)
	if mess.Data.IsActive {
		state.InitSyncT(mess.Data)
	}
}

// Этот метод запускает горутины синхронизаций!!!
// обработчик сообщений из монго, работает с сообщниями InputData
// запуск горутины произойдет только после записи о синхронизации в mongo
func (state *State) mdbInputData(mess mongo.MessCommand) {
	StorageChanKey := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
	if mess.Error != nil {
		log.Error("Данные не добавлены в Mongo: ", mess.Error)
		state.ResponseAPIRequest(StorageChanKey, mess.Error, api.InputData)
		return
	}
	state.AddInfoToStorage(mess.Data)
	if mess.Data.IsActive {
		state.InitSyncT(mess.Data)
	}
}

// метод обработчик для сообщений UpdateData из модуля MongoDB
func (state *State) mdbUpdateData(mess mongo.MessCommand) {
	key := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
	itemSync := state.stateStorage[key]
	if mess.Error != nil {
		log.Error("Данные не обновлены в Mongo: ", mess.Error)
		state.StopSyncState(key, mess.Error, true)
		return
	}

	// Останавливает синхронизацю если флаг IsActive false
	if itemSync.IsActive == false {
		state.StopSyncState(key, nil, true)
		return
	}

	itemSync.IsSave = true
	state.stateStorage[key] = itemSync

	// если данные обновлены то в горутину отпрвляется сообщение о продолжении работы
	// если нет у sync нет канала (она была не активна) она инициализируется заного
	if itemSync.syncChan != nil {
		itemSync.syncChan <- Continue
	} else {
		state.InitSyncT(mess.Data)
	}
}

// Метод для отправки изменений в состоянии в Mongo.
// Для отправки изменений нужно сначала записать изменения в локальный map stateStorage
// далее вызвать этот метод передав в него ключ
func (state *State) updateDataMongo(id_sync string) {
	itemSync := state.stateStorage[id_sync]
	newMess := mongo.MessCommand{
		Info: mongo.UpdateData,
		Data: mongo.StateMess{
			Oid:      id_sync,
			DataBase: itemSync.DataBase,
			Schema:   itemSync.Schema,
			Table:    itemSync.Table,
			Offset:   fmt.Sprintf("%s", itemSync.Offset),
			IsActive: itemSync.IsActive,
		},
	}
	state.mdbInput <- newMess
}
