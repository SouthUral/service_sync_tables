package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type State struct {
	table        string
	status       int
	mongoError   interface{}
	StorageChanI StorageChanInput
	mdbInput     chan MessCommand
	mdbOutput    chan MessCommand
	syncOutput   chan syncMessChan
	ApiInputCh   StateAPIChan
	// если нужно отправить StateStorage в целиком в другую горутину
	// то нужно отправлять глубокую копию StateStorage иначе будет состояние гонки
	stateStorage StateStorage
}

// создает структуру State и запускает горутину StateWorker
func InitState(mongoChInput chan MessCommand, mongoChOutput chan MessCommand, ApiCh StateAPIChan) {

	w_state := State{
		StorageChanI: make(StorageChanInput),
		mdbInput:     mongoChInput,
		mdbOutput:    mongoChOutput,
		ApiInputCh:   ApiCh,
		stateStorage: make(StateStorage),
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
		case mess := <-state.ApiInputCh:
			state.ApiHandler(mess)
		}
	}
}

// Функция для обработки сообщений из API
func (state *State) ApiHandler(mess APImessage) {
	switch mess.Message {
	case GetAll:
		mess.ApiChan <- StateAnswer{
			Err:  nil,
			Data: CopyMap(state.stateStorage),
		}
	case InputData:
		key_sync := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
		_, ok := state.stateStorage[key_sync]
		if ok {
			ErrString := fmt.Sprintf("Sync '%s' is already", key_sync)
			mess.ApiChan <- StateAnswer{
				Err: ErrString,
			}
			log.Error(ErrString)
			return
		}
		messComand := MessCommand{
			Info: InputData,
			Data: StateMess{
				Table:    mess.Data.Table,
				DataBase: mess.Data.DataBase,
				Offset:   mess.Data.Offset,
				IsActive: mess.Data.IsActive,
			},
		}
		state.mdbInput <- messComand
		// в словарь StorageChanI записывается канал, до момента получения ответа о записи из mdb
		state.StorageChanI[key_sync] = mess.ApiChan
	case StopSync:
		key_sync := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
		itemSync, ok := state.stateStorage[key_sync]
		if !ok {
			ErrString := fmt.Sprintf("There is no such sync: %s", key_sync)
			mess.ApiChan <- StateAnswer{
				Err: ErrString,
			}
			log.Error(ErrString)
			return
		}
		itemSync.IsActive = false
		state.stateStorage[key_sync] = itemSync
		state.StorageChanI[key_sync] = mess.ApiChan
	case StartSync:
		key_sync := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
		itemSync, ok := state.stateStorage[key_sync]
		if !ok {
			ErrString := fmt.Sprintf("There is no such sync: %s", key_sync)
			mess.ApiChan <- StateAnswer{
				Err: ErrString,
			}
			log.Error(ErrString)
			return
		}
		itemSync.IsActive = true
		state.stateStorage[key_sync] = itemSync
		state.StorageChanI[key_sync] = mess.ApiChan
		state.updateDataMongo(key_sync)
	}

}

// обработчик для сообщений которые приходят из горутин
func (state *State) handlerSyncThreads(mess syncMessChan) {
	itemSync := state.stateStorage[mess.id]
	if mess.Error != nil {
		itemSync.Err = mess.Error
		log.Error(mess.Error)
		return
	}
	itemSync.Offset = mess.Offset
	itemSync.IsSave = false
	state.stateStorage[mess.id] = itemSync
	state.updateDataMongo(mess.id)
	log.Debug("Данные из горутины отправлены на сохранение в MongoDB")
}

