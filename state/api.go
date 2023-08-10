package state

import (
	"fmt"

	api "github.com/SouthUral/service_sync_tables/api"
	mongo "github.com/SouthUral/service_sync_tables/database/mongodb"

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
				DataBase: mess.Data.DataBase,
				Offset:   mess.Data.Offset,
				IsActive: mess.Data.IsActive,
			},
		}
		state.mdbInput <- messComand
		// в словарь StorageChanI записывается канал, до момента получения ответа о записи из mdb
		state.StorageChanI[key_sync] = mess.ApiChan
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
		state.StorageChanI[key_sync] = mess.ApiChan
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
				Err: "sync is already start",
			}
			return
		}
		itemSync.IsActive = true
		state.stateStorage[key_sync] = itemSync
		state.StorageChanI[key_sync] = mess.ApiChan
		state.updateDataMongo(key_sync)
	}

}

// отправляет сообщение API если есть канал для этого
func (state *State) ResponseAPIRequest(key string, err interface{}, status string) {
	ch, ok := state.StorageChanI[key]
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
	case StartSync:
		if err != nil {
			apiMess.Info = "sync did not start due to an error"
			apiMess.Err = err
		} else {
			apiMess.Info = "sync has started successfully"
		}
	case StopSync:
		apiMess.Info = "sync has been stopped"
	}

	ch <- apiMess

	delete(state.StorageChanI, key)
}
