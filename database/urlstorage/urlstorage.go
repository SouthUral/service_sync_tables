package urlstorage

// пакет для хранения и работы с url для подключениям к БД

import (
	// "fmt"

	"fmt"

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

// обработчик для сообщений GetAll
func (url *urlStorage) handlerMessGetAll(mess UrlMessInput) {
	switch mess.Message.Format {

	case FormatURL:
		answerMess := AnswerMessAPI{
			Error: ErrorAnswerURL{
				textError: "Данный формат недоступен для множественного вывода",
			},
			AnswerData: nil,
		}
		mess.ReverseCh <- answerMess
		log.Debug("Неверный формат для отправки параметров БД")

	case FormatStruct:
		answerData := CopyMap(url.storage)
		answerMess := AnswerMessAPI{
			Error:      nil,
			AnswerData: answerData,
		}
		mess.ReverseCh <- answerMess
		log.Debug("Данные о параметрах подключения отправлены")

	default:
		answerMess := AnswerMessAPI{
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
	switch mess.Message.Format {

	case FormatURL:
		data, err := url.getOneConn(mess.Message.SearchFor)
		if err != nil {
			answerMess := AnswerMessAPI{
				Error:      err,
				AnswerData: nil,
			}
			mess.ReverseCh <- answerMess
		} else {
			urlConn := CreateUrlFromStruct(data)
			answerMess := AnswerMessAPI{
				Error:      nil,
				AnswerData: urlConn,
			}
			mess.ReverseCh <- answerMess
		}

	case FormatStruct:
		data, err := url.getOneConn(mess.Message.SearchFor)
		if err != nil {
			answerMess := AnswerMessAPI{
				Error:      err,
				AnswerData: nil,
			}
			mess.ReverseCh <- answerMess
		} else {
			answerMess := AnswerMessAPI{
				Error:      nil,
				AnswerData: data,
			}
			mess.ReverseCh <- answerMess
		}

	default:
		answerMess := AnswerMessAPI{
			Error: ErrorAnswerURL{
				textError: "Неизвестный формат",
			},
			AnswerData: nil,
		}
		mess.ReverseCh <- answerMess
		log.Debug("Неизвестный формат")
	}
}

// обработчик для сообщений ChangeOne
func (url *urlStorage) handlerMessChangeOne(mess UrlMessInput) {
	mess.ReverseCh <- AnswerMessAPI{
		Error: ErrorAnswerURL{
			textError: "Метод не готов",
		},
		AnswerData: nil,
	}
}

// обработчик для сообщений AddOne
func (url *urlStorage) handlerMessAddOne(mess UrlMessInput) {
	mess.ReverseCh <- AnswerMessAPI{
		Error: ErrorAnswerURL{
			textError: "Метод не готов",
		},
		AnswerData: nil,
	}
}

func (url *urlStorage) getDBConnData(mess UrlMessInput) {

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

// Функция преобразует парметры подключения к БД в URL
func CreateUrlFromStruct(connData ConnDBData) ConnDBURL {
	return ConnDBURL(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			connData.User,
			connData.Password,
			connData.Host,
			connData.Port,
			connData.NameDB,
		),
	)
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
