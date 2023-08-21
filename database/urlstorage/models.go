package urlstorage

// канал для получения соообщений в urlstorage от других модулей
type InputUrlStorageCh chan interface{}

// структура для хранения данных о подключении к БД
type ConnDBData struct {
	Host     string
	Port     string
	NameDB   string
	User     string
	Password string
}

// словарь содержащий в качестве ключа alias БД а в качетсве значения структуру ConnDBData
type StorageConnDB map[DBAlias]ConnDBData

// alias для БД
type DBAlias string
