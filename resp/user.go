package resp

type QueryUserResp struct {
	Identity string `json:"_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Sex      uint8  `json:"sex"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
	IsFriend bool   `json:"is_friend"`
}
