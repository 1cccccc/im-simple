package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"im/commons"
	"im/config"
	"im/models"
	req "im/req"
	resp "im/resp"
	"im/utils"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	var userLoginReq req.UserLoginReq
	if err := c.ShouldBindBodyWithJSON(&userLoginReq); err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))

		return
	}

	//参数校验
	if userLoginReq.Username == "" || userLoginReq.Password == "" {
		c.JSON(http.StatusOK, commons.Error("username or password is empty"))

		return
	}

	//密码比对
	user, err := models.GetUserByUsernamePassword(userLoginReq.Username, utils.GetMD5(userLoginReq.Password))
	if err != nil {
		c.JSON(http.StatusOK, commons.Error("username or password is error"))
		return
	}

	//生成token
	jwt, err := utils.GenerateJWT(user.Identity, user.Email)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	c.JSON(http.StatusOK, commons.Ok(gin.H{
		commons.TOKEN_HEADER: jwt,
	}))
}

func Detail(c *gin.Context) {
	//获取用户id
	uc, exist := c.Get(commons.USER_CLAIM)
	userClaim := uc.(*utils.UserClaim)
	if !exist {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	//获取用户信息
	user, err := models.GetUserByIdentity(userClaim.Identity)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error("user not found"))
		return
	}

	//返回用户信息
	c.JSON(http.StatusOK, commons.Ok(user))
}

func Register(c *gin.Context) {
	var userRegisterReq req.UserRegisterReq
	err := c.ShouldBindBodyWithJSON(&userRegisterReq)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	//参数校验
	if userRegisterReq.Username == "" || userRegisterReq.Password == "" || userRegisterReq.Email == "" || userRegisterReq.Code == "" ||
		len(userRegisterReq.Username) > commons.USERNAME_PASSWORD_INVAILD_LEN || len(userRegisterReq.Password) > commons.USERNAME_PASSWORD_INVAILD_LEN ||
		!utils.ValidateEmail(userRegisterReq.Email) || len(userRegisterReq.Code) != commons.CODE_NUM {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))
		return
	}

	//验证码校验
	result, err := config.Redis.Get(context.TODO(), commons.TOKEN_PREFIX+userRegisterReq.Email).Result()
	if err != nil || result != userRegisterReq.Code { //验证码错误或其他
		c.JSON(http.StatusOK, commons.Error(commons.CODE_ERROR_MSG))
		return
	}

	//用户是否已经存在
	count, err := models.GetUserCountByUsername(userRegisterReq.Username)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, commons.Error(commons.USER_EXISTS_MSG))
		return
	}

	//创建实体，用户注册
	nowUnix := time.Now().Unix()
	user := &models.User{
		Username: userRegisterReq.Username,
		Password: utils.GetMD5(userRegisterReq.Password),
		Nickname: userRegisterReq.Nickname,
		Sex:      userRegisterReq.Sex,
		Email:    userRegisterReq.Email,
		Avatar:   userRegisterReq.Avatar,
		CreateAt: nowUnix,
		UpdateAt: nowUnix,
	}

	err = models.CreateOneUser(user)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
	}

	//将验证码删除
	go config.Redis.Del(context.TODO(), commons.TOKEN_PREFIX+userRegisterReq.Email).Result()

	return
}

func QueryUser(c *gin.Context) {
	userIdentity := c.Query("user_identity") //目标用户的 identity

	if userIdentity == "" {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))
		return
	}

	//获取用户信息
	userClaim := c.MustGet(commons.USER_CLAIM).(*utils.UserClaim)

	//先查询用户的基础信息
	user, err := models.GetUserByIdentity(userIdentity)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	//将信息拷贝到resp对象中
	queryUserResp := resp.QueryUserResp{
		Identity: user.Identity,
		Username: user.Username,
		Nickname: user.Nickname,
		Sex:      user.Sex,
		Email:    user.Email,
		Avatar:   user.Avatar,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
		IsFriend: false,
	}

	//查看是否是好友
	//查询自己的所有单聊房间
	queryUserResp.IsFriend = UserIsFriend(userClaim.Identity, userIdentity)

	c.JSON(http.StatusOK, commons.Ok(queryUserResp))
}

