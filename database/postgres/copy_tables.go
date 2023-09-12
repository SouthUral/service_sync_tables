package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pgx "github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
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
	Fields     []string
	Data       [][]any
	LastOffset string
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
		offsetCh,
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

	oldOffset := offset
	waitingTime := 0

	err := doQuery(
		chToProcessing,
		responseCh,
		conn,
		table,
		schema,
		offset,
		chunk)
	if err != nil {
		log.Debug("Работа горутины чтения завершена из-за ошибки")
		return
	}

	for {
		select {
		case _ = <-contolCh:
			log.Debug("Работа горутины чтения завершена по команде")
			return
		case messOffset := <-offsetCh:
			if oldOffset == messOffset && waitingTime <= 10 {
				waitingTime++
				time.Sleep(time.Duration(waitingTime) * time.Second)
			} else {
				oldOffset = messOffset
				waitingTime = 0
			}
			err := doQuery(
				chToProcessing,
				responseCh,
				conn,
				table,
				schema,
				messOffset,
				chunk)
			if err != nil {
				log.Debug("Работа горутины чтения завершена из-за ошибки")
				return
			}
		}
	}
}

// Горутина записи данных в БД2.
// Горутина должна быть блокируемая, т.е. запись в БД должна происходить только после команды Continue из канала contolCh
func writeData(chIncomData outgoingTransmissCh, responseCh responseCh, conn *pgx.Conn, table, schema string, contolCh controlGorutinCh) {
	control := Continue
	for {
		select {
		case messControl := <-contolCh:
			switch messControl {
			case Stop:
				log.Debug("Работа горутины записи завершена по команде")
				return
			case Continue:
				control = Continue
			}
		case messData := <-chIncomData:
			if control == Continue {
				err := writer(messData, conn, table, schema)
				if err != nil {
					sendErrorSync(GorWriteData, err, responseCh)
					log.Debug("Работа горутины записи завершена из-за ошибки")
					return
				}
				responseCh <- responseMessGorutine{
					InfoGorutine: GorWriteData,
					Offset:       messData.LastOffset,
				}
				control = Waiting
			} else {
				// Момент, когда сообщение с новыми данными пришло, но ответа от центрального потока еще нет
				messControl := <-contolCh
				switch messControl {
				case Stop:
					log.Debug("Работа горутины записи завершена по команде")
					return
				case Continue:
					err := writer(messData, conn, table, schema)
					if err != nil {
						sendErrorSync(GorWriteData, err, responseCh)
						log.Debug("Работа горутины записи завершена из-за ошибки")
						return
					}
					responseCh <- responseMessGorutine{
						InfoGorutine: GorWriteData,
						Offset:       messData.LastOffset,
					}
				}
			}
		}
	}
}

// Функция записи данных в таблицу
func writer(mess dataForRecording, conn *pgx.Conn, table, schema string) error {
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{schema, table},
		mess.Fields,
		pgx.CopyFromRows(mess.Data),
	)
	if err != nil {
		log.Error("Ошибка записи в БД", err.Error())
		return err
	}
	return nil
}

// Горутина обработки данных для их последующей записи
func processingData(chIncomData incomTransmissinCh, chOutgoinData outgoingTransmissCh, responseCh responseCh, contolCh controlGorutinCh, offsetCh offsetReturnsCh) {
	oldOffset := ""

	for {
		select {
		case messData := <-chIncomData:
			switch messData.MessInfo {
			case Last:
				// выдление оффсета и отправка в модуль чтения
				rows, err := rowsToMap(messData.Rows)
				if err != nil {
					sendErrorSync("processingData", err, responseCh)
					log.Debug("Работа горутины обработки завершена из-за ошибки")
					return
				}
				if len(rows) == 0 {
					offsetCh <- Last
					continue
				} else {
					lastOffset := getLastOffset(rows)
					offsetCh <- lastOffset
				}
			default:
				// обработанные данные отправить в модуль записи вместе с оффсетом
				// оффсет отправить в модуль чтения
				rows, err := rowsToMap(messData.Rows)
				if err != nil {
					sendErrorSync("processingData", err, responseCh)
					log.Debug("Работа горутины обработки завершена из-за ошибки")
					return
				}
				// Если новых записей не было, то горутине чтения отправяляется прошлый оффсет
				if len(rows) == 0 {
					offsetCh <- oldOffset
					continue
				} else {
					oldOffset = getLastOffset(rows)
					offsetCh <- oldOffset
					fields := getFiled(rows)
					resRows := dictionaryConverter(rows, fields)
					chOutgoinData <- dataForRecording{
						Fields:     fields,
						Data:       resRows,
						LastOffset: oldOffset,
					}
				}
			}
		case _ = <-contolCh:
			// при получении сообщения из управляющего канала завершить горутину
			log.Debug("Работа горутины обработки завершена по команде")
			return
		}
	}
}

// Словарь с данными сконвертированными из pgx.Rows
type rowTable map[string]any

// Функция читает данные из rows, конвертирует их в map и записывает в слайс
func rowsToMap(rows pgx.Rows) ([]rowTable, error) {
	res := make([]rowTable, 0)
	for rows.Next() {
		mapRes, err := pgx.RowToMap(rows)
		if err != nil {
			// Нужно отправить ошибку вверх по стеку вызовов
			log.Error(err.Error())
			return res, err
		}
		res = append(res, mapRes)
	}
	return res, nil
}

// Функция получает []rowTable и возвращает последний оффсет в формате строки
func getLastOffset(rows []rowTable) string {
	lastId := len(rows) - 1
	lastItem := rows[lastId]
	return fmt.Sprintf("%d", lastItem["id"])
}

// Функция генерирует слайс с именами полей, слайс с полями необходим для записи данных в таблицу
func getFiled(data []rowTable) []string {
	resSlice := []string{}
	for key := range data[0] {
		resSlice = append(resSlice, key)
	}
	return resSlice
}

// Функция генерирует слайс слайсов из []rowTable, который необходим для записи в таблицу
func dictionaryConverter(data []rowTable, sliceField []string) [][]any {
	res := make([][]any, 0)
	for _, row := range data {
		rowSlice := make([]any, 0)
		for _, key := range sliceField {
			rowSlice = append(rowSlice, row[key])
		}
		res = append(res, rowSlice)
	}
	return res
}
