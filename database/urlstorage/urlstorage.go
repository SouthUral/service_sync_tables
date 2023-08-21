package urlstorage

// пакет для хранения и работы с url для подключениям к БД

// структура необходимая для работы модуля urlstorage
type urlStorage struct {
	storage   StorageConnDB
	inputChan InputUrlStorageCh
}

// функция для инициализации urlstorage, принимает на вход канал для получения сообщений
// и канал для отправки сообщений в Mongo
func InitUrlStorage(InputCh InputUrlStorageCh) {
	urlStorage := urlStorage{
		storage:   make(StorageConnDB),
		inputChan: InputCh,
	}

	go urlStorage.urlMain()
}

// метод для запуска логики модуля в отдельном потоке
func (url *urlStorage) urlMain() {

}
