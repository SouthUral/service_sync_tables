package main

// type MessageDB struct {
// 	Table    string
// 	DataBase string
// 	Offset   string
// 	IsActive bool
// }

// структура возвращается из монго
type StateMess struct {
	oid      interface{}
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
	id string
}

// type MessAnswer struct {
// 	Status string
// 	Data   StateMess
// }
