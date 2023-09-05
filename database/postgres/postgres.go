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

func (pg *postgresMain) mainStreamSync(mainUrl, secondUrl string, incomMess IncomingMess) {
	connects := initConnPg(mainUrl, secondUrl)
	if connects.Error != nil {
		sendErrorMess(incomMess, connects.Error, pg.outgoingChan, StartSync)
		return
	}
	resComparison, checkError := CheckingStrucTables(connects, incomMess)
	if checkError != nil {
		sendErrorMess(incomMess, checkError, pg.outgoingChan, StartSync)
		return
	} else if !resComparison {
		comparisonError := fmt.Errorf("Таблицы %s в mainDB и в %s не совпадают", incomMess.Table, incomMess.Database)
		sendErrorMess(incomMess, comparisonError, pg.outgoingChan, StartSync)
	}

}
