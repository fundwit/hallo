package serveHttp

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"hallo/domain"
	"hallo/service/auth"
	"net/http"
)

type RegistryHandler struct {
	AccountRepository domain.AccountRepository
}

type EmailOccupiedQuery struct {
	Email string `json:"email"  binding:"required"`
}

type EmailOccupiedInfo struct {
	Email    string `json:"email"  binding:"required"`
	Occupied bool   `json:"occupied" binding:"required"`
}

type UsernameOccupiedQuery struct {
	Name string `json:"name"  binding:"required"`
}

type UsernameOccupiedInfo struct {
	Name     string `json:"name"  binding:"required"`
	Occupied bool   `json:"occupied"  binding:"required"`
}

func (handler *RegistryHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/emails", handler.emailOccupied)
	r.POST("/names", handler.usernameOccupied)
	r.POST("/email_register_tokens", handler.acquireEmailRegisterToken)
}

func (handler *RegistryHandler) emailOccupied(c *gin.Context) {
	var query EmailOccupiedQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOccupied, err := handler.AccountRepository.IsEmailOccupied(query.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &EmailOccupiedInfo{Email: query.Email, Occupied: isOccupied})
}

func (handler *RegistryHandler) acquireEmailRegisterToken(c *gin.Context) {
	var query EmailOccupiedQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOccupied, err := handler.AccountRepository.IsEmailOccupied(query.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if isOccupied {
		c.JSON(http.StatusConflict, gin.H{"error": ""})
		return
	}

	token, found := auth.RegisterTokenCache.Get(query.Email)
	if !found {
		token = uuid.New().String()
		auth.RegisterTokenCache.Set(query.Email, token, cache.DefaultExpiration)
	}

	// send token to
	// ...

	c.JSON(http.StatusOK, gin.H{"email": query.Email})
}

func (handler *RegistryHandler) usernameOccupied(c *gin.Context) {
	var query UsernameOccupiedQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOccupied, err := handler.AccountRepository.IsAccountNameOccupied(query.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &UsernameOccupiedInfo{Name: query.Name, Occupied: isOccupied})
}
