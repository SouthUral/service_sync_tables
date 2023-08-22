package urlstorage

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// Функция для чтения и декодирования данных из json-файла.
// На вход получает тип слайс со структурой, в которую нужно загрузить данные,
// и название файла, в который нужно произвести запись.
// Возвращает пару (результат, ошибку).
func JsonRead(resultTypeData []interface{}, FileName string) (any, any) {
	data, err := os.ReadFile(FileName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = json.Unmarshal(data, &resultTypeData)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return resultTypeData, nil
}

// Функция для кодирования и записи в json-файл.
// На вход получает слайс с заполненными данными структурами, которые будут записаны в файл,
// и название файла, в который нужно произвести запись.
// Возвращает либо nil если запись произошла успешно, либо err.
func JsonWrite(data []interface{}, FileName string) interface{} {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error(err)
		return err
	}

	err = os.WriteFile(FileName, jsonData, 0666)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
