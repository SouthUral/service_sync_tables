package tools

import (
	log "github.com/sirupsen/logrus"
)

type ConfEnum struct {
	UrlStoragePass string `json:"url_storage_path"`
}

// Структура с переменными окружение, которые используются в модуле Mongo.
// Структура заполняется автоматически при вызове метода GetMongoEnvs()
type MongoEnvs struct {
	Host       string
	Port       string
	Collection string
	DataBase   string
}

// Структура для хранения и передачи переменных окружения
type EnvironmentsVars struct {
	UrlStoragePass  string
	MongoHost       string
	MongoPort       string
	MongoCollection string
	MongoDataBase   string
	ApiServerPort   string
	LogLevel        log.Level
}

func (envs *EnvironmentsVars) GetMongoEnvs() MongoEnvs {
	return MongoEnvs{
		Host:       envs.MongoHost,
		Port:       envs.MongoPort,
		Collection: envs.MongoCollection,
		DataBase:   envs.MongoDataBase,
	}
}

// метод для загрузки значений из map словаря в структуру EnvironmentsVars
func (envs *EnvironmentsVars) loadingVal(vals map[string]string) {
	envs.UrlStoragePass = vals["urlStoragePass"]
	envs.MongoHost = vals["mongoHost"]
	envs.MongoPort = vals["mongoPort"]
	envs.MongoCollection = vals["mongoCollection"]
	envs.MongoDataBase = vals["mongoDataBase"]
	envs.ApiServerPort = vals["apiServerPort"]

	switch vals["logLevel"] {
	case "DebugLevel":
		envs.LogLevel = log.DebugLevel
	case "ErrorLevel":
		envs.LogLevel = log.ErrorLevel
	case "WarnLevel":
		envs.LogLevel = log.WarnLevel
	case "TraceLevel":
		envs.LogLevel = log.TraceLevel
	case "PanicLevel":
		envs.LogLevel = log.PanicLevel
	default:
		envs.LogLevel = log.InfoLevel
	}

}
