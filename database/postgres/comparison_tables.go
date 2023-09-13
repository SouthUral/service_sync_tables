package postgres

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func CheckingStrucTables(connects ConnectsPG, mess IncomingMess) (bool, error) {
	chMainDb := StartGettingStructure(connects.MainConn, mess)
	chSeconDb := StartGettingStructure(connects.SecondConn, mess)

	structsTables := StructsTablesDb{}
	var errorMess string

	select {

	case messMainAnswer := <-chMainDb:
		if messMainAnswer.errorAnswer != nil {
			errorMess = fmt.Sprintf("Ошибка чтения структуры таблицы в mainDB: %s.\n%s", messMainAnswer.errorAnswer.Error(), errorMess)
		}
		structsTables.mainStructTable = messMainAnswer.answer

	case messSecondAnswer := <-chSeconDb:
		if messSecondAnswer.errorAnswer != nil {
			errorMess = fmt.Sprintf("Ошибка чтения структуры таблицы в secondDB: %s.\n%s", messSecondAnswer.errorAnswer.Error(), errorMess)
		}
		structsTables.secondStructTable = messSecondAnswer.answer
	}
	if errorMess != "" {
		err := fmt.Errorf(errorMess)
		return false, err
	}
	result := ComparisonTables(structsTables)
	return result, nil
}

func ComparisonTables(tables StructsTablesDb) bool {
	if len(tables.mainStructTable) != len(tables.secondStructTable) {
		return false
	}
	for key, value := range tables.mainStructTable {
		valSecond, ok := tables.secondStructTable[key]
		if !ok {
			return false
		} else if value != valSecond {
			return false
		}
	}
	return true
}

type ChanTabsStruct chan AnswerTableRequestStruct

type AnswerTableRequestStruct struct {
	answer      map[string]string
	errorAnswer error
}

type StructsTablesDb struct {
	mainStructTable   map[string]string
	secondStructTable map[string]string
}

func StartGettingStructure(conn *pgx.Conn, mess IncomingMess) ChanTabsStruct {
	ch := make(ChanTabsStruct)
	err := conn.Ping(context.Background())
	if err != nil {
		log.Error(conn.Ping(context.Background()))
	}
	go RequestStructTable(conn, ch, mess.Table)
	return ch
}

func RequestStructTable(conn *pgx.Conn, ch ChanTabsStruct, table string) {
	result := make(map[string]string)
	var rowFirst, rowSecond string
	answer := AnswerTableRequestStruct{}
	readyQuery := fmt.Sprintf(QueryTableStruct, table)
	rows, err := conn.Query(context.Background(), readyQuery)
	if err != nil {
		answer.errorAnswer = err
		ch <- answer
		return
	}
	for rows.Next() {
		err = rows.Scan(&rowFirst, &rowSecond)
		if err != nil {
			answer.errorAnswer = err
			ch <- answer
			return
		}
		result[rowFirst] = rowSecond
	}
	answer.answer = result
	ch <- answer
	return
}
