package main

const (
	StatusActive   = 0
	StatusInactive = 1
)

type State struct {
	table  string
	status int
}

// создает структуру State и запускает горутину StateWorker
func InitState(mongoChInput chan MessCommand, mongoChOutput chan MessCommand) {
	w_state := State{}
	go w_state.StateWorker(mongoChInput, mongoChOutput)
}

// основная функция для работы с состоянием, при запуске она стартует и ждет ответа от монго
// состояние должно принимать парметрами 6 каналов, которые создаются в main
// для обработки сообщений с каждым из модулей можно создать функции (3 функции)
func (state *State) StateWorker(mongoChInput chan MessCommand, mongoChOutput chan MessCommand) {
	state.startWork(mongoChInput, mongoChOutput)

}

func (state *State) startWork(mongoChInput chan MessCommand, mongoChOutput chan MessCommand) {
	for {
		startCommand := MessCommand{
			Info: "get_all",
		}
		mongoChInput <- startCommand
	}

}
