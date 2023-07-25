package main

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

// Структура для хранении информации о синхронизации таблиц в программе.
type StateSyncStorage struct {
	id        string
	table     string
	dataBase  string
	offset    interface{}
	err       interface{}
	isSave    bool
	isActive  bool
	syncChan  chan string
	dateStart interface{}
	dateEnd   interface{}
}
