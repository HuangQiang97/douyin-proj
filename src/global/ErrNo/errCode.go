package ErrNo

import "douyin-proj/src/types"

const (
	Success        int32 = 0
	ParamInvalid   int32 = 1
	UserHasExisted int32 = 2
	UserNotExisted int32 = 3
	WrongPassword  int32 = 4
	UnknownError   int32 = 255
	AuthFailed     int32 = 255
)

var (
	SuccessResp        = types.Response{StatusCode: Success, StatusMsg: "Success"}
	ParamInvalidResp   = types.Response{StatusCode: ParamInvalid, StatusMsg: "Parameters are Invalid!"}
	UserHasExistedResp = types.Response{StatusCode: UserHasExisted, StatusMsg: "User has been Registered!"}
	UserNotExistedResp = types.Response{StatusCode: UserNotExisted, StatusMsg: "User does not exist!"}
	WrongPasswordResp  = types.Response{StatusCode: WrongPassword, StatusMsg: "Password is wrong!"}
	UnknownErrorResp   = types.Response{StatusCode: UnknownError, StatusMsg: "Unknown error!"}
)
