package urlstorage

import "fmt"

// Функция для отправки сообщения в urlstorage и получения ответа
func SendingMess(mess InputMessage, outputCh InputUrlStorageAPIch) (StorageConnDB, any) {
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
func AllConn(outputCh InputUrlStorageAPIch) (StorageConnDB, any) {
	outputMess := InputMessage{
		Method: GetAll,
	}
	answer, err := SendingMess(outputMess, outputCh)
	return answer, err
}

// Внешний метод для получения параметров одного подключения
func GetOneConn(alias string, outputCh InputUrlStorageAPIch) (ConnDBData, any) {
	answer, err := AllConn(outputCh)
	if err != nil {
		return ConnDBData{}, err
	}
	value, ok := answer[DBAlias(alias)]
	if ok {
		return value, nil
	}
	return ConnDBData{}, ErrorAnswerURL{
		textError: fmt.Sprintf("Конфигурация для подключения по ключу %s не найдена", alias),
	}
}

// switch format {
// case FormatStruct:
// 	return answer, err
// case FormatURL:
// 	answerURLs := CreateMapURLs(answer)
// 	return answerURLs, err
// default:
// 	err := ErrorAnswerURL{
// 		textError: "Неверный формат",
// 	}
// 	return nil, err
// }
