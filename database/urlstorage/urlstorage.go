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
	inputChan      InputUrlStorageCh
	urlStoragePath string
}

// функция для инициализации urlstorage.
// На вход функция получает канал для получения сообщений
// и map с набором env переменных, которые будут использоваться в этом модуле
func InitUrlStorage(InputCh InputUrlStorageCh, ConfVars Config.ConfEnum) {
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
			url.getUrldata(mess)
		}
	}
}

func (url *urlStorage) getUrldata(mess UrlMessInput) {

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
