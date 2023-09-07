package postgres

import (
	pgx "github.com/jackc/pgx/v5"
)

// Функция для запуска и контроля синхронизаций.
// Принимает на вход структуру с коннекторами к базам, сообщение IncomingMess,
// chankRead - нужен для определения размера чанка чтения,
// chankWrite - нужен для определения размера чанка записи.
func sync(connects ConnectsPG, mess IncomingMess, chankRead, chankWrite string, OutgoingChan OutgoingChanSync) {

}

// Канал для передачи горутине обработки данных
// Канал нужно передать горутине чтения и горутине обработки
type incomTransmissinCh chan pgx.Rows

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
type responseCh chan responseGotutine

// Структура сообщения из горутин (записи, обработки, чтения) в центральный поток синхронизации.
type responseGotutine struct {
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

	answer := controlsChGorutines{}

	// каналы для контроля горутин извне
	answer.chWriteData = make(controlGorutinCh)
	answer.chReadData = make(controlGorutinCh)
	answer.chProcessingData = make(controlGorutinCh)

	go readData(
		transmissionChForProcessing,
		responseChGorutines,
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

// Горутина чтения данных из БД1
func readData(chToProcessing incomTransmissinCh, responseCh responseCh, conn *pgx.Conn, table, schema, offset, chunk string, contolCh controlGorutinCh) {

}

// Горутина записи данных в БД2
func writeData(chIncomData outgoingTransmissCh, responseCh responseCh, conn *pgx.Conn, table, schema string, contolCh controlGorutinCh) {

}

// Горутина обработки данных для их последующей записи
func processingData(chIncomData incomTransmissinCh, chOutgoinData outgoingTransmissCh, responseCh responseCh, contolCh controlGorutinCh) {

}
