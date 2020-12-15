package auth

import "github.com/gin-gonic/gin"

type SecurityContext struct {
	Token     string
	Principal Principal
}

const securityContextKey = "SECURITY_CONTEXT"

func LoadFromRequestContext(c *gin.Context) *SecurityContext {
	if sc, find := c.Get(securityContextKey); find {
		return sc.(*SecurityContext)
	}
	return nil
}

func SaveToRequestContext(c *gin.Context, securityContext *SecurityContext) {
	c.Set(securityContextKey, securityContext)
}
