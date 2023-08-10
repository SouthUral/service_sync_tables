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
	newOffset := intOffset + 1
	if err != nil {
		log.Error(err)
		test_answer := syncMessChan{
			Offset: "0",
			id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
			Error:  err.Error(),
		}
		outputChan <- test_answer
		return
	}

	for answer != Stop {

		test_answer := syncMessChan{
			Offset: strconv.Itoa(newOffset),
			id:     fmt.Sprintf("%s_%s", data.DataBase, data.Table),
			Error:  nil,
		}
		outputChan <- test_answer
		log.Debug("Сообщение отправлено, offset: ", newOffset)
		answer = <-inputChan
		time.Sleep(1 * time.Second)
		newOffset++
	}
	// close(outputChan)

	defer log.Debug("Получено сообщение, цикл прекращен: ", answer)
	return
}
