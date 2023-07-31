package main

// Канал для отправки сообщения в State
type StateAPIChan chan APImessage

// Структура для упаковки сообщения для State из API
type APImessage struct {
	Data    InputDataApi
	Message string
	ApiChan APImessChan
}

type InputDataApi struct {
	Table    string
	DataBase string
	IsActive bool
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
	Id        string      `json:"id"`
	Table     string      `json:"table"`
	DataBase  string      `json:"data_base"`
	Offset    interface{} `json:"offset"`
	Err       interface{} `json:"error"`
	IsSave    bool        `json:"is_save"`
	IsActive  bool        `json:"is_active"`
	syncChan  chan string
	DateStart interface{} `json:"date_start"`
	DateEnd   interface{} `json:"date_end"`
}