// Метод для отправки изменений в состоянии в Mongo.
// Для отправки изменений нужно сначала записать изменения в локальный map stateStorage
// далее вызвать этот метод передав в него ключ
func (state *State) updateDataMongo(id_sync string) {
	itemSync := state.stateStorage[id_sync]
	newMess := MessCommand{
		Info: UpdateData,
		Data: StateMess{
			oid:      itemSync.Id,
			DataBase: itemSync.DataBase,
			Table:    itemSync.Table,
			Offset:   fmt.Sprintf("%s", itemSync.Offset),
			IsActive: itemSync.IsActive,
		},
	}
	state.mdbInput <- newMess
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
	key := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
	itemSync := state.stateStorage[key]
	if mess.Error != nil {
		log.Error("Данные не обновлены в Mongo: ", mess.Error)
		// добавляет данные об ошибке в хранилище
		itemSync.Err = mess.Error
		itemSync.IsSave = false
		itemSync.DateEnd = time.Now()
		state.stateStorage[key] = itemSync
		// отправляет сообщение горутине об остановке
		itemSync.syncChan <- Stop
		return
	}
	// Останавливает синхронизацю если флаг IsActive false
	if itemSync.IsActive == false {
		itemSync.IsSave = true
		state.stateStorage[key] = itemSync
		itemSync.syncChan <- Stop
		itemSync.syncChan = nil
		answ := make(StateStorage)
		answ[key] = itemSync
		ch := state.StorageChanI[key]
		ch <- StateAnswer{
			Err:  nil,
			Data: answ,
		}
		state.stateStorage[key] = itemSync
		return
	}
	itemSync.IsSave = true
	state.stateStorage[key] = itemSync
	// если данные обновлены то в горутину отпрвляется сообщение о продолжении работы
	if itemSync.syncChan != nil {
		itemSync.syncChan <- Continue
	} else {
		state.InitSyncT(mess.Data)
	}
}

// Этот метод запускает горутины синхронизаций!!!
// обработчик сообщений из монго, работает с сообщниями InputData
// запуск горутины произойдет только после записи о синхронизации в mongo
func (state *State) mdbInputData(mess MessCommand) {
	StorageChanKey := fmt.Sprintf("%s_%s", mess.Data.DataBase, mess.Data.Table)
	ch := state.StorageChanI[StorageChanKey]
	if mess.Error != nil {
		log.Error("Данные не добавлены в Mongo: ", mess.Error)
		// отправка сообщения в канал REST о неудачном запуске
		ch <- StateAnswer{
			Err: mess.Error,
		}
		delete(state.StorageChanI, StorageChanKey)
		return
	}
	state.AddInfoToStorage(mess.Data)
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
	state.AddInfoToStorage(mess.Data)
	if mess.Data.IsActive {
		state.InitSyncT(mess.Data)
	}
}

// создает новую запись в словаре stateStorage
func (state *State) AddInfoToStorage(data StateMess) {
	key := fmt.Sprintf("%s_%s", data.DataBase, data.Table)
	state.stateStorage[key] = StateSyncStorage{
		Id:        data.oid,
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
	log.Debug(fmt.Sprintf("%+v\n", state.stateStorage[data.oid]))
}

// функция для запуска горутин синхронизации
func (state *State) InitSyncT(data StateMess) {
	// создает канал для связи с горутиной, запускает горутину и записывает канал в структуру по id
	StorageKey := fmt.Sprintf("%s_%s", data.DataBase, data.Table)
	syncInput := make(chan string)
	go SyncTables(data, syncInput, state.syncOutput)
	SyncData := state.stateStorage[StorageKey]
	SyncData.syncChan = syncInput
	state.stateStorage[StorageKey] = SyncData

	// отправляет сообщение API если есть канал для этого
	ch, ok := state.StorageChanI[StorageKey]
	if ok {
		answerMap := make(StateStorage)
		answerMap[StorageKey] = state.stateStorage[StorageKey]
		ch <- StateAnswer{
			Err:  nil,
			Data: answerMap,
		}
		delete(state.StorageChanI, StorageKey)
	}
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
			id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
			Error:  err,
		}
		outputChan <- test_answer
		log.Debug("Сообщение отправлено, offset: ", newOffset)
		answer = <-inputChan
		time.Sleep(1 * time.Second)
		newOffset++
	}
	// close(outputChan)

	defer log.Debug("Получено сообщение, цикл прекращен: ", answer)
	return
}
