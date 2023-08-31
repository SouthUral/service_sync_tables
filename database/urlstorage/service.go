package urlstorage

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// обработчик для сообщений GetAll
func (url *urlStorage) handlerMessGetAll(mess UrlMessInput) {
	answerData := CopyMap(url.storage)
	sendAnswerMess(answerData, mess.ReverseCh)
	log.Debug("Данные о параметрах подключения отправлены")
}

// Обработчик для сообщений ChangeOne
func (url *urlStorage) handlerMessChangeOne(mess UrlMessInput) {
	key := mess.Message.ChangeData.DBAlias
	data := mess.Message.ChangeData.ConnData

	urlData, ok := url.storage[DBAlias(key)]
	if !ok {
		errorMess := fmt.Errorf("Указанный alias: %s не найден", key)
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
		url.storage[DBAlias(key)] = urlData
	}

	url.WriteDataToJson()
	sendAnswerMess(nil, mess.ReverseCh)
}

// обработчик для сообщений AddOne
func (url *urlStorage) handlerMessAddOne(mess UrlMessInput) {
	key := mess.Message.ChangeData.DBAlias
	data := mess.Message.ChangeData.ConnData

	check, _ := checkConnDBData(data)
	if check {
		url.storage[DBAlias(key)] = data
	} else {
		sendErrorMess(errors.New("Есть незаполенные обязательные поля"), mess.ReverseCh)
		return
	}
	url.WriteDataToJson()
	sendAnswerMess(nil, mess.ReverseCh)
}

// метод получает данные конфигураций БД и записывает их в urlStorage.storage
func (url *urlStorage) GetDataFromJson() {
	bootJsonData := BootJsonData{}
	err := JsonRead(&bootJsonData, url.urlStoragePath)
	if err != nil {
		log.Error("Не удалось получить данные конфигураций БД")
		return
	}
	for _, item := range bootJsonData {
		url.storage[DBAlias(item.DBAlias)] = item.ConnData
	}
}

// Метод для записи данных в json
func (url *urlStorage) WriteDataToJson() {
	dataForWrite := make(BootJsonData, 0)
	for key, item := range url.storage {
		writeItem := JsonFormat{
			DBAlias:  string(key),
			ConnData: item,
		}
		dataForWrite = append(dataForWrite, writeItem)
	}
	JsonWrite(dataForWrite, url.urlStoragePath)
}

// Функция для отправки сообщения с ошибкой в канал.
func sendErrorMess(err error, ch ReverseAPIch) {
	var answer AnswerMessAPI
	answer.Error = err
	ch <- answer
}

// Функция для отправки ответного сообщения (без ошибки = положительный ответ)
func sendAnswerMess(answerData StorageConnDB, ch ReverseAPIch) {
	answer := AnswerMessAPI{
		Error:      nil,
		AnswerData: answerData,
	}
	ch <- answer
}
