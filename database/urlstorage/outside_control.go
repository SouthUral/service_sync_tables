package urlstorage

// Функция для отправки сообщения в urlstorage и получения ответа
func SendingMess(mess InputMessage, outputCh InputUrlStorageAPIch) (any, any) {
	reverseCh := make(ReverseAPIch, 0)
	outgoingMess := UrlMessInput{
		Message:   mess,
		ReverseCh: reverseCh,
	}
	outputCh <- outgoingMess
	answer, _ := <-reverseCh
	if answer.Error != nil {
		return nil, answer.Error
	}
	return answer.AnswerData, nil
}
