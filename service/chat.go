package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im/commons"
	"im/models"
	"im/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{}

var conns = make(map[string]*websocket.Conn)

type WsMessage struct {
	RoomIdentity string `json:"room_identity"`
	Message      string `json:"message"`
}

func GetChatList(c *gin.Context) {
	room_identity := c.Query("room_identity")
	p, err := strconv.ParseInt(c.Query("p"), 10, 64) //page
	s, err := strconv.ParseInt(c.Query("s"), 10, 64) //size

	if room_identity == "" {
		c.JSON(http.StatusOK, commons.Error(commons.INVALID_ARGS_MSG))
		return
	}

	if p == 0 || s == 0 {
		p = commons.PAGE_DEFAULT_VALUE
		s = commons.PAGE_SIZE_DEFAULT_VALUE
	}

	//拿到用户信息
	us := c.MustGet(commons.USER_CLAIM).(*utils.UserClaim)

	// 查看用户是否在该房间内
	_, err = models.GetUserRoomByUserIdAndRoomId(us.Identity, room_identity)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.USER_NOT_EXISTS_ROOM_MSG))
		return
	}

	//返回聊天列表
	ml, err := models.FindMessageListByRoomIdentity(room_identity, s, p)
	if err != nil {
		c.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	c.JSON(http.StatusOK, commons.Ok(ml))
}

func WebSocketConnect(ctx *gin.Context) {
	//协议升级
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, commons.Error(commons.SYSTEM_ERROR_MSG))
		return
	}

	defer conn.Close()

	//拿到用户信息
	us := ctx.MustGet(commons.USER_CLAIM).(*utils.UserClaim)

	//是否已经连接,一个用户只能有一个连接
	if _, ok := conns[us.Identity]; ok {
		//剔除原来的连接
		oldConn := conns[us.Identity]

		//关闭原来的连接
		oldConn.Close()

	}

	//保存连接
	conns[us.Identity] = conn

	defer delete(conns, us.Identity) //删除连接

	for {
		if conn != conns[us.Identity] {
			//连接已经被关闭
			log.Printf("conn has been closed\n")
			return
		}

		//接收消息
		var msg WsMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("conn.ReadJSON err:%v\n", err)
			return
		}

		// 查看用户是否在该房间内
		_, err = models.GetUserRoomByUserIdAndRoomId(us.Identity, msg.RoomIdentity)
		if err != nil {
			log.Printf("The user has not joined the room:%v\n", err)
			return
		}

		//TODO 保存消息到数据库
		nowUnix := time.Now().Unix()

		message := &models.Message{
			UserIdentity: us.Identity,
			RoomIdentity: msg.RoomIdentity,
			Data:         msg.Message,
			CreateAt:     nowUnix,
			UpdateAt:     nowUnix,
		}

		models.InsertOneMessage(message)

		//TODO 查询房间内在线的用户
		urs, err := models.GetUserRoomsByRoomIdentity(msg.RoomIdentity)
		if err != nil {
			log.Printf("models.GetUserRoomsByRoomIdentity err:%v\n", err)
			return
		}

		// 给房间内在线的用户发送消息
		for _, ur := range urs {
			if cc, ok := conns[ur.UserIdentity]; ok {
				if cc == conn { //跳过自己
					continue
				}

				err := cc.WriteMessage(websocket.TextMessage, []byte(msg.Message))
				if err != nil {
					log.Printf("WriteMessage err:%v\n", err)
				}
			}
		}

		//发送消息
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg.Message))
		if err != nil {
			log.Printf("conn.WriteMessage err:%v\n", err)
			return
		}
	}
}
