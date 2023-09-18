package tools

import (
	"log"
	"os"
	"strconv"
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
