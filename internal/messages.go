package internal

var errorMessages = make([]error, 0)

func AddError(err error) {
	if err == nil {
		return
	}
	errorMessages = append(errorMessages, err)
}

func GetErrors() []error {
	return errorMessages
}
func FlushErros() {
	errorMessages = make([]error, 0)
}
