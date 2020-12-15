package serveHttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"hallo/domain"
	"hallo/domain/entity"
	"hallo/service/auth"
	"log"
	"net/http"
)

type AccountHandler struct {
	AccountManager    domain.AccountManager
	AccountRepository domain.AccountRepository
}

type AccountCreateForm struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Secret        string `json:"secret"   binding:"required"`
	RegisterToken string `json:"register_token" binding:"required"`
}

func (handler *AccountHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("", handler.createAccount)
}

func (handler *AccountHandler) createAccount(c *gin.Context) {
	var form AccountCreateForm
	// with validating
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request body"})
		return
	}

	// clean
	token, found := auth.RegisterTokenCache.Get(form.Email)
	if !found || token != form.RegisterToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": (&domain.ErrRegisterTokenInvalid{}).Error()})
		return
	}
	auth.RegisterTokenCache.Delete(form.Email)

	account, err := handler.AccountManager.CreateAccount(entity.EmailAccountCreateRequest{
		Name: form.Name, Email: form.Email, Secret: form.Secret})
	if err != nil {
		log.Printf("error: %v\n", err)

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request body"})
			return
		} else if errors.Is(err, &domain.AccountNameIsOccupied{}) {
			c.JSON(http.StatusConflict, gin.H{"error": (&domain.AccountNameIsOccupied{}).Error()})
			return
		} else if errors.Is(err, &domain.AccountEmailIsOccupied{}) {
			c.JSON(http.StatusConflict, gin.H{"error": (&domain.AccountEmailIsOccupied{}).Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
	}

	c.JSON(http.StatusCreated, gin.H{"user": account})
}
