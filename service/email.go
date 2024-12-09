package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"im/commons"
	"im/config"
	"im/req"
	"im/utils"
	"net/http"
)

func SendCode(c *gin.Context) {
	var sendCodeReq req.SendCodeReq

	err := c.ShouldBindBodyWithJSON(&sendCodeReq)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))
		return
	}

	//查看邮箱是否已经发送验证码
	oldCode := config.Redis.Get(context.TODO(), commons.TOKEN_PREFIX+sendCodeReq.To)
	if oldCode.Err() == nil && oldCode.Val() != "" { //已经存在验证码
		c.JSON(http.StatusOK, commons.Error(commons.CODE_EXISTS_MSG))
		return
	}

	//生成验证码
	code := utils.GetNumCode(commons.CODE_NUM)

	//发送邮件,异步发送
	go utils.SendCode(code, sendCodeReq.To)

	//保存验证码到redis
	config.Redis.Set(context.TODO(), commons.TOKEN_PREFIX+sendCodeReq.To, code, commons.TOKEN_EXPIRE_TIME)

	c.JSON(http.StatusOK, commons.Ok("验证码已发送成功，请注意查收"))
}
