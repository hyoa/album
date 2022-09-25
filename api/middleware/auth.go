package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/internal/user"
)

type kindAuth string

const (
	BasicAuth kindAuth = "basicAuth"
	AdminAuth kindAuth = "adminAuth"
)

func Auth(userRepo user.UserRepo, kind kindAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header["Authorization"]
		jwt := strings.ReplaceAll(bearer[0], "Bearer ", "")

		aTokenizer := user.CreateAuthTokenizer()
		t, _ := aTokenizer.Decode(jwt)

		if (kind == AdminAuth && t.Role != user.RoleAdmin) || (kind == BasicAuth && t.Role != user.RoleNormal) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		u, errFind := userRepo.FindByEmail(t.Email)

		if u == (user.User{}) || errFind != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}

		c.Set("user", u)
	}
}
