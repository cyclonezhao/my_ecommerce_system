package errorhandler

type BusinessError struct {
	Code int
	Message string
	HttpCode int
}

func (e *BusinessError) Error() string {
	return e.Message
}