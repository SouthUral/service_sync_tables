package postgres

import (
	"context"
	"fmt"
	"strconv"

	pgx "github.com/jackc/pgx/v5"
)

// Функция для запуска и контроля синхронизаций.
// Принимает на вход структуру с коннекторами к базам, сообщение IncomingMess,
// chankRead - нужен для определения размера чанка чтения,
// chankWrite - нужен для определения размера чанка записи.
func sync(connects ConnectsPG, mess IncomingMess, chankRead, chankWrite string, OutgoingChan OutgoingChanSync) {

}

// Канал для возвращения оффсета горутине чтения из горутины обработки
type offsetReturnsCh chan string

// Канал для передачи горутине обработки данных.
// Канал нужно передать горутине чтения и горутине обработки
type incomTransmissinCh chan MessForProcessingData

// Сообщение для горутины обработки данных от горутины чтения
type MessForProcessingData struct {
	Rows     pgx.Rows
	MessInfo string
}

// Канал для передачи данных от горутины обработки для горутины записи.
// Канал нужно передать горутине обработки и горутине записи.
type outgoingTransmissCh chan dataForRecording

// Структура для передачи данных от горутины обработки для горутины записи
type dataForRecording struct {
	Fields []string
	Data   [][]any
}

// Канал для управления горутинами
type controlGorutinCh chan string

// Канал для сообщений из горутин (записи, обработки, чтения) в центральный поток синхронизации.
type responseCh chan responseMessGorutine

// Структура сообщения из горутин (записи, обработки, чтения) в центральный поток синхронизации.
type responseMessGorutine struct {
	InfoGorutine string // информация от какой горутины пришло сообщение ();
	ErrorMess    error  // Ошибка, которую возвращает горутина;
	Offset       string // Offset возвращает только горутина записи;
}

// Структура содержащая каналы для контроля горутин синхронизации.
type controlsChGorutines struct {
	chReadData       controlGorutinCh
	chWriteData      controlGorutinCh
	chProcessingData controlGorutinCh
}

// Функция для старта горутин синхронизации (горутины: readData, writeData, processingData)
func startFuncsSync(mess IncomingMess, connects ConnectsPG, chankRead string) controlsChGorutines {
	// каналы для общения между горутинами синхронизации
	transmissionChForProcessing := make(incomTransmissinCh)
	transmissionChForWriting := make(outgoingTransmissCh)
	responseChGorutines := make(responseCh, 100)
	offsetCh := make(offsetReturnsCh)

	answer := controlsChGorutines{}

	// каналы для контроля горутин извне
	answer.chWriteData = make(controlGorutinCh)
	answer.chReadData = make(controlGorutinCh)
	answer.chProcessingData = make(controlGorutinCh)

	go readData(
		transmissionChForProcessing,
		responseChGorutines,
		offsetCh,
		connects.MainConn,
		mess.Table,
		mess.Schema,
		mess.Offset,
		chankRead,

		answer.chReadData,
	)

	go writeData(
		transmissionChForWriting,
		responseChGorutines,
		connects.SecondConn,
		mess.Table,
		mess.Schema,
		answer.chWriteData,
	)

	go processingData(
		transmissionChForProcessing,
		transmissionChForWriting,
		responseChGorutines,
		answer.chProcessingData,
	)

	return answer
}

// Функция для отправки сообщения ошибки в центральную горутину синхронизации
func sendErrorSync(Info string, ErrorSync error, responseCh responseCh) {
	responseCh <- responseMessGorutine{
		InfoGorutine: Info,
		ErrorMess:    ErrorSync,
	}
}

// Функция запроса к БД оффсету
func doQuery(chToProcessing incomTransmissinCh, responseCh responseCh, conn *pgx.Conn, table, schema, offset, chunk string) error {
	switch offset {
	case First:
		query := fmt.Sprintf("SELECT * FROM %s.%s ORDER BY id limit $1::int;", schema, table)
		rows, err := conn.Query(context.Background(), query, chunk)
		if err != nil {
			sendErrorSync(GorReadData, err, responseCh)
			return err
		}
		chToProcessing <- MessForProcessingData{
			Rows: rows,
		}
	case Last:
		query := fmt.Sprintf("SELECT * FROM %s.%s ORDER BY id DESC limit 1;", schema, table)
		rows, err := conn.Query(context.Background(), query)
		if err != nil {
			sendErrorSync(GorReadData, err, responseCh)
			return err
		}
		chToProcessing <- MessForProcessingData{
			Rows:     rows,
			MessInfo: Last,
		}
	default:
		_, err := strconv.Atoi(offset)
		if err != nil {
			sendErrorSync(GorReadData, err, responseCh)
			return err
		}
		query := fmt.Sprintf("SELECT * FROM %s.%s WHERE id > $1::int ORDER BY id limit $2::int;", schema, table)
		rows, err := conn.Query(context.Background(), query, offset, chunk)
		if err != nil {
			sendErrorSync(GorReadData, err, responseCh)
			return err
		}
		chToProcessing <- MessForProcessingData{
			Rows: rows,
		}
	}
	return nil
}

// Горутина чтения данных из БД1
func readData(chToProcessing incomTransmissinCh, responseCh responseCh, offsetCh offsetReturnsCh, conn *pgx.Conn, table, schema, offset, chunk string, contolCh controlGorutinCh) {

	err := doQuery(chToProcessing, responseCh, conn, table, schema, offset, chunk)
	if err != nil {
		return
	}

	for {
		select {
		case _ = <-contolCh:
			return
		case messOffset := <-offsetCh:
			err := doQuery(chToProcessing, responseCh, conn, table, schema, messOffset, chunk)
			if err != nil {
				return
			}
		}
	}
}

// Горутина записи данных в БД2
func writeData(chIncomData outgoingTransmissCh, responseCh responseCh, conn *pgx.Conn, table, schema string, contolCh controlGorutinCh) {

}

// Горутина обработки данных для их последующей записи
func processingData(chIncomData incomTransmissinCh, chOutgoinData outgoingTransmissCh, responseCh responseCh, contolCh controlGorutinCh) {

}