func UserIsFriend(userIdentityA, userIdentityB string) bool {
	userRooms, err := models.GetUserRoomsByUserIdentityRoomType(userIdentityA, 1)
	if err != nil {
		return false
	}

	rooms := make([]string, 0)
	for _, userRoom := range userRooms {
		rooms = append(rooms, userRoom.RoomIdentity)
	}

	//查看目标用户是否在自己的单聊房间中
	count, err := models.GetUserRoomCountByUserIdentityRoomIdentity(userIdentityB, rooms, 1)
	if err != nil {
		return false
	}

	if count <= 0 {
		return false
	}

	return true
}

func GetRoomIdentity(userIdentityA, userIdentityB string) string {
	userRooms, err := models.GetUserRoomsByUserIdentityRoomType(userIdentityA, 1)
	if err != nil {
		return ""
	}

	rooms := make([]string, 0)
	for _, userRoom := range userRooms {
		rooms = append(rooms, userRoom.RoomIdentity)
	}

	//查看目标用户是否在自己的单聊房间中
	roomIdentity, err := models.GetRoomIdentityCountByUserIdentityRoomIdentity(userIdentityB, rooms, 1)
	if err != nil {
		return ""
	}

	return roomIdentity
}

func AddFriend(c *gin.Context) {
	m := make(map[string]interface{})
	c.ShouldBindBodyWithJSON(&m)

	userIdentity := m["user_identity"].(string)

	if userIdentity == "" {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))
		return
	}

	//获取用户信息
	userClaim := c.MustGet(commons.USER_CLAIM).(*utils.UserClaim)

	//查看是否已经是好友
	isFriend := UserIsFriend(userClaim.Identity, userIdentity)
	if isFriend {
		c.JSON(http.StatusOK, commons.Error(commons.ALREADY_FRIEND_MSG))
		return
	}

	//创建单聊房间
	nowUnix := time.Now().Unix()
	roomIdentity := uuid.New().String()
	room := &models.Room{
		ID:           roomIdentity,
		UserIdentity: userClaim.Identity,
		CreateAt:     nowUnix,
		UpdateAt:     nowUnix,
	}

	err := models.CreateRoom(room)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	//创建双方的好友关系
	myur := &models.UserRoom{
		UserIdentity: userClaim.Identity,
		RoomIdentity: roomIdentity,
		RoomType:     1,
		CreateAt:     nowUnix,
		UpdateAt:     nowUnix,
	}

	youur := &models.UserRoom{
		UserIdentity: userIdentity,
		RoomIdentity: roomIdentity,
		RoomType:     1,
		CreateAt:     nowUnix,
		UpdateAt:     nowUnix,
	}

	err = models.CreateUserRoom(myur)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	err = models.CreateUserRoom(youur)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	c.JSON(http.StatusOK, commons.Ok("add friend success"))
}

func DelFriend(c *gin.Context) {
	m := make(map[string]interface{})
	c.ShouldBindBodyWithJSON(&m)

	userIdentity := m["user_identity"].(string)

	if userIdentity == "" {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))
		return
	}

	//获取用户信息
	userClaim := c.MustGet(commons.USER_CLAIM).(*utils.UserClaim)

	//查看是否已经是好友
	isFriend := UserIsFriend(userClaim.Identity, userIdentity)
	if !isFriend {
		c.JSON(http.StatusOK, commons.Error(commons.NOT_FRIEND_MSG))
		return
	}

	//查询两者关联的好友关系
	roomIdentity := GetRoomIdentity(userClaim.Identity, userIdentity)
	if roomIdentity == "" {
		c.JSON(http.StatusOK, commons.Error(commons.NOT_FRIEND_MSG))
		return
	}

	//删除好友单聊房间
	err := models.DeleteRoom(roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	//删除双方的好友关系
	err = models.DeleteUserRoomByUserIdentityAndRoomIdentity(userClaim.Identity, roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}
	models.DeleteUserRoomByUserIdentityAndRoomIdentity(userIdentity, roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	c.JSON(http.StatusOK, commons.Ok("del friend success"))
}
