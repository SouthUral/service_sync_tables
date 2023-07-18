package main

const (
	StatusActive   = 0
	StatusInactive = 1
)

type State struct {
	table     string
	status    int
	mdbInput  chan MessCommand
	mdbOutput chan MessCommand
}

// создает структуру State и запускает горутину StateWorker
func InitState(mongoChInput chan MessCommand, mongoChOutput chan MessCommand) {
	w_state := State{
		mdbInput:  mongoChInput,
		mdbOutput: mongoChOutput,
	}
	go w_state.StateWorker()
}

// основная функция для работы с состоянием, при запуске она стартует и ждет ответа от монго
// состояние должно принимать парметрами 6 каналов, которые создаются в main
// для обработки сообщений с каждым из модулей можно создать функции (3 функции)
func (state *State) StateWorker() {
	// метод для старта, запрашивает из монго все документы
	state.mdbInput <- MessCommand{Info: "get_all"}

	// запускается бесконечный цикл обработки сообщений
	for {
		select {
		case mess := <-state.mdbOutput:
			state.MongoWorker(mess)
		}
	}

}

// Обработчик сообщений приходящих от модуля MongoDB
func (state *State) MongoWorker(mess MessCommand) {

}
