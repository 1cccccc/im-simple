package req

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Sex      uint8  `json:"sex"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Code     string `json:"code"`
}
