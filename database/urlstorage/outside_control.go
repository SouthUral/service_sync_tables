package urlstorage

import (
	"fmt"
)

// Функция для отправки сообщения в urlstorage и получения ответа
func SendingMess(mess InputMessage, outputCh InputUrlStorageAPIch) (StorageConnDB, error) {
	reverseCh := make(ReverseAPIch, 0)
	outgoingMess := UrlMessInput{
		Message:   mess,
		ReverseCh: reverseCh,
	}
	outputCh <- outgoingMess
	answer, _ := <-reverseCh
	if answer.Error != nil {
		return nil, answer.Error
	}
	return answer.AnswerData, nil
}

// Внешний метод для получения параметров всех БД
func AllConn(outputCh InputUrlStorageAPIch) (StorageConnDB, error) {
	outputMess := InputMessage{
		Method: GetAll,
	}
	answer, err := SendingMess(outputMess, outputCh)
	return answer, err
}

// Внешний метод для получения параметров одного подключения
func GetOneConn(alias string, outputCh InputUrlStorageAPIch) (ConnDBData, error) {
	answer, err := AllConn(outputCh)
	if err != nil {
		return ConnDBData{}, err
	}
	value, ok := answer[DBAlias(alias)]
	if ok {
		return value, nil
	}
	return ConnDBData{}, fmt.Errorf("Конфигурация для подключения по ключу %s не найдена", alias)
}

// Внешний метод для получения параметров одного подключения в виде URL
func GetOneConnURL(alias string, outputCh InputUrlStorageAPIch) (string, error) {
	resConn, err := GetOneConn(alias, outputCh)
	if err != nil {
		return "", err
	}
	URL := CreateUrlFromStruct(resConn)
	return string(URL), nil
}

// Внешний метод для изменения параметров подключения
func ChangeOneConn(data JsonFormat, outputCh InputUrlStorageAPIch) error {
	return absAddMethod(data, outputCh, ChangeOne)
}

// Внешний метод для добавления параметров подключения
func AddOneConn(data JsonFormat, outputCh InputUrlStorageAPIch) error {
	return absAddMethod(data, outputCh, AddOne)
}

// Абстрактный метод для изменения или добавления данных подключения
func absAddMethod(data JsonFormat, outputCh InputUrlStorageAPIch, mess string) error {
	outputMess := InputMessage{
		Method:     mess,
		ChangeData: data,
	}
	_, err := SendingMess(outputMess, outputCh)
	return err
}
