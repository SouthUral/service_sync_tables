package postgres

// Входящее сообщение
type IncomingMess struct {
	Table      string
	Database   string
	Offset     string
	ChCommSync CommToSync
}

// Канал для входящих сообщений
type IncomCh chan IncomingMess

// Канал для отправке сообщений горутине синхронизации в процессе ее работы
type CommToSync chan string

// Структура для передачи сообщений из горутины с синхронизацией в горутину состояния.
type OutgoingMessSync struct {
	Info     string
	Offset   string
	Error    error
	Database string
	Table    string
}

type URLsDB struct {
	urlMainDb   string
	urlSecondDb string
	err         error
}

// канал для исходящих сообщений от горутин синхронизаций для модуля state
type OutgoingChanSync chan OutgoingMessSync
