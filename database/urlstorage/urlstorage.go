package urlstorage

// пакет для хранения и работы с url для подключениям к БД

import (
	Config "github.com/SouthUral/service_sync_tables/config"
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

}
