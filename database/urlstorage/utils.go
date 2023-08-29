package urlstorage

import (
	"fmt"
)

func CopyMap(data StorageConnDB) StorageConnDB {
	copyMap := make(StorageConnDB)
	for key, value := range data {
		copyMap[key] = value
	}
	return copyMap
}

// Функция преобразует парметры подключения к БД в URL
func CreateUrlFromStruct(connData ConnDBData) ConnDBURL {
	return ConnDBURL(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			connData.User,
			connData.Password,
			connData.Host,
			connData.Port,
			connData.NameDB,
		),
	)
}

// Создает словарь с URLs вместо структуры ConnDBData
func CreateMapURLs(data StorageConnDB) ConnDBURLs {
	resData := ConnDBURLs{}
	for key, value := range data {
		resData[key] = CreateUrlFromStruct(value)
	}
	return resData
}

// Функция для проверки структуры ConnDBData на заполенность полей.
// Первым элементом возвращает bool=True если все поля заполены bool=False если не все поля заполены.
// Вторым элементом возвращает map[string]string с полями которые заполнены.
func checkConnDBData(data ConnDBData) (bool, map[string]string) {
	var result map[string]string
	fieldIn := true
	for key, item := range data.makeMapStruct() {
		if item != "" {
			result[key] = item
		} else {
			fieldIn = false
		}
	}
	return fieldIn, result
}
