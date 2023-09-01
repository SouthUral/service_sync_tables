package postgres

import (
	"context"
	"fmt"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"

	pgx "github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

// метод для подключения к БД
func pgConnect(dbURL string) interface{} {
	connect, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Error(fmt.Sprintf("Connect error: %s", dbURL))
		return err
	}
	log.Info(fmt.Sprintf("Connect is ready: %s", dbURL))
	return connect
}

// Структура для хранения каналов, которые переданы при инициализации в main
type postgresMain struct {
	OutgoingChan OutgoingChanSync         // канал для сообщений от горутин синхронизаций
	urlIncomCh   url.InputUrlStorageAPIch // канал для запросов к модулю urlstorage
	pgInputch    IncomCh                  // канал для получения сообщений от других модулей
	mainDB       string                   // alias БД из которой будет производится выгрузка
}

func InitPostgres(urlIncomCh url.InputUrlStorageAPIch, pgInputch IncomCh) {
	pg := postgresMain{
		urlIncomCh: urlIncomCh,
		pgInputch:  pgInputch,
		mainDB:     "main_db",
	}
	// Запуск основного потока postgres
	go pg.mainWorkPg()
	log.Debug(fmt.Sprintf("Создана структура %T", pg))
}

func (pg *postgresMain) mainWorkPg() {
	for {
		select {
		case mess := <-pg.pgInputch:
			mainUrlConn, err := url.GetOneConnURL(pg.mainDB, pg.urlIncomCh)
			if err != nil {
				log.Error(fmt.Sprintf("Не найден URL для mainDB: %s", err.Error()))
				return
			}
			urlConn, err := url.GetOneConnURL(mess.Database, pg.urlIncomCh)
			if err != nil {
				log.Error(err.Error())
				return
			}
			go mainStreamSync(mainUrlConn, urlConn)
		}
	}
}

func mainStreamSync(mainUrl, secondUrl string) {

}
