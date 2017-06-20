package system

type Error interface {
	Error() string
	Map() Map
	ResultMap(args ...string) Map
}

type QError struct {
	Code    int
	Message string
}

func (e *QError) Error() string {
	return e.Message
}

func (e *QError) Map() Map {
	return Map{"code": e.Code, "message": e.Message}
}

func (e *QError) ResultMap(args ...string) Map {
	if args != nil {
		e.Message = args[0]

	}
	return Map{"status": e.Map()}
}

func ErrorResult(msg string, code int) Map {
	e := QError{code, msg}
	return Map{"status": e.Map()}
}

func NewQError(code int, msg string) *QError {
	return &QError{Code: code, Message: msg}
}

func SqlQError(err error) *QError {
	return &QError{Code: SqlErrorID, Message: err.Error()}
}

func TokenQError(err error) *QError {
	return &QError{Code: TokenErrorID, Message: err.Error()}
}
