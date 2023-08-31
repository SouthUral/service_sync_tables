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
	urlIncomCh url.InputUrlStorageAPIch
	pgInputch  IncomCh
	mainDB     string // alias БД из которой будет производится выгрузка
}

func InitPostgres(urlIncomCh url.InputUrlStorageAPIch, pgInputch IncomCh) {
	pg := postgresMain{
		urlIncomCh: urlIncomCh,
		pgInputch:  pgInputch,
	}
	log.Debug(fmt.Sprintf("Создана структура %T", pg))
}
