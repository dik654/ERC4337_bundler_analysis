package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dik654/Go_projects/SNS_SERVER/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func EnsureLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		googleSessionID := session.Get("google_user")
		regularSessionID := session.Get("regular_user")
		if (googleSessionID == nil && regularSessionID == nil) || !validateCookieSignature(ctx, os.Getenv("SECRET_KEY")) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "unauthorized"})
			return
		}
		session.Options(sessions.Options{
			MaxAge: int(30 * time.Minute),
		})
		session.Save()

		ctx.Next()
	}
}

func validateCookieSignature(c *gin.Context, secretKey string) bool {
	cookieValue, err := c.Cookie("session_cookie")
	if err != nil {
		return false
	}

	// 쿠키 값과 서명 분리
	parts := strings.Split(cookieValue, "|")
	if len(parts) != 2 {
		return false
	}

	value, signature := parts[0], parts[1]

	// 서명 검증
	expectedSignature := utils.CreateSignature(value, secretKey)
	return signature == expectedSignature
}
