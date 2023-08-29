package api

// Канал для отправки сообщения в State
type OutputAPIChan chan APImessage

// Канал для возврата сообщения из State в Api
type InputAPIChan chan StateAnswer

// Структура в которой содержится ответ для API запроса от State
type StateAnswer struct {
	Info string
	Err  interface{}
	Data interface{}
}

// Структура для получения и расшифровки данных от клиента
type InputDataApi struct {
	Table    string `json:"table"`
	DataBase string `json:"data_base"`
	IsActive bool   `json:"is_active"`
	Offset   string `json:"offset"`
}

// Структура для упаковки сообщения для State из API
type APImessage struct {
	Data    InputDataApi
	Message string
	ApiChan InputAPIChan
}

// Структура для возврата ошибки клиенту
type ErrorResponse struct {
	Status bool        `json:"status"`
	Error  interface{} `json:"error"`
}

type RequestDBConn struct {
	Alias string `json:"alias"`
}

// // Канал для возврата сообщения из State в Api
// type APImessChan chan StateAnswer
