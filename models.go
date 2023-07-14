package main

type MessageDB struct {
	Table    string
	DataBase string
	Offset   string
	IsActive bool
}

// структура возвращается из монго
type MessDBAnsw struct {
	oid      string
	Table    string
	DataBase string
	Offset   string
	IsActive bool
}

type MessCommand struct {
	Command string
	Data    MessageDB
}

type MessAnswer struct {
	Status string
	Data   interface{}
}
