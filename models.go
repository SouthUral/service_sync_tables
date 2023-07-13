package main

type MessageDB struct {
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
