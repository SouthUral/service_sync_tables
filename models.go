package main

// type MessageDB struct {
// 	Table    string
// 	DataBase string
// 	Offset   string
// 	IsActive bool
// }

// структура возвращается из монго
type StateMess struct {
	oid      string
	Table    string
	DataBase string
	Offset   string
	IsActive bool
}

type MessCommand struct {
	Info  string
	Data  StateMess
	Error interface{}
}

type syncMessChan struct {
	Offset string
	Error  interface{}
	id     string
}

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

// type MessAnswer struct {
// 	Status string
// 	Data   StateMess
// }
