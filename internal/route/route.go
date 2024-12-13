package route

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"meta-node-ficam/internal/handler"
)

func InitialRoutes(engine *gin.Engine, handler handler.EmailHandler) {
	r := engine.Group("/api/v1",PreflightHandler())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	nameRoute := r.Group("/email")
	{
		nameRoute.POST("/verification", handler.EmailVerification)
		nameRoute.POST("/authentication", handler.EmailAuthentication)
	}

}
func PreflightHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
	  c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	  c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	  c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
  
	  if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	  }
	  c.Next()
	}
  }
