package mongodb

// структура возвращается из монго
type StateMess struct {
	Oid      string
	Table    string
	Schema   string
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

// Канал для получения сообщений из модуля Mongo
type MongoOutputChan chan MessCommand

// Канал для отправки сообщений в модуль Mongo
type MongoInputChan chan MessCommand
