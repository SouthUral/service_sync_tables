package urlstorage

import (
	"fmt"

	tools "github.com/SouthUral/service_sync_tables/tools"

	log "github.com/sirupsen/logrus"
)

// обработчик для сообщений GetAll
func (url *urlStorage) handlerMessGetAll(mess UrlMessInput) {
	switch mess.Message.Format {

	case FormatURL:
		sendErrorMess("Данный формат недоступен для множественного вывода", mess.ReverseCh)
		log.Debug("Неверный формат для отправки параметров БД")

	case FormatStruct:
		answerData := CopyMap(url.storage)
		sendAnswerMess[StorageConnDB](answerData, mess.ReverseCh)
		log.Debug("Данные о параметрах подключения отправлены")

	default:
		sendErrorMess("Неизвестный формат", mess.ReverseCh)
		log.Debug("Неизвестный формат")
	}
}

// обработчик для сообщений GetOne
func (url *urlStorage) handlerMessGetOne(mess UrlMessInput) {
	switch mess.Message.Format {

	case FormatURL:
		data, err := url.getOneConn(mess.Message.SearchFor)
		if err != nil {
			sendErrorMess(err, mess.ReverseCh)
		} else {
			urlConn := CreateUrlFromStruct(data)
			sendAnswerMess[ConnDBURL](urlConn, mess.ReverseCh)
		}

	case FormatStruct:
		data, err := url.getOneConn(mess.Message.SearchFor)
		if err != nil {
			sendErrorMess(err, mess.ReverseCh)
		} else {
			sendAnswerMess[ConnDBData](data, mess.ReverseCh)
		}

	default:
		sendErrorMess("Неизвестный формат", mess.ReverseCh)
		log.Debug("Неизвестный формат")
	}
}

// Обработчик для сообщений ChangeOne
func (url *urlStorage) handlerMessChangeOne(mess UrlMessInput) {
	key := mess.Message.ChangeData.DBAlias
	data := mess.Message.ChangeData.ConnData

	urlData, ok := url.storage[DBAlias(key)]
	if !ok {
		errorMess := fmt.Sprintf("Указанный alias: %s не найден", key)
		sendErrorMess(errorMess, mess.ReverseCh)
		return
	}

	check, fieldInData := checkConnDBData(data)
	if check {
		url.storage[DBAlias(key)] = data
	} else {
		err := urlData.fillingStruct(fieldInData)
		if err != nil {
			sendErrorMess(err, mess.ReverseCh)
			return
		}
	}

	sendAnswerMess[any](nil, mess.ReverseCh)
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

// Функция для отправки сообщения с ошибкой в канал.
func sendErrorMess(err interface{}, ch ReverseAPIch) {
	var answer AnswerMessAPI
	switch err.(type) {
	case string:
		answer.Error = ErrorAnswerURL{
			textError: fmt.Sprintf("%s", err),
		}
	case ErrorAnswerURL:
		answer.Error = err
	}
	ch <- answer
}

// Функция для отправки ответного сообщения (без ошибки = положительный ответ)
func sendAnswerMess[DataType ConnDBData | any | StorageConnDB | ConnDBURL](answerData DataType, ch ReverseAPIch) {
	answer := AnswerMessAPI{
		Error:      nil,
		AnswerData: answerData,
	}
	ch <- answer
}
