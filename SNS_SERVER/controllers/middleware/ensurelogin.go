package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func EnsureLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionID := session.Get("user")
		if sessionID == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "unauthorized"})
			return
		}
		session.Options(sessions.Options{
			MaxAge: int(30 * time.Minute),
		})
		session.Save()

		ctx.Next() //
	}
}
