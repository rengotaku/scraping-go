package controllers

type FormError struct {
	messages []string
}

func (fe FormError) addMessage(message string) {
	if fe.messages == nil {
		fe.messages = []string{}
	}

	fe.messages = append(fe.messages, message)
}

func (fe FormError) hasMessage(key string) bool {
	return fe.messages == nil || len(fe.messages) > 0
}
