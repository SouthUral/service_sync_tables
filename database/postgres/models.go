package postgres

import (
	"fmt"

	pgx "github.com/jackc/pgx/v5"
)

// Структура для внутренней передачи коннектов к БД между функциями модуля Postgres
type ConnectsPG struct {
	MainConn   *pgx.Conn
	SecondConn *pgx.Conn
	Error      error
}

// Входящее сообщение
type IncomingMess struct {
	Table      string
	Schema     string
	Database   string
	Offset     string
	Clean      bool
	ChCommSync CommToSync
}

// Канал для входящих сообщений
type IncomCh chan IncomingMess

// Канал для отправки сообщений горутине синхронизации в процессе ее работы
type CommToSync chan string

// Структура для передачи сообщений из горутины с синхронизацией в горутину состояния.
type OutgoingMessSync struct {
	Info     string
	Offset   string
	Error    error
	Database string
	Schema   string
	Table    string
}

func (mess *OutgoingMessSync) GetID() string {
	return fmt.Sprintf("%s_%s", mess.Database, mess.Table)
}

type URLsDB struct {
	urlMainDb   string
	urlSecondDb string
	err         error
}

// канал для исходящих сообщений от горутин синхронизаций для модуля state
type OutgoingChanSync chan OutgoingMessSync
