package commons

import "time"

const (
	//Redis
	TOKEN_PREFIX      = "token/"
	TOKEN_EXPIRE_TIME = time.Minute * 5

	//Code
	CODE_NUM = 6

	//Default values
	PAGE_DEFAULT_VALUE      = 1
	PAGE_SIZE_DEFAULT_VALUE = 10

	//User
	USER_CLAIM = "userClaim"

	//Request Header
	TOKEN_HEADER = "token"

	//Invaild Value
	USERNAME_PASSWORD_INVAILD_LEN = 20
)
