package postgres

// Входящее сообщение
type IncomingMess struct {
	Table    string
	Database string
	Offset   string
}

// Канал для входящих сообщений
type IncomCh chan IncomingMess

// Структура для передачи сообщений из горутины с синхронизацией в горутину состояния.
type OutgoingMessSync struct {
	Info   string
	Offset string
	Error  error
	id     string
}

// канал для исходящих сообщений от горутин синхронизаций для модуля state
type OutgoingChanSync chan OutgoingMessSync
