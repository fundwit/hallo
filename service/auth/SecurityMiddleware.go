package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthenticateByToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			token := auth[7:]
			if securityContext, find := TokenCache.Get(token); find {
				SaveToRequestContext(context, securityContext.(*SecurityContext))
			}
		}
		context.Next()
	}
}
func AuthenticatedCheck() gin.HandlerFunc {
	return func(context *gin.Context) {
		securityContext := LoadFromRequestContext(context)
		if securityContext != nil {
			context.Next()
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "authentication is required"})
		}
	}
}
