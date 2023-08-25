package urlstorage

// пакет для хранения и работы с url для подключениям к БД

import (
	// "fmt"

	"fmt"

	Config "github.com/SouthUral/service_sync_tables/config"

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
		url.handlerMessGetAll(mess)
	case GetOne:
		url.handlerMessGetOne(mess)
	case ChangeOne:
		log.Debug("заглушка processMess.ChangeOne")
		url.handlerMessChangeOne(mess)
	case AddOne:
		log.Debug("заглушка processMess.AddOne")
		url.handlerMessAddOne(mess)
	}
}

// Метод для получения парметров одного подключения
func (url *urlStorage) getOneConn(key DBAlias) (ConnDBData, interface{}) {
	data, ok := url.storage[key]
	if ok {
		return data, nil
	} else {
		return ConnDBData{}, ErrorAnswerURL{
			textError: fmt.Sprintf("Не найдены параметры подключения по ключу: %s", key),
		}
	}
}
