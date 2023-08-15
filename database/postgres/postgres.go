package postgres

import (
	"context"
	"fmt"

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
