package serveHttp

import "github.com/gin-gonic/gin"

type Routable interface {
	RegisterRoutes(r *gin.RouterGroup)
}
