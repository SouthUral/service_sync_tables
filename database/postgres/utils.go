package postgres

import (
	"context"
	"fmt"

	url "github.com/SouthUral/service_sync_tables/database/urlstorage"
	pgx "github.com/jackc/pgx/v5"
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

// Метод для отправки сообщения об ошибке в State
func sendErrorMess(mess IncomingMess, err error, outgoingChan OutgoingChanSync, infoMess string) {
	outgoingMess := OutgoingMessSync{
		Info:     infoMess,
		Error:    err,
		Offset:   mess.Offset,
		Database: mess.Database,
		Table:    mess.Table,
	}
	outgoingChan <- outgoingMess
}

// Функция для отправки сообщения об старте или остановке синхронизации
func sendMessForApi(mess IncomingMess, outgoingChan OutgoingChanSync, infoMess string) {
	outgoingMess := OutgoingMessSync{
		Info:     infoMess,
		Offset:   mess.Offset,
		Database: mess.Database,
		Table:    mess.Table,
	}
	outgoingChan <- outgoingMess
	log.Info(fmt.Sprintf("Сообщение об %s отправлено в State", infoMess))
}

// Функция для закрытия коннекта к БД
func closeConn(conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		log.Warning("Коннект к БД уже был закрыт")
		return
	}
	log.Info("Коннект к БД закрыт")
}

// Функция для получения поля наименования ID поля из таблицы
func getFieldIDFromTable(conn *pgx.Conn, schema, table string) (string, error) {
	var NameId string
	query := fmt.Sprintf("SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '%s' and table_name = '%s' and ordinal_position = 1;", schema, table)
	err := conn.QueryRow(context.Background(), query).Scan(&NameId)
	if err != nil {
		return NameId, err
	}
	return NameId, nil
}

// Функция для очистки таблицы перед записью
func cleanTable(conn *pgx.Conn, dataBase, schema, table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s.%s;", schema, table)
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Error(fmt.Sprintf("Таблица %s.%s в БД %s не очищена из-за ошибки : %s", schema, table, dataBase, err.Error()))
		return err
	}
	log.Info(fmt.Sprintf("Таблица %s.%s в БД %s очищена", schema, table, dataBase))
	return nil
}
