package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode       = 0
	ParamErrCode      = 10001
	LoginErrCode            = 10002
	UserNotExistErrCode     = 10003
	UserAlreadyExistErrCode = 10004
	ServerErrCode           = 10005
	TokenInvalidErrCode     = 10006
)

type ErrNo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrNo) Error() string{
	return fmt.Sprintf("err_code=%d err_msg=%s",e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int32, msg string) ErrNo{
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo{
	e.ErrMsg = msg
	return e
}

var (
	Success             = NewErrNo(SuccessCode, "Success")
	ParamErr            = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	LoginErr            = NewErrNo(LoginErrCode, "Wrong username or password")
	UserNotExistErr     = NewErrNo(UserNotExistErrCode, "User does not exists")
	UserAlreadyExistErr = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	TokenInvalidErr     = NewErrNo(TokenInvalidErrCode, "Token is invalid")
)


func ConvertErr(err error) ErrNo{
	if errors.As(err, &ErrNo{}) {
		return err.(ErrNo)
	}

	Err := ErrNo{
		ErrCode: ServerErrCode,
		ErrMsg: err.Error(),
	}
	return Err

}