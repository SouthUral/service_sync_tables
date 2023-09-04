package postgres

import (
	"context"
	"fmt"

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

type ConnMess struct {
	Conn    *pgx.Conn
	ErrConn error
}

type ChConnMess chan ConnMess

// Функция для запуска в горутине
func ConcurentConnPg(dbURL string, ch ChConnMess) {
	conn, err := pgConnect(dbURL)
	ch <- ConnMess{
		Conn:    conn,
		ErrConn: err,
	}
}

func StartConcurentConn(dbURL string) (*pgx.Conn, error) {
	ConnCh := make(ChConnMess)
	go ConcurentConnPg(dbURL, ConnCh)
	mess := <-ConnCh
	return mess.Conn, mess.ErrConn
}

type ConnectsPG struct {
	MainConn   *pgx.Conn
	SecondConn *pgx.Conn
	Error      error
}

func initConnPg(URLmainDB, URLsecondDb string) ConnectsPG {
	mainConn, errMain := StartConcurentConn(URLmainDB)
	secondConn, errSecond := StartConcurentConn(URLsecondDb)

	var answer ConnectsPG
	var errString string
	if errMain != nil {
		errString = fmt.Sprintf("Ошибка подключения к основной БД: %s.", errMain.Error())
	}
	if errSecond != nil {
		errString = fmt.Sprintf("Ошибка подключения к БД: %s.\n%s", errSecond.Error(), errString)
	}
	if errString != "" {
		answer = ConnectsPG{
			MainConn:   mainConn,
			SecondConn: secondConn,
			Error:      fmt.Errorf(errString),
		}
	} else {
		answer = ConnectsPG{
			MainConn:   mainConn,
			SecondConn: secondConn,
			Error:      nil,
		}
	}
	return answer
}
