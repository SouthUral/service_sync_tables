package postgres

import (
	"context"
	"fmt"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"

	pgx "github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

// метод для подключения к БД
func pgConnect(dbURL string) (*pgx.Conn, error) {
	connect, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Error(fmt.Sprintf("Connect error: %s", dbURL))
		return connect, err
	}
	log.Info(fmt.Sprintf("Connect is ready: %s", dbURL))
	return connect, nil
}

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
				// Отправить сообщение об ошибке в state
			}
			go pg.mainStreamSync(URLs.urlMainDb, URLs.urlSecondDb)
		}
	}
}

func (pg *postgresMain) mainStreamSync(mainUrl, secondUrl string) {

}
