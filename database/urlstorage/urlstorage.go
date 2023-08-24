package urlstorage

// пакет для хранения и работы с url для подключениям к БД

import (
	// "fmt"

	Config "github.com/SouthUral/service_sync_tables/config"
	tools "github.com/SouthUral/service_sync_tables/tools"
	log "github.com/sirupsen/logrus"
)

// структура необходимая для работы модуля urlstorage
type urlStorage struct {
	storage        StorageConnDB
	inputChan      InputUrlStorageAPIch
	urlStoragePath string
}

// функция для инициализации urlstorage.
// На вход функция получает канал для получения сообщений
// и map с набором env переменных, которые будут использоваться в этом модуле
func InitUrlStorage(InputCh InputUrlStorageAPIch, ConfVars Config.ConfEnum) {
	urlStorage := urlStorage{
		storage:        make(StorageConnDB),
		inputChan:      InputCh,
		urlStoragePath: ConfVars.UrlStoragePass,
	}

	go urlStorage.urlMain()
}

// метод для запуска логики модуля в отдельном потоке
func (url *urlStorage) urlMain() {
	url.GetDataFromJson()

	for {
		select {
		// ловит сообщение из mdb
		case mess := <-url.inputChan:
			url.processMess(mess)
		}
	}
}

// Метод для распределения сообщений между методами обработчиками
func (url *urlStorage) processMess(mess UrlMessInput) {
	switch mess.Message.Method {
	case GetAll:
		// log.Debug("заглушка processMess.GetAll")
		url.handlerMessGetAll(mess)
	case GetOne:
		log.Debug("заглушка processMess.GetOne")
		url.handlerMessGetOne(mess)
	case ChangeOne:
		log.Debug("заглушка processMess.ChangeOne")
		url.handlerMessChangeOne(mess)
	case AddOne:
		log.Debug("заглушка processMess.AddOne")
		url.handlerMessAddOne(mess)
	}
}

// обработчик для сообщений GetAll
func (url *urlStorage) handlerMessGetAll(mess UrlMessInput) {
	switch mess.Message.Format {

	case FormatURL:
		answerMess := AnswerMessAPI[StorageConnDB]{
			Error: ErrorAnswerURL{
				textError: "Данный формат недоступен для API",
			},
			AnswerData: nil,
		}
		mess.ReverseCh <- answerMess
		log.Debug("Неверный формат для отправки параметров БД")

	case FormatStruct:
		answerData := CopyMap(url.storage)
		answerMess := AnswerMessAPI[StorageConnDB]{
			Error:      nil,
			AnswerData: answerData,
		}
		mess.ReverseCh <- answerMess
		log.Debug("Данные о параметрах подключения отправлены")

	default:
		answerMess := AnswerMessAPI[StorageConnDB]{
			Error: ErrorAnswerURL{
				textError: "Неизвестный формат",
			},
			AnswerData: nil,
		}
		mess.ReverseCh <- answerMess
		log.Debug("Неизвестный формат")
	}
}

// обработчик для сообщений GetOne
func (url *urlStorage) handlerMessGetOne(mess UrlMessInput) {

}

// обработчик для сообщений ChangeOne
func (url *urlStorage) handlerMessChangeOne(mess UrlMessInput) {

}

// обработчик для сообщений AddOne
func (url *urlStorage) handlerMessAddOne(mess UrlMessInput) {

}

func (url *urlStorage) getDBConnData(mess UrlMessInput) {

}

// метод получает данные конфигураций БД и записывает их в urlStorage.storage
func (url *urlStorage) GetDataFromJson() {
	bootJsonData := BootJsonData{}
	err := tools.JsonRead(&bootJsonData, url.urlStoragePath)
	if err != nil {
		log.Error("Не удалось получить данные конфигураций БД")
		return
	}
	for _, item := range bootJsonData {
		url.storage[DBAlias(item.DBAlias)] = item.ConnData
	}
}
