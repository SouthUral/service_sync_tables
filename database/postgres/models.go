package postgres

// Входящее сообщение
type IncomingMess struct {
	Table    string
	Database string
	Offset   string
}

// Канал для входящих сообщений
type IncomCh chan IncomingMess
