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
