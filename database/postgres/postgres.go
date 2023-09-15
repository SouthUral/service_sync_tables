package postgres

import (
	"fmt"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"

	log "github.com/sirupsen/logrus"
)

// Структура для хранения каналов, которые переданы при инициализации в main
type postgresMain struct {
	outgoingChan OutgoingChanSync         // канал для сообщений от горутин синхронизаций
	urlIncomCh   url.InputUrlStorageAPIch // канал для запросов к модулю urlstorage
	pgInputch    IncomCh                  // канал для получения сообщений от других модулей
	mainDB       string                   // alias БД из которой будет производится выгрузка
}

// Функция инициализации модуля Postgres, создает и заполняет структуру postgresMain, запускает центральный поток модуля
func InitPostgres(urlIncomCh url.InputUrlStorageAPIch, pgInputch IncomCh, pgOutGoingCh OutgoingChanSync) {
	pg := postgresMain{
		outgoingChan: pgOutGoingCh,
		urlIncomCh:   urlIncomCh,
		pgInputch:    pgInputch,
		mainDB:       "main_db",
	}
	// Запуск основного потока postgres
	go pg.mainWorkPg()
	log.Debug(fmt.Sprintf("Создана структура %T", pg))
}

// Основной поток модуля postgres
func (pg *postgresMain) mainWorkPg() {
	for {
		select {
		case mess := <-pg.pgInputch:
			URLs := checkDBalias(mess, pg.mainDB, pg.urlIncomCh)
			if URLs.err != nil {
				sendErrorMess(mess, URLs.err, pg.outgoingChan, StartSync)
			}
			go pg.mainStreamSync(URLs.urlMainDb, URLs.urlSecondDb, mess)
		}
	}
}

// Поток с проверкой и синхронизацией таблиц
func (pg *postgresMain) mainStreamSync(mainUrl, secondUrl string, incomMess IncomingMess) {
	connects := initConnPg(mainUrl, secondUrl)
	if connects.Error != nil {
		sendErrorMess(incomMess, connects.Error, pg.outgoingChan, StartSync)
		return
	}
	resComparison, checkError := CheckingStrucTables(connects, incomMess)
	if checkError != nil {
		sendErrorMess(incomMess, checkError, pg.outgoingChan, StartSync)
		log.Error(checkError)
		return
	} else if !resComparison {
		comparisonError := fmt.Errorf("Таблицы %s в mainDB и в %s не совпадают", incomMess.Table, incomMess.Database)
		log.Error(comparisonError)
		sendErrorMess(incomMess, comparisonError, pg.outgoingChan, StartSync)
		return
	}
	// Нужно получить поле эквивалентное id, по которому можно отсортироваться
	nameFieldId, err := getFieldIDFromTable(connects.SecondConn, incomMess.Schema, incomMess.Table)
	if err != nil {
		FieldNameError := fmt.Errorf("Имя поля ID не получено из-за ошибки: %s", err.Error())
		log.Error(FieldNameError)
		sendErrorMess(incomMess, FieldNameError, pg.outgoingChan, StartSync)
	}
	// Очищает таблицу перед синхронизацией если флаг IncomingMess.Clean true
	if incomMess.Clean == true {
		cleanTable(connects.SecondConn, incomMess.Database, incomMess.Schema, incomMess.Table)
		// если возникнет ошибка при очистке таблицы, синхронизация не отвалится, но есть риск возникновения ошибки далее
	}
	// запуск цикла синхронизации
	log.Info("Запуск синхронизации")
	shunkTest := "10000"
	sync(connects, incomMess, shunkTest, pg.outgoingChan, nameFieldId)
	return
}
