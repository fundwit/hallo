package serveHttp

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"hallo/domain"
	"hallo/service/auth"
	"log"
	"net/http"
)

type SessionHandler struct {
	AccountManager domain.AccountManager
}

type LoginRequest struct {
	Name   string `json:"name"   binding:"required" pact:"example=sally"`
	Secret string `json:"secret" binding:"required" pact:"example=secret"`
}

func (handler *SessionHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("", handler.newSession)
	r.DELETE("", auth.AuthenticateByToken(), deleteSession)
	r.GET("/me", auth.AuthenticateByToken(), auth.AuthenticatedCheck(), currentSession)
}

// authentication
func (handler *SessionHandler) newSession(c *gin.Context) {
	var login LoginRequest
	// json.SyntaxError, validate error
	if paramErr := c.ShouldBindJSON(&login); paramErr != nil {
		log.Println(paramErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request body"})
		return
	}

	_, err := handler.AccountManager.AuthenticateInternalIdentity(login.Name, login.Secret)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "account not exist or secret is not match"})
		return
	}

	token := uuid.NewV4().String()

	sc := &auth.SecurityContext{Token: token, Principal: auth.Principal{Name: login.Name}}
	auth.TokenCache.Set(token, sc, cache.DefaultExpiration)
	auth.SaveToRequestContext(c, sc)

	c.Header("Authentication", token)
	c.JSON(http.StatusOK, gin.H{"token": sc.Token, "principal": gin.H{"name": sc.Principal.Name}})
}

func deleteSession(c *gin.Context) {
	securityContext := auth.LoadFromRequestContext(c)
	if securityContext != nil {
		auth.TokenCache.Delete(securityContext.Token)
	}
	c.Status(http.StatusNoContent)
}

func currentSession(c *gin.Context) {
	sc := auth.LoadFromRequestContext(c)
	if sc != nil {
		c.JSON(http.StatusOK, gin.H{"token": sc.Token, "principal": gin.H{"name": sc.Principal.Name}})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": (&domain.ErrUnauthorized{}).Error()})
	}
}
