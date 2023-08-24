package urlstorage

const (
	// Форматы получения параметров подключения к БД
	FormatURL    = "format_url"
	FormatStruct = "format_struct"

	// Методы для выбора логики ответа на сообщение в urlstorage
	GetAll    = "get_all"
	GetOne    = "get_one"
	ChangeOne = "change_one"
	AddOne    = "add_one"
)
