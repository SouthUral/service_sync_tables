package urlstorage

import (
	"fmt"
	// "errors"
)

// канал для получения соообщений в urlstorage от других модулей
type InputUrlStorageAPIch chan UrlMessInput

// Канал для отправки ответа от urlstorage в  API
type ReverseAPIch chan AnswerMessAPI

// Структура для ответа модулю API
type AnswerMessAPI struct {
	Error      error
	AnswerData StorageConnDB
}

// параметры для подключения к БД в формате url
type URLConnData string

// Сообщение которое должен передать модуль API (или другой модуль), с каналом для отправки ответа
type UrlMessInput struct {
	Message   InputMessage
	ReverseCh ReverseAPIch
}

// Сообщение которое будет передано в структуре UrlMessInput для модуля urlstorage
type InputMessage struct {
	Method     string     // метод получить один | получить все | изменить | добавить
	ChangeData JsonFormat // сюда записываются либо данные для изменения либо данные для добавления
	SearchFor  DBAlias    // при установленом методе (получить один) в это поле нужно записать alias, если не указать это поле то?
}

// структура для хранения данных о подключении к БД
type ConnDBData struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	NameDB   string `json:"bd_name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Метод структуры ConnDBData для генерации map[string]string по полям структуры
func (conn *ConnDBData) makeMapStruct() map[string]string {
	return map[string]string{
		"host":     conn.Host,
		"port":     conn.Port,
		"name_db":  conn.NameDB,
		"user":     conn.User,
		"password": conn.Password,
	}
}

// Метод для заполения структура данными из map[string]string
func (conn *ConnDBData) fillingStruct(data map[string]string) error {
	for key, item := range data {
		switch key {
		case "host":
			conn.Host = item
		case "port":
			conn.Port = item
		case "name_db":
			conn.NameDB = item
		case "user":
			conn.User = item
		case "password":
			conn.Password = item
		default:
			return fmt.Errorf("Неизвестный ключ: %s", key)
		}
	}
	return nil
}

// структура для выгрузки и загрузки в json
// Так же эта структура может использоваться для изменения параметров
// подключения со стороны API
type JsonFormat struct {
	DBAlias  string     `json:"db_alias"`
	ConnData ConnDBData `json:"conn_data"`
}

// слайс для загрузочных данных
type BootJsonData []JsonFormat

// словарь содержащий в качестве ключа alias БД а в качетсве значения структуру ConnDBData
type StorageConnDB map[DBAlias]ConnDBData

// Словарь для хранения данных о подключениях в формате URL
type ConnDBURLs map[DBAlias]ConnDBURL

// alias для БД
type DBAlias string

// формат параметров для подключения к БД в формате URL
type ConnDBURL string

type ConfEnum struct {
	UrlStoragePass string `json:"url_storage_path"`
}

type errA interface {
	Error() string
}

type ErrorAnswerURL struct {
	textError string
}

func (errStruct *ErrorAnswerURL) Error() string {
	return errStruct.textError
}

// func Example(N bool) error {
// 	if N {
// 		err := errors.New("some error")
// 		return err
// 	}
// 	return nil
// }
