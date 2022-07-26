package config

import "douyin-proj/src/types"

const (
	Success              int32 = 0 // 成功响应
	ParamInvalid         int32 = 1 // 请求参数错误
	UserHasExisted       int32 = 2 // 用户名已经存在
	UserNotExisted       int32 = 3 // 用户名不存在
	WrongPassword        int32 = 4 // 密码错误
	AuthFailed           int32 = 5 // token校验失败
	VideoUploadFailed    int32 = 6 // 视频上传失败
	DuplicateFavorite    int32 = 7
	NotInFavorite        int32 = 8
	CommentAddFailed     int32 = 9
	CommentDeleteFailed  int32 = 10
	RelationAddFailed    int32 = 12
	RelationDeleteFailed int32 = 13
	VideoNotExisted      int32 = 14 // 视频不存在
	UnknownError         int32 = 255
)

var (
	SuccessResp             = types.Response{StatusCode: Success, StatusMsg: "Success..."}
	ParamInvalidResp        = types.Response{StatusCode: ParamInvalid, StatusMsg: "Parameters are Invalid!"}
	UserHasExistedResp      = types.Response{StatusCode: UserHasExisted, StatusMsg: "User has been Registered!"}
	UserNotExistedResp      = types.Response{StatusCode: UserNotExisted, StatusMsg: "User does not exist!"}
	WrongPasswordResp       = types.Response{StatusCode: WrongPassword, StatusMsg: "Password is wrong!"}
	DuplicateFavoriteResp   = types.Response{StatusCode: DuplicateFavorite, StatusMsg: "Duplicate favorite!"}
	NotInFavoriteResp       = types.Response{StatusCode: NotInFavorite, StatusMsg: "Not in favorite list!"}
	UnknownErrorResp        = types.Response{StatusCode: UnknownError, StatusMsg: "Unknown error!"}
	AuthFailedResp          = types.Response{StatusCode: AuthFailed, StatusMsg: "Token is invalid!"}
	VideoUploadFailedResp   = types.Response{StatusCode: VideoUploadFailed, StatusMsg: "Video upload failed!"}
	CommentDeleteFailedResp = types.Response{StatusCode: CommentDeleteFailed, StatusMsg: "Comment delete failed!"}
	CommentAddFailedResp    = types.Response{StatusCode: CommentAddFailed, StatusMsg: "Comment add failed!"}
	VideoNotExistedResp     = types.Response{StatusCode: VideoNotExisted, StatusMsg: "video does not exist!"}
)
