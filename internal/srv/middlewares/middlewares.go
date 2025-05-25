package middlewares

import (
	"gophermart-points/internal/datacrypt"
	"gophermart-points/internal/srv/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golodash/godash/strings"
)

func RestrictText() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.StartsWith(c.ContentType(), "text/plain") {
			c.JSON(http.StatusBadRequest, handlers.RsDef{Err: "Unsupported Content-Type"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RestrictJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.StartsWith(c.ContentType(), "application/json") {
			c.JSON(http.StatusBadRequest, handlers.RsDef{Err: "Unsupported Content-Type"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckAuth(authKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		unparsed, err := c.Cookie(handlers.USER_COOKIE_NAME)
		if err != nil {
			c.JSON(http.StatusUnauthorized, handlers.RsDef{
				Err: "Unauthorized",
			})
			c.Abort()
			return
		}

		userId, err := datacrypt.GetUserID(unparsed, authKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, handlers.RsDef{
				Err: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set(handlers.USER_ID_KEY, userId)

		c.Next()
	}
}
