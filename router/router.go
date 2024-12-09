package router

import (
	"github.com/gin-gonic/gin"
	"im/middleware"
	"im/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/sendCode", service.SendCode)

	//用户路由
	authRouter := r.Group("/", middleware.AuthCheck())
	notAuthRouter := r.Group("/")

	//根路径需要验证token
	chatAuthRouter := authRouter.Group("/chat")
	chatAuthRouter.GET("/ws", service.WebSocketConnect)
	chatAuthRouter.GET("/list", service.GetChatList)

	//用户路径不需要验证token
	userNotAuthRouter := notAuthRouter.Group("/user")
	userNotAuthRouter.POST("/login", service.Login)
	userNotAuthRouter.PUT("/register", service.Register)

	//用户路径需要验证token
	userAuthRouter := authRouter.Group("/user")
	userAuthRouter.GET("/detail", service.Detail)
	userAuthRouter.GET("/query", service.QueryUser)
	userAuthRouter.PUT("/addFriend", service.AddFriend)
	userAuthRouter.DELETE("/delFriend", service.DelFriend)

	return r
}
