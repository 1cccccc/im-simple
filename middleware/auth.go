package middleware

import (
	"github.com/gin-gonic/gin"
	"im/commons"
	"im/utils"
	"net/http"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(commons.TOKEN_HEADER)

		if token == "" {
			c.Abort()
			c.JSON(http.StatusOK, commons.Error(commons.NOT_AUTH_MSG))
			return
		}

		userClaim, err := utils.ParseJWT(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, commons.Error(commons.NOT_AUTH_MSG))
			return
		}

		c.Set(commons.USER_CLAIM, userClaim)

		c.Next()
	}
}
