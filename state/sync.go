package state

import (
	"fmt"
	"strconv"
	"time"

	mongo "github.com/SouthUral/service_sync_tables/database/mongodb"

	log "github.com/sirupsen/logrus"
)

// функция запускается в отдельном потоке, ее задача подключиться к БД1 и БД2 и синхронизировать их
// останавливается горутина сообщением из канала inputChan
func SyncTables(data mongo.StateMess, inputChan chan string, outputChan chan syncMessChan) {
	//  здесь должно быть подключение к БД1 и БД2
	// в случае неудачного подключения нужно отправить ошибку в канал outputChan и завершить работу горутины
	var answer string

	intOffset, err := strconv.Atoi(data.Offset)

	if err != nil {
		log.Error(err)
		test_answer := syncMessChan{
			Info:   StartSync,
			Offset: data.Offset,
			id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
			Error:  err.Error(),
		}
		outputChan <- test_answer
		return
	}

	// Отправка сообщения об успешном старте
	// Такое сообщение не будет записываться в DB а сразу уйдет API
	outputChan <- syncMessChan{
		Info:   StartSync,
		Offset: data.Offset,
		id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
		Error:  nil,
	}

	for answer != Stop {

		intOffset++
		test_answer := syncMessChan{
			Info:   RegularSync,
			Offset: strconv.Itoa(intOffset),
			id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
			Error:  nil,
		}
		outputChan <- test_answer
		log.Debug("Сообщение отправлено, offset: ", intOffset)
		answer = <-inputChan
		time.Sleep(1 * time.Second)
	}

	// Отправляет сообщение об остановке синхронизации
	outputChan <- syncMessChan{
		Info:   StopSync,
		Offset: strconv.Itoa(intOffset),
		id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
		Error:  nil,
	}

	// close(outputChan)

	defer log.Debug("Получено сообщение, цикл прекращен: ", answer)
	return
}
