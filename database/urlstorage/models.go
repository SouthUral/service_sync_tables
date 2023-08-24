package urlstorage

// канал для получения соообщений в urlstorage от других модулей
type InputUrlStorageAPIch chan UrlMessInput

// Канал для отправки ответа от urlstorage в  API
type ReverseAPIch chan AnswerMessAPI

// Структура для ответа модулю API
type AnswerMessAPI struct {
	Error      interface{}
	AnswerData any
}

type AnswerData[typeAnswer StorageConnDB | URLConnData | any] struct {
	data typeAnswer
}

// параметры для подключения к БД в формате url
type URLConnData string

// Сообщение которое должен передать модуль API (или другой модуль), с каналом для отправки ответа
type UrlMessInput struct {
	Message   InputMessage[any]
	ReverseCh ReverseAPIch
}

// Сообщение которое будет передано в структуре UrlMessInput для модуля urlstorage
type InputMessage[dataForChange JsonFormat | any] struct {
	Method     string        // метод получить один | получить все | изменить | добавить
	Format     string        // формат для получения (json или url)
	ChangeData dataForChange // сюда записываются либо данные для изменения либо данные для добавления
	SearchFor  DBAlias       // при установленом методе (получить один) в это поле нужно записать alias, если не указать это поле то?
}

// структура для хранения данных о подключении к БД
type ConnDBData struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	NameDB   string `json:"bd_name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// структура для выгрузки и загрузки в json
type JsonFormat struct {
	DBAlias  string     `json:"db_alias"`
	ConnData ConnDBData `json:"conn_data"`
}

// слайс для загрузочных данных
type BootJsonData []JsonFormat

// словарь содержащий в качестве ключа alias БД а в качетсве значения структуру ConnDBData
type StorageConnDB map[DBAlias]ConnDBData

// alias для БД
type DBAlias string

// формат параметров для подключения к БД в формате URL
type ConnDBURL string

type ConfEnum struct {
	UrlStoragePass string `json:"url_storage_path"`
}

type ErrorAnswerURL struct {
	textError string
}

func (errStruct *ErrorAnswerURL) Error() string {
	return errStruct.textError
}
