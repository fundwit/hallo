package meta

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(r *gin.RouterGroup) {
	r.GET("", info)
}

func info(c *gin.Context) {
	c.JSON(http.StatusOK, GetServiceInfo())
}
