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
		log.Error(fmt.Sprintf("Connect error: %s, %s", dbURL, err.Error()))
		return connect, err
	}
	if connect == nil {
		err := fmt.Errorf(fmt.Sprintf("Connect error: %s", dbURL))
		log.Error(err)
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

func StartConcurentConn(dbURL string) ChConnMess {
	ConnCh := make(ChConnMess)
	go ConcurentConnPg(dbURL, ConnCh)
	return ConnCh
}

func initConnPg(URLmainDB, URLsecondDb string) ConnectsPG {
	ChMain := StartConcurentConn(URLmainDB)
	ChSecond := StartConcurentConn(URLsecondDb)

	errString := ""
	answer := ConnectsPG{}
	checkAnswer := 0

	messMainConn := ConnMess{}
	messSecondConn := ConnMess{}
	for checkAnswer < 2 {
		select {

		case messMainConn = <-ChMain:
			answer.MainConn = messMainConn.Conn
			mainError := messMainConn.ErrConn
			if mainError != nil {
				errString = fmt.Sprintf("Ошибка подключения к основной БД: %s.\n%s", mainError.Error(), errString)
			}
			checkAnswer++
		case messSecondConn = <-ChSecond:
			answer.SecondConn = messSecondConn.Conn
			secondError := messSecondConn.ErrConn
			if secondError != nil {
				errString = fmt.Sprintf("Ошибка подключения к БД: %s.\n%s", secondError.Error(), errString)
			}
			checkAnswer++
		default:
			continue
		}
	}

	if errString != "" {
		answer.Error = fmt.Errorf(errString)
		// Закрытие коннектов если хотя бы один коннект не сработал
		if messMainConn.ErrConn == nil {
			closeConn(messMainConn.Conn)
		}
		if messSecondConn.ErrConn == nil {
			closeConn(messSecondConn.Conn)
		}
	} else {
		answer.Error = nil
	}
	return answer
}
