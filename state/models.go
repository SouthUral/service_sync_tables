package state

import "github.com/SouthUral/service_sync_tables/api"

// словарь для хранения обратного канала API до момента пока придет ответ от mdb
// о добавлении в Mongo информации о новой синхронизации
// хранимый канал нужен для ответа клиенту об усппешном добавлении
type StorageChanInput map[string]api.InputAPIChan

// словарь для хранения всех структур с информацией о синхронизациях
type StateStorage map[string]StateSyncStorage

// Структура для передачи сообщений из горутины с синхронизацией в горутину состояния.
type syncMessChan struct {
	Offset string
	Error  interface{}
	id     string
}

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
