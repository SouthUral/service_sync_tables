package urlstorage

import (
	tools "github.com/SouthUral/service_sync_tables/tools"

	log "github.com/sirupsen/logrus"
)

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
