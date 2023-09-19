package urlstorage

// пакет для хранения и работы с url для подключениям к БД

// структура необходимая для работы модуля urlstorage
type urlStorage struct {
	storage        StorageConnDB
	inputChan      InputUrlStorageAPIch
	urlStoragePath string
}

// функция для инициализации urlstorage.
// На вход функция получает канал для получения сообщений
// и map с набором env переменных, которые будут использоваться в этом модуле
func InitUrlStorage(InputCh InputUrlStorageAPIch, UrlStoragePass string) {
	urlStorage := urlStorage{
		storage:        make(StorageConnDB),
		inputChan:      InputCh,
		urlStoragePath: UrlStoragePass,
	}

	go urlStorage.urlMain()
}

// метод для запуска логики модуля в отдельном потоке
func (url *urlStorage) urlMain() {
	url.GetDataFromJson()

	for {
		select {
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
	case ChangeOne:
		url.handlerMessChangeOne(mess)
	case AddOne:
		url.handlerMessAddOne(mess)
	}
}
