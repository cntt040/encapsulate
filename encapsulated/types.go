package encapsulated

import "fmt"

const (
	CodeInternal = "500"
	CodeNetWork  = "598"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetError() error {

	if e.Message != "" {
		return Error{
			Message: e.Message,
			Code:    "503",
		}
	}
	return nil
}

func WrapError(err error, code string, msg string, args ...interface{}) error {
	if err, ok := err.(Error); ok {
		return err
	}
	return Error{
		Code:    code,
		Message: fmt.Sprintf(msg, args...),
	}
}

type Response interface {
	GetError() error
}
