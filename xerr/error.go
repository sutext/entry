package xerr

import "fmt"

type Error int

const (
	UserNotFound   Error = 10001
	InvalidToken   Error = 10002
	InvalidRequest Error = 10003
	InternalError  Error = 10004
)

func (e Error) Error() string {
	switch e {
	case UserNotFound:
		return "user not found"
	case InvalidToken:
		return "invalid token"
	case InvalidRequest:
		return "invalid request"
	case InternalError:
		return "internal error"
	default:
		return fmt.Sprintf("unknown error %d", e)
	}
}
