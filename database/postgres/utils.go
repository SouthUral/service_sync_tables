package postgres

import (
	"fmt"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"
	log "github.com/sirupsen/logrus"
)

// Функция генерирует каналы которые необходимы для общения с модулем postgres
func GenerateChan() (OutgoingChanSync, IncomCh) {
	outgoingSyncCh := make(OutgoingChanSync, 100)
	incomCh := make(IncomCh, 100)
	return outgoingSyncCh, incomCh
}

// Метод для получения URL из urlstorage
func checkDBalias(mess IncomingMess, aliasMainDb string, urlCh url.InputUrlStorageAPIch) URLsDB {
	mainUrlConn, err := url.GetOneConnURL(aliasMainDb, urlCh)
	if err != nil {
		log.Error(fmt.Sprintf("Не найден URL для mainDB: %s", err.Error()))
	}
	urlConn, err := url.GetOneConnURL(mess.Database, urlCh)
	if err != nil {
		log.Error(fmt.Sprintf("Нет найдено подключение по переданному ключу: %s", mess.Database))
	}
	res := URLsDB{
		urlMainDb:   mainUrlConn,
		urlSecondDb: urlConn,
		err:         err,
	}
	return res
}
