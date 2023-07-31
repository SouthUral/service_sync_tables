package main

// Канал для отправки сообщения в State
type StateAPIChan chan APImessage

// Структура для упаковки сообщения для State из API
type APImessage struct {
	message string
	ApiChan APImessChan
}

// Канал для возврата сообщения из State в Api
type APImessChan chan StateAnswer

// Структура в которой содержится ответ для API запроса от State
type StateAnswer struct {
	Err  interface{}
	Data StateStorage
}

// Структура для возврата ошибки клиенту
type ErrorResponse struct {
	Status bool        `json:"status"`
	Error  interface{} `json:"error"`
}

// структура возвращается из монго
type StateMess struct {
	oid      string
	Table    string
	DataBase string
	Offset   string
	IsActive bool
}

// Структура для отправки сообщений в каналах между горутиной состояния
// и горутиной MongoDB
type MessCommand struct {
	Info  string
	Data  StateMess
	Error interface{}
}

// Структура для передачи сообщений из горутины с синхронизацией в горутину состояния.
type syncMessChan struct {
	Offset string
	Error  interface{}
	id     string
}

type StateStorage map[string]StateSyncStorage

// Структура для хранении информации о синхронизации таблиц в программе.
type StateSyncStorage struct {
	Id        string
	Table     string
	DataBase  string
	Offset    interface{}
	Err       interface{}
	IsSave    bool
	IsActive  bool
	syncChan  chan string
	DateStart interface{}
	DateEnd   interface{}
}
