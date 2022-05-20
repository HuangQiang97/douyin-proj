package types

type RelationResponse Response

type UserFollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type UserFansListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
