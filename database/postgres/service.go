package postgres

import (
	pgx "github.com/jackc/pgx/v5"
)

func CheckingStrucTables(connects ConnectsPG) error {
	chMainDb := StartGettingStructure(connects.MainConn)
	chSeconDb := StartGettingStructure(connects.SecondConn)
}

type ChanTabsStruct chan map[string]string

func StartGettingStructure(conn *pgx.Conn) ChanTabsStruct {
	ch := make(ChanTabsStruct)
	go RequestStructTable(conn, ch)
	return ch
}

func RequestStructTable(conn *pgx.Conn, ch ChanTabsStruct) {

}
