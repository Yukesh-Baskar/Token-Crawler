package errorshandler

type NewError struct {
	Error  interface{}
	Status int
}

func (e *NewError) HandleError() *NewError {
	return &NewError{
		Error:  e.Error,
		Status: e.Status,
	}
}

func GetErrors(e error, status int) *NewError {
	return &NewError{
		Error:  e,
		Status: status,
	}
}
