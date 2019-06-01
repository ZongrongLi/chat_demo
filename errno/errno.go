package errno

var (
	// Common errors
	OK               = &Errno{Code: 0, Message: "OK"}
	ErrUserNotExist  = &Errno{Code: 1, Message: "user not exist"}
	ErrInvalidPasswd = &Errno{Code: 2, Message: "Passwd or username not right"}
	ErrInvalidParams = &Errno{Code: 3, Message: "Invalid params"}
	ErrUserExist     = &Errno{Code: 4, Message: "user exist"}
)

type Errno struct {
	Code    int
	Message string
}

func (e Errno) Error() string {
	return e.Message
}
