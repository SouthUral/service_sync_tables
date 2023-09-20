package tools

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Метод для получения переменной окружения по имени (параметр key).
// При отсуствии переменной вернется пустая строка.
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

// Метод для получения переменной окружения по имени (параметр key).
// Переменная преобразуется в int значение.
// При возникновении ошибки или отсутствии переменной вернется 0.
func GetEnvInt(key string) int {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
			return 0
		}
		return val
	}
	return 0
}

// Функция загрузки переменных окружения.
// Если не все переменные загружены, то функция вернет ошибку
func LoadingEnvVars() (EnvironmentsVars, error) {
	envs := EnvironmentsVars{}
	envVars := make(map[string]string)
	messErrors := make([]string, 0)

	envVars["urlStoragePass"] = GetEnv("URL_STORAGE_PASS")
	envVars["mongoHost"] = GetEnv("MONGO_HOST")
	envVars["mongoPort"] = GetEnv("MONGO_PORT")
	envVars["mongoCollection"] = GetEnv("MONGO_COLLECTION")
	envVars["mongoDataBase"] = GetEnv("MONGO_DATABASE")
	envVars["apiServerPort"] = GetEnv("API_SERVER_PORT")
	envVars["logLevel"] = GetEnv("LOG_LEVEL")

	for key, value := range envVars {
		if value == "" {
			switch key {
			case "logLevel":
				continue
			case "urlStoragePass":
				envVars["urlStoragePass"] = "/etc/config/dbconfig.json"
			default:
				messErrors = append(messErrors, key)
			}
		}
	}

	if len(messErrors) > 0 {
		mess := strings.Join(messErrors, ", ")
		errLoading := fmt.Errorf("Не загружены переменные окружения : %s", mess)
		return envs, errLoading
	}

	envs.loadingVal(envVars)
	return envs, nil
}
