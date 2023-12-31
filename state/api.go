package state

import (
	"fmt"

	api "github.com/SouthUral/service_sync_tables/api"
	mongo "github.com/SouthUral/service_sync_tables/database/mongodb"
	pg "github.com/SouthUral/service_sync_tables/database/postgres"

	log "github.com/sirupsen/logrus"
)

// Функция для обработки сообщений из API
func (state *State) ApiHandler(mess api.APImessage) {
	switch mess.Message {
	case api.GetAll:
		mess.ApiChan <- api.StateAnswer{
			Err:  nil,
			Data: CopyMap(state.stateStorage),
		}
	case api.InputData:
		key_sync := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
		_, ok := state.stateStorage[key_sync]
		if ok {
			ErrString := fmt.Sprintf("Sync '%s' is already", key_sync)
			mess.ApiChan <- api.StateAnswer{
				Err: ErrString,
			}
			log.Error(ErrString)
			return
		}
		messComand := mongo.MessCommand{
			Info: mongo.InputData,
			Data: mongo.StateMess{
				Table:    mess.Data.Table,
				Schema:   mess.Data.Schema,
				DataBase: mess.Data.DataBase,
				Offset:   mess.Data.Offset,
				IsActive: mess.Data.IsActive,
				Clean:    mess.Data.Clean,
			},
		}
		state.mdbInput <- messComand
		// в словарь StorageChanI записывается канал, до момента получения ответа о записи из mdb
		state.StorageChanI[key_sync] = &CountChanUse{
			Chanal:  mess.ApiChan,
			Сounter: 1,
		}
	case api.StopSync:
		key_sync := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
		itemSync, ok := state.stateStorage[key_sync]
		if !ok {
			ErrString := fmt.Sprintf("There is no such sync: %s", key_sync)
			mess.ApiChan <- api.StateAnswer{
				Err: ErrString,
			}
			log.Error(ErrString)
			return
		}
		if itemSync.IsActive == false {
			mess.ApiChan <- api.StateAnswer{
				Err: "sync is already stop",
			}
			return
		}
		itemSync.IsActive = false
		state.stateStorage[key_sync] = itemSync
		state.StorageChanI[key_sync] = &CountChanUse{
			Chanal:  mess.ApiChan,
			Сounter: 1,
		}
		state.updateDataMongo(key_sync)
	case api.StartSync:
		key_sync := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
		itemSync, ok := state.stateStorage[key_sync]
		if !ok {
			ErrString := fmt.Sprintf("There is no such sync: %s", key_sync)
			mess.ApiChan <- api.StateAnswer{
				Err: ErrString,
			}
			log.Error(ErrString)
			return
		}
		if itemSync.IsActive == true {
			mess.ApiChan <- api.StateAnswer{
				Err:  "sync is already start",
				Data: itemSync,
			}
			return
		}
		itemSync.IsActive = true
		state.stateStorage[key_sync] = itemSync
		state.StorageChanI[key_sync] = &CountChanUse{
			Chanal:  mess.ApiChan,
			Сounter: 1,
		}
		state.updateDataMongo(key_sync)
	case api.StartAll:
		state.apiChangeActiveSyncs(mess, true)
	case api.StopAll:
		state.apiChangeActiveSyncs(mess, false)
	}
}

// метод для изменения состояния всех синхронизаций
func (state *State) apiChangeActiveSyncs(mess api.APImessage, active bool) {
	counterChanUse := CountChanUse{
		Chanal:  mess.ApiChan,
		Сounter: 0,
	}
	keySyncItems := make([]string, 0)
	for key, itemSync := range state.stateStorage {
		switch active {
		case true:
			if itemSync.IsActive {
				continue
			}
			itemSync.IsActive = true
		case false:
			if !itemSync.IsActive {
				continue
			}
			itemSync.IsActive = false
		}

		counterChanUse.Сounter++
		state.stateStorage[key] = itemSync

		keySyncItems = append(keySyncItems, key)
	}

	// Отправка данных об уже остановленных или уже запущенных синхронизациях
	// После отправки канал обязательно должен быть закрыт иначе API заблокируется и не отдаст ответа
	if counterChanUse.Сounter == 0 {
		switch active {
		case true:
			mess.ApiChan <- api.StateAnswer{
				Info: "Все синхронизации уже запущены",
			}
		case false:
			mess.ApiChan <- api.StateAnswer{
				Info: "Все синхронизации уже остановлены",
			}
		}
		close(mess.ApiChan)
		return
	}

	for _, key := range keySyncItems {
		state.StorageChanI[key] = &counterChanUse
		state.updateDataMongo(key)
	}
}

// отправляет сообщение API если есть канал для этого
func (state *State) ResponseAPIRequest(key string, err interface{}, status string) {
	chanCount, ok := state.StorageChanI[key]
	if !ok {
		log.Info("Channel not found, API message not sent")
		return
	}
	answerMap := make(StateStorage)
	answerMap[key] = state.stateStorage[key]

	apiMess := api.StateAnswer{
		Data: answerMap,
	}

	switch status {
	case pg.StartSync:
		if err != nil {
			apiMess.Info = "sync did not start due to an error"
			apiMess.Err = err
		} else {
			apiMess.Info = "sync has started successfully"
		}
	case pg.StopSync:
		apiMess.Info = "sync has been stopped"
	}

	chanCount.Chanal <- apiMess

	// проверка сколько сообщений ожидает канал, если остается одно сообщение
	// то оно отправляется и канал закрывается
	if chanCount.Сounter > 1 {
		chanCount.Сounter--
	} else {
		chanCount.Сounter--
		close(chanCount.Chanal)
	}

	// Удаление канала из словаря по ключу,
	// канал не будет полностью удален из словаря пока есть ключи имеющие ссылки на канал
	delete(state.StorageChanI, key)
}
